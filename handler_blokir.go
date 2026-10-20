package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func slotOverlap(aMulai, aSelesai, bMulai, bSelesai int) bool {
	return aMulai <= bSelesai && bMulai <= aSelesai
}

func (s *Store) buatBlokir(w http.ResponseWriter, r *http.Request) {
	var body struct {
		SumberDayaID int    `json:"sumber_daya_id"`
		Tanggal      string `json:"tanggal"`
		SlotMulai    int    `json:"slot_mulai"`
		SlotSelesai  int    `json:"slot_selesai"`
		Keterangan   string `json:"keterangan"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, 400, "body tidak valid")
		return
	}

	body.Tanggal = trim(body.Tanggal)
	body.Keterangan = trim(body.Keterangan)

	if body.SumberDayaID < 1 {
		writeError(w, 400, "sumber_daya_id wajib diisi")
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
	if body.Keterangan == "" {
		writeError(w, 400, "keterangan wajib diisi")
		return
	}
	if utf8RuneCount(body.Keterangan) > 100 {
		writeError(w, 400, "keterangan maksimal 100 karakter")
		return
	}

	for _, b := range s.blokir {
		if b.SumberDayaID == body.SumberDayaID && b.Tanggal == body.Tanggal &&
			slotOverlap(b.SlotMulai, b.SlotSelesai, body.SlotMulai, body.SlotSelesai) {
			writeError(w, 409, "slot beririsan dengan blokir lain")
			return
		}
	}

	bl := Blokir{
		ID:           s.nextBlokirID,
		SumberDayaID: body.SumberDayaID,
		Tanggal:      body.Tanggal,
		SlotMulai:    body.SlotMulai,
		SlotSelesai:  body.SlotSelesai,
		Keterangan:   body.Keterangan,
	}
	s.nextBlokirID++
	s.blokir = append(s.blokir, bl)

	var dibatalkan []int
	for i := range s.peminjaman {
		p := &s.peminjaman[i]
		if p.SumberDayaID == body.SumberDayaID && p.Tanggal == body.Tanggal &&
			p.Status == "aktif" &&
			slotOverlap(p.SlotMulai, p.SlotSelesai, body.SlotMulai, body.SlotSelesai) {
			p.Status = "dibatalkan"
			dibatalkan = append(dibatalkan, p.ID)
		}
	}
	if dibatalkan == nil {
		dibatalkan = []int{}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"blokir":      bl,
		"dibatalkan": dibatalkan,
	})
}

func (s *Store) daftarBlokir(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.blokir == nil {
		s.blokir = []Blokir{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.blokir)
}

func (s *Store) hapusBlokir(w http.ResponseWriter, r *http.Request) {
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
	for i := range s.blokir {
		if s.blokir[i].ID == id {
			idx = i
			found = true
			break
		}
	}
	if !found {
		writeError(w, 404, "blokir tidak ditemukan")
		return
	}

	s.blokir = append(s.blokir[:idx], s.blokir[idx+1:]...)

	w.WriteHeader(204)
}

func utf8RuneCount(s string) int {
	return len([]rune(s))
}
