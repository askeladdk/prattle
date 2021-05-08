package prattle

import (
	"testing"
	"unicode"
)

func TestScanner(t *testing.T) {
	scan := func(s *Scanner) ScanFunc {
		for !s.Done() {
			s.Any(unicode.IsSpace)
			s.Skip()
			switch {
			case s.One(unicode.IsLetter):
				s.Any(unicode.IsLetter)
				s.Emit(1)
			case s.One(unicode.IsDigit):
				s.Any(unicode.IsDigit)
				s.Emit(2)
			case s.One(OneOf("+=")):
				s.Emit(3)
			case s.OneRune('€'):
				s.Emit(4)
			}
		}

		s.Emit(0)
		return nil
	}

	expected := []Token{
		{Kind: 1, Text: "result", Position: Position{Line: 1, Column: 1}},
		{Kind: 3, Text: "=", Position: Position{Line: 1, Column: 8}},
		{Kind: 4, Text: "€", Position: Position{Line: 2, Column: 1}},
		{Kind: 2, Text: "1337", Position: Position{Line: 2, Column: 2}},
		{Kind: 3, Text: "+", Position: Position{Line: 2, Column: 7}},
		{Kind: 1, Text: "BlAbLa", Position: Position{Line: 2, Column: 9}},
		{Position: Position{Line: 2, Column: 15}},
	}

	s := NewScanner("", "result =\n€1337 + BlAbLa", scan)
	defer s.Flush()

	for _, x := range expected {
		tok := s.Next()
		if x.Kind != tok.Kind || x.Text != tok.Text || x.Line != tok.Line || x.Column != tok.Column {
			t.Fatal(tok)
		}
	}
}

func TestScannerPrefixes(t *testing.T) {
	scan := func(s *Scanner) ScanFunc {
		for !s.Done() {
			s.Any(unicode.IsSpace)
			s.Skip()
			if s.OneRune('+') {
				if s.OneRune('+') {
					s.Emit(2)
				} else {
					s.Emit(1)
				}
			}
		}

		s.Emit(0)
		return nil
	}

	expected := []Kind{1, 2, 2, 1, 2, 2, 0}

	s := NewScanner("", "+ ++ +++ ++++", scan)
	defer s.Flush()

	for _, x := range expected {
		tok := s.Next()
		if tok.Kind != x {
			t.Fatal(tok)
		}
	}
}
