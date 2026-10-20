package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func (s *Store) buatSumberDaya(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Nama  string `json:"nama"`
		Jenis string `json:"jenis"`
		Stok  int    `json:"stok"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, 400, "body tidak valid")
		return
	}

	body.Nama = trim(body.Nama)
	body.Jenis = trim(body.Jenis)

	if !validNama(body.Nama) {
		writeError(w, 400, "nama wajib diisi, maksimal 50 karakter")
		return
	}
	if !validJenis(body.Jenis) {
		writeError(w, 400, "jenis harus ruang atau alat")
		return
	}
	if body.Stok < 1 {
		writeError(w, 400, "stok minimal 1")
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	sd := SumberDaya{
		ID:    s.nextSumberDayaID,
		Nama:  body.Nama,
		Jenis: body.Jenis,
		Stok:  body.Stok,
	}
	s.nextSumberDayaID++
	s.sumberDaya = append(s.sumberDaya, sd)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(sd)
}

func (s *Store) daftarSumberDaya(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.sumberDaya == nil {
		s.sumberDaya = []SumberDaya{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.sumberDaya)
}

func (s *Store) detailSumberDaya(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, 400, "id harus angka")
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, sd := range s.sumberDaya {
		if sd.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(sd)
			return
		}
	}

	writeError(w, 404, "sumber daya tidak ditemukan")
}
