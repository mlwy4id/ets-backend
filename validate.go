package main

import (
	"strings"
	"time"
	"unicode/utf8"
)

func trim(s string) string {
	return strings.TrimSpace(s)
}

func validNama(s string) bool {
	s = trim(s)
	return s != "" && utf8.RuneCountInString(s) <= 50
}

func validJenis(s string) bool {
	return s == "ruang" || s == "alat"
}

func validTanggal(s string) bool {
	s = trim(s)
	_, err := time.Parse("2006-01-02", s)
	return err == nil
}

func validSlot(mulai, selesai int) bool {
	return mulai >= 1 && selesai <= 8 && mulai <= selesai
}
