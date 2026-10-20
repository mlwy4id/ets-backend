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

## Handler peminjaman B2 (handler_peminjaman.go)
- Prompt: "handler Go net/http untuk POST peminjaman dengan validasi 6 field, deteksi bentrok per slot vs stok, respons 201/400/409"
- Dipakai: buatPeminjaman, daftarPeminjaman, detailPeminjaman di handler_peminjaman.go
- Diubah: deteksi bentrok dilakukan per slot dalam rentang, hanya hitung peminjaman status "aktif", semua dalam satu kunci mutex. Filter GET /peminjaman supports sumber_daya_id, tanggal, status.

## Router (router.go dan main.go)
- Prompt: "router Go http.NewServeMux dengan NewHandler() http.Handler dan main yang jalankan ListenAndServe :8080"
- Dipakai: NewHandler di router.go, main di main.go
- Diubah: menambahkan route POST/GET/GET{id} untuk /sumber-daya dan /peminjaman.

## Verifikasi
- go vet ./... dan go build ./... berhasil tanpa peringatan. Uji manual pakai curl untuk semua kasus B1 dan B2 lolos (201, 400 nama/jenis/stok/tanggal/slot/sumber daya/jumlah, 409 bentrok, 404, 400 id bukan angka, 400 filter sumber_daya_id non-angka).
