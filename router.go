package main

import "net/http"

func NewHandler() http.Handler {
	s := NewStore()
	mux := http.NewServeMux()

	mux.HandleFunc("POST /sumber-daya", s.buatSumberDaya)
	mux.HandleFunc("GET /sumber-daya", s.daftarSumberDaya)
	mux.HandleFunc("GET /sumber-daya/{id}", s.detailSumberDaya)

	return mux
}
