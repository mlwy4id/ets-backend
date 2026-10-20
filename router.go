package main

import "net/http"

func NewHandler() http.Handler {
	s := NewStore()
	mux := http.NewServeMux()

	mux.HandleFunc("POST /sumber-daya", s.buatSumberDaya)
	mux.HandleFunc("GET /sumber-daya", s.daftarSumberDaya)
	mux.HandleFunc("GET /sumber-daya/{id}", s.detailSumberDaya)

	mux.HandleFunc("POST /peminjaman", s.buatPeminjaman)
	mux.HandleFunc("GET /peminjaman", s.daftarPeminjaman)
	mux.HandleFunc("GET /peminjaman/{id}", s.detailPeminjaman)
	mux.HandleFunc("PUT /peminjaman/{id}", s.ubahPeminjaman)
	mux.HandleFunc("DELETE /peminjaman/{id}", s.hapusPeminjaman)

	mux.HandleFunc("POST /blokir", s.buatBlokir)
	mux.HandleFunc("GET /blokir", s.daftarBlokir)
	mux.HandleFunc("DELETE /blokir/{id}", s.hapusBlokir)

	return mux
}
