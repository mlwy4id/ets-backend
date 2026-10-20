# AI-LOG

Catat setiap penggunaan AI yang menghasilkan kode. Jika tidak memakai AI sama sekali, tulis satu baris: **tidak memakai AI**.

Log yang jujur tidak mengurangi nilai, sebanyak apa pun AI dipakai. Yang mengurangi nilai adalah log yang tidak sesuai dengan kode, atau isi log yang tidak bisa dijelaskan.

## Store dan struktur data (store.go)
- Prompt: "buat struct Go untuk sumber daya, peminjaman, dan blokir dengan tag JSON sesuai kontrak PinjamKampus, plus Store dengan mutex dan id otomatis"
- Dipakai: struct SumberDaya, Peminjaman, Blokir, dan Store di store.go
- Diubah: -

## Helper error JSON (errors.go)
- Prompt: "fungsi Go untuk menulis response error JSON dengan format {"error":"pesan"} dan kode status HTTP"
- Dipakai: fungsi writeError di errors.go
- Diubah: -

## Validasi input dasar (validate.go)
- Prompt: "fungsi validasi nama maksimal 50 karakter dan jenis ruang/alat di Go"
- Dipakai: validNama, validJenis di validate.go
- Diubah: -

## Handler sumber daya B1 (handler_sumberdaya.go)
- Prompt: "handler Go net/http untuk POST, GET daftar, dan GET detail sumber daya dengan validasi nama, jenis, stok"
- Dipakai: buatSumberDaya, daftarSumberDaya, detailSumberDaya di handler_sumberdaya.go
- Diubah: -

## Handler peminjaman B2+B3+B5+B6 (handler_peminjaman.go)
- Prompt: "handler Go net/http untuk POST, GET daftar, GET detail, PUT ubah, DELETE hapus peminjaman dengan validasi 6 field, deteksi bentrok per slot vs stok, PUT tidak bentrok dengan diri sendiri, respons 201/200/204/400/404/409"
- Dipakai: buatPeminjaman, daftarPeminjaman, detailPeminjaman, ubahPeminjaman, hapusPeminjaman di handler_peminjaman.go
- Diubah: -

## Handler blokir jadwal Persona E (handler_blokir.go)
- Prompt: "handler Go net/http untuk POST, GET, DELETE blokir jadwal dengan validasi, cek irisan blokir, batalkan peminjaman yang terpengaruh, respons 201/200/204/400/404/409"
- Dipakai: buatBlokir, daftarBlokir, hapusBlokir di handler_blokir.go; fungsi slotOverlap juga digunakan
- Diubah: -

## Change Request C: DELETE peminjaman jadi pembatalan (soft delete) (handler_peminjaman.go)
- Prompt: "ubah DELETE /peminjaman/{id} dari hard delete menjadi soft delete: set status dibatalkan, data tetap ada, cegah double-cancel dengan 409"
- Dipakai: hapusPeminjaman di handler_peminjaman.go
- Diubah: -

## Router (router.go dan main.go)
- Prompt: "router Go http.NewServeMux dengan NewHandler() http.Handler dan main yang jalankan ListenAndServe :8080"
- Dipakai: NewHandler di router.go, main di main.go
- Diubah: -

## Verifikasi
- go vet ./... dan go build ./... berhasil tanpa peringatan. Uji manual pakai curl untuk semua kasus B1–B7, Persona E, dan CR-C lolos (201, 200, 204, 400, 404, 409; konkurensi 1×201 + 19×409 dari 20 request paralel; blokir membatalkan peminjaman dan mencegah peminjaman baru; DELETE peminjaman → 204, GET → 200 status dibatalkan, double delete → 409, slot tidak ditahan oleh yang dibatalkan).
