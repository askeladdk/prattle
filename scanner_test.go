package prattle

import (
	"bufio"
	"errors"
	"strings"
	"testing"
	"testing/iotest"
	"unicode"
)

func TestScanner(t *testing.T) {
	scan := func(s *Scanner) Kind {
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
	s := (&Scanner{Scan: scan}).Init(strings.NewReader(source))

	for _, x := range expected {
		tok := s.Next()
		if tok != x {
			t.Fatal(tok)
		}
	}
}

func TestScannerPrefixes(t *testing.T) {
	scan := func(s *Scanner) Kind {
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

	expected := []Kind{1, 2, 2, 1, 2, 2, 0}

	source := "+ ++ +++ ++++"
	s := (&Scanner{Scan: scan}).Init(strings.NewReader(source))

	for _, x := range expected {
		tok := s.Next()
		if tok.Kind != x {
			t.Fatal(tok)
		}
	}
}

func TestScannerErr(t *testing.T) {
	scan := func(s *Scanner) Kind {
		return 1
	}

	r := iotest.ErrReader(errors.New("i/o error"))
	s := (&Scanner{Scan: scan}).Init(bufio.NewReader(r))
	k := s.Next()
	if k.Kind != 1 || s.Err() == nil {
		t.Fatal(k, s.Err())
	}
}

func TestScannerPanic(t *testing.T) {
	defer func() {
		if s, _ := recover().(string); s == "" {
			t.Fatal()
		}
	}()
	(&Scanner{}).Init(nil)
}

func Test_matchKeyword(t *testing.T) {
	keywords := [][]rune{
		[]rune("a"),
		[]rune("i"),
		[]rune("if"),
		[]rune("ifelsd"),
		[]rune("ifelse"),
		[]rune("var"),
	}

	scan := func(s *Scanner) Kind {
		s.ExpectAny(unicode.IsSpace)
		s.Skip()
		switch {
		case s.Done():
			return 0
		case s.ExpectOne(unicode.IsLetter):
			s.ExpectAny(unicode.IsLetter)
			if i := Kind(s.MatchKeyword(keywords)); i >= 0 {
				return 1 + i
			}
			return -1
		}
		s.Advance()
		return -1
	}

	expected := []Kind{3, 5, -1, 6, -1}
	source := "if ifelse ifels var varr"

	s := NewScanner(strings.NewReader(source), scan)
	for _, x := range expected {
		tok := s.Next()
		if tok.Kind != x {
			t.Fatal(tok)
		}
	}
}
