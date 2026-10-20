# AI-LOG

Catat setiap penggunaan AI yang menghasilkan kode. Jika tidak memakai AI sama sekali, tulis satu baris: **tidak memakai AI**.

Log yang jujur tidak mengurangi nilai, sebanyak apa pun AI dipakai. Yang mengurangi nilai adalah log yang tidak sesuai dengan kode, atau isi log yang tidak bisa dijelaskan.

## Store dan struktur data (store.go)
- Prompt: "buat struct Go untuk sumber daya, peminjaman, dan blokir dengan tag JSON sesuai kontrak PinjamKampus, plus Store dengan mutex dan id otomatis"
- Dipakai: struct SumberDaya, Peminjaman, Blokir, dan Store di store.go
- Diubah: tidak ada, kode diperiksa ulang sendiri agar pas dengan kontrak dan nama field yang tepat.

## Helper error JSON (errors.go)
- Prompt: "fungsi Go untuk menulis response error JSON dengan format {"error":"pesan"} dan kode status HTTP"
- Dipakai: fungsi writeError di errors.go
- Diubah: tidak ada, hanya menyesuaikan penulisan map agar sederhana.

## Validasi input dasar (validate.go)
- Prompt: "fungsi validasi nama maksimal 50 karakter dan jenis ruang/alat di Go"
- Dipakai: validNama, validJenis di validate.go
- Diubah: ditambahkan fungsi trim untuk membersihkan spasi sebelum dicek, sesuai aturan teks selalu di-trim. Ditambahkan validTanggal (time.Parse layout 2006-01-02) dan validSlot (1-8, mulai <= selesai) untuk B2.

## Handler sumber daya B1 (handler_sumberdaya.go)
- Prompt: "handler Go net/http untuk POST, GET daftar, dan GET detail sumber daya dengan validasi nama, jenis, stok"
- Dipakai: buatSumberDaya, daftarSumberDaya, detailSumberDaya di handler_sumberdaya.go
- Diubah: menambahkan lock mutex pada semua operasi, pathValue untuk parsing id, dan urutan cek id bukan angka (400) sebelum not found (404).

## Handler peminjaman B2+B3+B5+B6 (handler_peminjaman.go)
- Prompt: "handler Go net/http untuk POST, GET daftar, GET detail, PUT ubah, DELETE hapus peminjaman dengan validasi 6 field, deteksi bentrok per slot vs stok, PUT tidak bentrok dengan diri sendiri, respons 201/200/204/400/404/409"
- Dipakai: buatPeminjaman, daftarPeminjaman, detailPeminjaman, ubahPeminjaman, hapusPeminjaman di handler_peminjaman.go
- Diubah: deteksi bentrok dilakukan per slot dalam rentang, hanya hitung peminjaman status "aktif", di PUT skip peminjaman dengan id yang sedang diubah, semua dalam satu kunci mutex. PUT tidak mengubah id maupun status. Filter GET /peminjaman supports sumber_daya_id, tanggal, status. DELETE mengembalikan 204 tanpa body.

## Handler blokir jadwal Persona E (handler_blokir.go)
- Prompt: "handler Go net/http untuk POST, GET, DELETE blokir jadwal dengan validasi, cek irisan blokir, batalkan peminjaman yang terpengaruh, respons 201/200/204/400/404/409"
- Dipakai: buatBlokir, daftarBlokir, hapusBlokir di handler_blokir.go; fungsi slotOverlap juga digunakan
- Diubah: buatBlokir membatalkan peminjaman aktif yang beririsan slot (status → "dibatalkan"), respons 201 berisi {blokir, dibatalkan: []int}. Cek blokir juga ditambahkan di buatPeminjaman dan ubahPeminjaman (handler_peminjaman.go) → 409 jika slot bertabrakan blokir.

## Router (router.go dan main.go)
- Prompt: "router Go http.NewServeMux dengan NewHandler() http.Handler dan main yang jalankan ListenAndServe :8080"
- Dipakai: NewHandler di router.go, main di main.go
- Diubah: menambahkan route POST/GET/DELETE untuk /blokir dan /blokir/{id}.

## Verifikasi
- go vet ./... dan go build ./... berhasil tanpa peringatan. Uji manual pakai curl untuk semua kasus B1–B7 dan Persona E lolos (201, 200, 204, 400, 404, 409; konkurensi 1×201 + 19×409 dari 20 request paralel; blokir membatalkan peminjaman dan mencegah peminjaman baru).
