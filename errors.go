package main

import (
	"encoding/json"
	"net/http"
)

func writeError(w http.ResponseWriter, kode int, pesan string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(kode)
	json.NewEncoder(w).Encode(map[string]string{"error": pesan})
}
