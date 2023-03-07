package prattle

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"unicode"
)

var wordlist = []string{
	"aardvark",
	"babbling",
	"cabin",
	"dachshund",
	"eagerness",
	"fabric",
	"gadget",
	"habitat",
	"ibuprofen",
	"jabbering",
	"kabob",
	"laboratory",
	"macaroni",
	"nail",
	"oasis",
	"pacemaker",
	"quarters",
	"racoon",
	"sabotage",
	"tablespoon",
	"ultrasound",
	"vacuum",
	"waffle",
	"xerox",
	"yacht",
	"zealot",
}

func genWords(n int, rng *rand.Rand) string {
	var sb strings.Builder
	j := rng.Intn(len(wordlist))
	sb.WriteString(wordlist[j])
	for i := 1; i < n; i++ {
		j := rng.Intn(len(wordlist))
		fmt.Fprintf(&sb, "%s ", wordlist[j])
	}
	return sb.String()
}

func scanWords(s *Scanner) int {
	s.ExpectAny(unicode.IsSpace)
	s.Skip()
	switch {
	case s.Done():
		return 0
	case s.ExpectOne(unicode.IsLetter):
		s.ExpectAny(unicode.IsLetter)
		return 1
	}
	s.Advance()
	return -1
}

func BenchmarkScannerWithString(b *testing.B) {
	b.ReportAllocs()
	rng := rand.New(rand.NewSource(0))
	words := genWords(2048, rng)
	s := Scanner{Scan: scanWords}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.InitWithString(words)
		for _, ok := s.Next(); ok; _, ok = s.Next() {
		}
	}
}

func BenchmarkScannerWithReader(b *testing.B) {
	b.ReportAllocs()
	rng := rand.New(rand.NewSource(0))
	words := genWords(2048, rng)
	s := Scanner{Scan: scanWords}
	r := strings.NewReader(words)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Reset(words)
		s.InitWithReader(r)
		for _, ok := s.Next(); ok; _, ok = s.Next() {
		}
	}
}
