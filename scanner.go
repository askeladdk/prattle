package prattle

import (
	"io"
	"strings"
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

	source  io.RuneReader
	buffer  []rune
	peek    rune
	peekw   int
	cursor  int
	curline int
	curcoln int
	err     error
}

// Init initializes a Scanner with a new source input and returns it.
// Panics if Scan or source is nil.
func (s *Scanner) Init(source io.RuneReader) *Scanner {
	if s.Scan == nil || source == nil {
		panic("prattle.Scanner parameters cannot be nil")
	}

	if s.buffer == nil {
		s.buffer = make([]rune, 0, 256)
	} else {
		s.buffer = s.buffer[:0]
	}

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

// Err returns a non-nil value if the source reader returned an error.
func (s *Scanner) Err() error {
	return s.err
}

// Next returns the next token in the token stream.
func (s *Scanner) Next() (tok Token) {
	tok.Kind = s.Scan(s)
	tok.Text = string(s.buffer)
	tok.Position = s.Position
	s.Skip()
	return
}

// Skip swallows the next token.
func (s *Scanner) Skip() {
	s.Offset = s.cursor
	s.Line = s.curline
	s.Column = s.curcoln
	s.buffer = s.buffer[:0]
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

	if s.peekw > 0 {
		s.buffer = append(s.buffer, s.peek)
		s.cursor += s.peekw
	}

	s.peek, s.peekw, s.err = s.source.ReadRune()
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

// OneOf returns an AcceptFunc that accepts runes that appear in s.
func OneOf(s string) AcceptFunc {
	return func(r rune) bool {
		return strings.IndexRune(s, r) >= 0
	}
}
