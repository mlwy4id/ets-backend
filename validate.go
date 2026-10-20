package main

import (
	"strings"
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
