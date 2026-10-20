package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func (s *Store) buatPeminjaman(w http.ResponseWriter, r *http.Request) {
	var body struct {
		SumberDayaID int    `json:"sumber_daya_id"`
		Peminjam     string `json:"peminjam"`
		Tanggal      string `json:"tanggal"`
		SlotMulai    int    `json:"slot_mulai"`
		SlotSelesai  int    `json:"slot_selesai"`
		Jumlah       int    `json:"jumlah"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, 400, "body tidak valid")
		return
	}

	body.Peminjam = trim(body.Peminjam)
	body.Tanggal = trim(body.Tanggal)

	if body.SumberDayaID < 1 {
		writeError(w, 400, "sumber_daya_id wajib diisi")
		return
	}
	if !validNama(body.Peminjam) {
		writeError(w, 400, "peminjam wajib diisi, maksimal 50 karakter")
		return
	}
	if !validTanggal(body.Tanggal) {
		writeError(w, 400, "tanggal harus format YYYY-MM-DD dan benar")
		return
	}
	if !validSlot(body.SlotMulai, body.SlotSelesai) {
		writeError(w, 400, "slot harus 1 sampai 8 dan slot_mulai <= slot_selesai")
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	var sd *SumberDaya
	for i := range s.sumberDaya {
		if s.sumberDaya[i].ID == body.SumberDayaID {
			sd = &s.sumberDaya[i]
			break
		}
	}
	if sd == nil {
		writeError(w, 400, "sumber daya tidak ditemukan")
		return
	}
	if body.Jumlah < 1 || body.Jumlah > sd.Stok {
		writeError(w, 400, "jumlah antara 1 dan stok sumber daya")
		return
	}

	for slot := body.SlotMulai; slot <= body.SlotSelesai; slot++ {
		terpakai := 0
		for _, p := range s.peminjaman {
			if p.SumberDayaID == body.SumberDayaID && p.Tanggal == body.Tanggal && p.Status == "aktif" {
				for s2 := p.SlotMulai; s2 <= p.SlotSelesai; s2++ {
					if s2 == slot {
						terpakai += p.Jumlah
					}
				}
			}
		}
		if terpakai+body.Jumlah > sd.Stok {
			writeError(w, 409, "kapasitas slot terlampaui stok")
			return
		}
	}

	for _, b := range s.blokir {
		if b.SumberDayaID == body.SumberDayaID && b.Tanggal == body.Tanggal &&
			slotOverlap(b.SlotMulai, b.SlotSelesai, body.SlotMulai, body.SlotSelesai) {
			writeError(w, 409, "slot diblokir")
			return
		}
	}

	p := Peminjaman{
		ID:           s.nextPeminjamanID,
		SumberDayaID: body.SumberDayaID,
		Peminjam:     body.Peminjam,
		Tanggal:      body.Tanggal,
		SlotMulai:    body.SlotMulai,
		SlotSelesai:  body.SlotSelesai,
		Jumlah:       body.Jumlah,
		Status:       "aktif",
	}
	s.nextPeminjamanID++
	s.peminjaman = append(s.peminjaman, p)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(p)
}

func (s *Store) daftarPeminjaman(w http.ResponseWriter, r *http.Request) {
	sumberDayaIDStr := r.URL.Query().Get("sumber_daya_id")
	tanggal := trim(r.URL.Query().Get("tanggal"))
	status := trim(r.URL.Query().Get("status"))

	var sumberDayaID *int
	if sumberDayaIDStr != "" {
		id, err := strconv.Atoi(sumberDayaIDStr)
		if err != nil {
			writeError(w, 400, "sumber_daya_id harus angka")
			return
		}
		sumberDayaID = &id
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	var hasil []Peminjaman
	for _, p := range s.peminjaman {
		if sumberDayaID != nil && p.SumberDayaID != *sumberDayaID {
			continue
		}
		if tanggal != "" && p.Tanggal != tanggal {
			continue
		}
		if status != "" && p.Status != status {
			continue
		}
		hasil = append(hasil, p)
	}
	if hasil == nil {
		hasil = []Peminjaman{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hasil)
}

func (s *Store) detailPeminjaman(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, 400, "id harus angka")
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, p := range s.peminjaman {
		if p.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(p)
			return
		}
	}

	writeError(w, 404, "peminjaman tidak ditemukan")
}

func (s *Store) ubahPeminjaman(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, 400, "id harus angka")
		return
	}

	var body struct {
		SumberDayaID int    `json:"sumber_daya_id"`
		Peminjam     string `json:"peminjam"`
		Tanggal      string `json:"tanggal"`
		SlotMulai    int    `json:"slot_mulai"`
		SlotSelesai  int    `json:"slot_selesai"`
		Jumlah       int    `json:"jumlah"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, 400, "body tidak valid")
		return
	}

	body.Peminjam = trim(body.Peminjam)
	body.Tanggal = trim(body.Tanggal)

	if body.SumberDayaID < 1 {
		writeError(w, 400, "sumber_daya_id wajib diisi")
		return
	}
	if !validNama(body.Peminjam) {
		writeError(w, 400, "peminjam wajib diisi, maksimal 50 karakter")
		return
	}
	if !validTanggal(body.Tanggal) {
		writeError(w, 400, "tanggal harus format YYYY-MM-DD dan benar")
		return
	}
	if !validSlot(body.SlotMulai, body.SlotSelesai) {
		writeError(w, 400, "slot harus 1 sampai 8 dan slot_mulai <= slot_selesai")
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	var idx int
	var lama *Peminjaman
	for i := range s.peminjaman {
		if s.peminjaman[i].ID == id {
			idx = i
			lama = &s.peminjaman[i]
			break
		}
	}
	if lama == nil {
		writeError(w, 404, "peminjaman tidak ditemukan")
		return
	}

	var sd *SumberDaya
	for i := range s.sumberDaya {
		if s.sumberDaya[i].ID == body.SumberDayaID {
			sd = &s.sumberDaya[i]
			break
		}
	}
	if sd == nil {
		writeError(w, 400, "sumber daya tidak ditemukan")
		return
	}
	if body.Jumlah < 1 || body.Jumlah > sd.Stok {
		writeError(w, 400, "jumlah antara 1 dan stok sumber daya")
		return
	}

	for slot := body.SlotMulai; slot <= body.SlotSelesai; slot++ {
		terpakai := 0
		for _, p := range s.peminjaman {
			if p.ID == id {
				continue
			}
			if p.SumberDayaID == body.SumberDayaID && p.Tanggal == body.Tanggal && p.Status == "aktif" {
				for s2 := p.SlotMulai; s2 <= p.SlotSelesai; s2++ {
					if s2 == slot {
						terpakai += p.Jumlah
					}
				}
			}
		}
		if terpakai+body.Jumlah > sd.Stok {
			writeError(w, 409, "kapasitas slot terlampaui stok")
			return
		}
	}

	for _, b := range s.blokir {
		if b.SumberDayaID == body.SumberDayaID && b.Tanggal == body.Tanggal &&
			slotOverlap(b.SlotMulai, b.SlotSelesai, body.SlotMulai, body.SlotSelesai) {
			writeError(w, 409, "slot diblokir")
			return
		}
	}

	s.peminjaman[idx] = Peminjaman{
		ID:           id,
		SumberDayaID: body.SumberDayaID,
		Peminjam:     body.Peminjam,
		Tanggal:      body.Tanggal,
		SlotMulai:    body.SlotMulai,
		SlotSelesai:  body.SlotSelesai,
		Jumlah:       body.Jumlah,
		Status:       lama.Status,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.peminjaman[idx])
}

func (s *Store) hapusPeminjaman(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, 400, "id harus angka")
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	var idx int
	var found bool
	for i := range s.peminjaman {
		if s.peminjaman[i].ID == id {
			idx = i
			found = true
			break
		}
	}
	if !found {
		writeError(w, 404, "peminjaman tidak ditemukan")
		return
	}

	s.peminjaman = append(s.peminjaman[:idx], s.peminjaman[idx+1:]...)

	w.WriteHeader(204)
}
