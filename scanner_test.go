package prattle

import (
	"strings"
	"testing"
	"unicode"
)

func TestScanner(t *testing.T) {
	scan := func(s *Scanner) int {
		s.ExpectAny(unicode.IsSpace)
		s.Skip()
		switch {
		case s.Done():
			return 0
		case s.ExpectOne(unicode.IsLetter):
			s.ExpectAny(unicode.IsLetter)
			return 1
		case s.ExpectOne(unicode.IsDigit):
			s.ExpectAny(unicode.IsDigit)
			return 2
		case s.ExpectOne(OneOf("+=")):
			return 3
		case s.Expect('€'):
			return 4
		}
		s.Advance()
		return -1
	}

	expected := []Token{
		{Kind: 1, Text: "result", Position: Position{Line: 1, Column: 1, Offset: 0}},
		{Kind: 3, Text: "=", Position: Position{Line: 1, Column: 8, Offset: 7}},
		{Kind: 4, Text: "€", Position: Position{Line: 2, Column: 1, Offset: 9}},
		{Kind: 2, Text: "1337", Position: Position{Line: 2, Column: 2, Offset: 12}},
		{Kind: 3, Text: "+", Position: Position{Line: 2, Column: 7, Offset: 17}},
		{Kind: 1, Text: "BlAbLa", Position: Position{Line: 2, Column: 9, Offset: 19}},
		{Position: Position{Line: 2, Column: 15, Offset: 25}},
	}

	source := "result =\n€1337 + BlAbLa"
	s := (&Scanner{Scan: scan})
	s.InitWithString(source)
	s.InitWithString(source)

	for _, x := range expected {
		tok := s.Next()
		if tok != x {
			t.Fatal(tok)
		}
	}
}

func TestScannerPrefixes(t *testing.T) {
	scan := func(s *Scanner) int {
		s.ExpectAny(unicode.IsSpace)
		s.Skip()
		switch {
		case s.Done():
			return 0
		case s.Expect('+'):
			if s.Expect('+') {
				return 2
			}
			return 1
		}
		s.Advance()
		return -1
	}

	expected := []int{1, 2, 2, 1, 2, 2, 0}

	source := "+ ++ +++ ++++"
	s := (&Scanner{Scan: scan})
	s.InitWithReader(strings.NewReader(source))
	s.InitWithReader(strings.NewReader(source))

	for _, x := range expected {
		tok := s.Next()
		if tok.Kind != x {
			t.Fatal(tok)
		}
	}
}

func Test_matchKeyword(t *testing.T) {
	keywords := map[string]int{
		"a":      1,
		"i":      2,
		"if":     3,
		"ifelsd": 4,
		"ifelse": 5,
		"var":    6,
	}

	scan := func(s *Scanner) int {
		s.ExpectAny(unicode.IsSpace)
		s.Skip()
		switch {
		case s.Done():
			return 0
		case s.ExpectOne(unicode.IsLetter):
			s.ExpectAny(unicode.IsLetter)
			if i, ok := keywords[s.Text()]; ok {
				return i
			}
			return -1
		}
		s.Advance()
		return -1
	}

	expected := []int{3, 5, -1, 6, -1}
	source := "if ifelse ifels var varr"

	s := Scanner{Scan: scan}
	s.InitWithString(source)

	for _, x := range expected {
		tok := s.Next()
		if tok.Kind != x {
			t.Fatal(tok)
		}
	}
}
