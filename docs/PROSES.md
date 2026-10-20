# PROSES

## 1. Persona dan masalahnya

Persona yang didapat adalah E: Asisten laboratorium. Masalahnya: lab kadang dipakai praktikum atau sedang diservis pada jam tertentu, sehingga pada jam itu lab tidak boleh dipinjam dan peminjaman yang sudah ada harus dibatalkan. Fitur khas blokir jadwal memungkinkan asisten lab membuat blokir pada slot tertentu, sehingga peminjaman baru yang bertabrakan ditolak dan peminjaman aktif yang sudah ada di-slot tersebut otomatis berstatus dibatalkan.

## 2. Rancangan struktur data

Tiga slice di dalam struct `Store`, masing-masing dilindungi oleh satu `sync.Mutex`:
- `sumberDaya []SumberDaya` — menyimpan ruang/alat dengan field id, nama, jenis, stok. ID otomatis bertambah untuk setiap sumber daya baru.
- `peminjaman []Peminjaman` — menyimpan peminjaman dengan field id, sumber_daya_id, peminjam, tanggal, slot_mulai, slot_selesai, jumlah, status (aktif/dibatalkan/diajukan/disetujui/ditolak/selesai). ID otomatis bertambah.
- `blokir []Blokir` — menyimpan blokir jadwal dengan field id, sumber_daya_id, tanggal, slot_mulai, slot_selesai, keterangan. ID otomatis bertambah.

Counter ID tersimpan sebagai `nextSumberDayaID`, `nextPeminjamanID`, `nextBlokirID` agar setiap jenis data memiliki id sendiri yang mulai dari 1.

Slice dipilih karena data disimpan di memori dan jumlah data di hackathon kecil; pencarian linear sesuai kebutuhan. Semua operasi baca-tulis dilakukan dalam satu kunci mutex untuk kepastian konsisten saat permintaan bersamaan.

## 3. Satu keputusan desain

Pengecekan bentrok dilakukan **per slot** dalam rentang peminjaman, bukan hanya membandingkan ujung rentang (slot_mulai/slot_selesai). Alternatif yang dipertimbangkan adalah cek apakah dua rentang saling overlap; namun itu gagal menangani kasus stok > 1 di mana banyak peminjaman bisa menahan slot yang sama selama total tidak melebihi stok. Dengan pengecekan per slot, kita menghitung jumlah unit yang sudah terpakai pada setiap slot individu dan membandingkan dengan stok, sehingga akurat untuk stok lebih dari 1.

## 4. Bug atau kendala (minimal dua)

### Bug 1
- Gejala: Saat uji konkuren, dua request paralel ke slot yang sama dengan stok 1 keduanya mendapatkan 201.
- Cara menemukan: Mengirim 20 request bersamaan ke slot sama melalui curl background, lalu menghitung banyaknya 201.
- Penyebab dan solusi: Pemeriksaan bentrok dan penyimpanan belum dalam satu kunci mutex yang sama. Solusi: membungkus seluruh alur (cek bentrok + append peminjaman baru) dalam satu `s.mu.Lock()` / `defer s.mu.Unlock()`.

### Bug 2
- Gejala: PUT peminjaman ke slot yang sama dengan dirinya sendiri menghasilkan 409 meski seharusnya diizinkan (tidak boleh bentrok dengan diri sendiri).
- Cara menemukan: Mengubah peminjaman existentes ke slot yang sama lalu memastikan respons 200.
- Penyebab dan solusi: Pengecekan bentrok belum melewatkan peminjaman dengan id yang sedang diubah. Solusi: menambahkan `if p.ID == id { continue }` di dalam loop hitung terpakai pada `ubahPeminjaman`.

## 5. Change request

**Permintaan (CR-C):** `DELETE /peminjaman/{id}` diubah menjadi pembatalan (soft delete): status berubah menjadi `"dibatalkan"` dan data tetap ada, respons tetap 204 tanpa body. Peminjaman yang sudah dibatalkan tidak bisa dibatalkan lagi (409). Peminjaman dibatalkan tidak menahan slot.

**Bagian kode yang diubah:** `handler_peminjaman.go`, fungsi `hapusPeminjaman`. Diganti dari `append(s.peminjaman[:idx], s.peminjaman[idx+1:]...)` (hard delete) menjadi `s.peminjaman[idx].Status = "dibatalkan"` (soft delete), ditambah cek `if s.peminjaman[idx].Status == "dibatalkan"` → 409 sebelum mengubah status.

**Keputusan desain:** Memilih soft delete dengan field status `"dibatalkan"` daripada hard delete atau menambahkan field terpisah `deleted_at`. Alasan: kontrak Requirement sudah menentukan field `status` pada objek peminjaman, dan menambahkan status baru tidak memecah struktur JSON yang sudah disepakati frontend. Peminjaman berstatus `dibatalkan` sudah tidak dihitung dalam deteksi bentrok (cek hanya `status == "aktif"`), sehingga slot otomatis dibebaskan tanpa logika tambahan.
