package prattle

import (
	"strings"
	"unicode/utf8"
)

// ScanFunc returns the Kind of the next token.
type ScanFunc func(*Scanner) Kind

// AcceptFunc accepts a rune.
type AcceptFunc func(rune) bool

// Scanner produces a sequence of tokens from an io.RuneReader.
type Scanner struct {
	// Position of the last read token.
	Position

	// Scan scans tokens.
	Scan ScanFunc

	source  string
	peek    rune
	peekw   int
	cursor  int
	curline int
	curcoln int
	err     error
}

// Init initializes a Scanner with a new source input and returns it.
func (s *Scanner) Init(source string) *Scanner {
	s.Offset = 0
	s.Column = 0
	s.Line = 0
	s.source = source
	s.peek = 0
	s.peekw = 0
	s.cursor = 0
	s.curline = 0
	s.curcoln = 0
	s.err = nil

	s.Advance()
	return s
}

// Text returns the string matched so far.
func (s *Scanner) Text() string {
	return s.source[s.Offset:s.cursor]
}

// Err returns a non-nil value if the source reader returned an error.
func (s *Scanner) Err() error {
	return s.err
}

// Next returns the next token in the token stream.
func (s *Scanner) Next() Token {
	var tok Token
	tok.Kind = s.Scan(s)
	tok.Text = s.Text()
	tok.Position = s.Position
	s.Skip()
	return tok
}

// Skip swallows the next token.
func (s *Scanner) Skip() {
	s.Offset = s.cursor
	s.Line = s.curline
	s.Column = s.curcoln
}

// Peek returns the current rune.
func (s *Scanner) Peek() rune {
	return s.peek
}

// Done reports whether the entire input has been consumed.
func (s *Scanner) Done() bool {
	return s.peekw == 0
}

// Advance advances the cursor by one rune.
func (s *Scanner) Advance() {
	if s.err != nil {
		return
	}

	if s.curline == 0 || s.peek == '\n' {
		s.curline++
		s.curcoln = 1
	} else if s.peekw > 0 {
		s.curcoln++
	}

	s.cursor += s.peekw
	s.peek, s.peekw = utf8.DecodeRuneInString(s.source[s.cursor:])
}

// Expect advances the cursor if the current rune matches.
func (s *Scanner) Expect(r rune) bool {
	if s.Peek() == r {
		s.Advance()
		return true
	}
	return false
}

// ExpectOne advances the cursor if the current rune is accepted.
func (s *Scanner) ExpectOne(accept AcceptFunc) bool {
	if accept(s.Peek()) {
		s.Advance()
		return true
	}
	return false
}

// ExpectAny advances the cursor zero or more times.
func (s *Scanner) ExpectAny(accept AcceptFunc) {
	for accept(s.Peek()) {
		s.Advance()
	}
}

// OneOf returns an AcceptFunc that reports whether a rune appears in chars.
func OneOf(chars string) AcceptFunc {
	return func(r rune) bool {
		return strings.ContainsRune(chars, r)
	}
}
