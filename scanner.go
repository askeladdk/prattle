package prattle

import (
	"io"
	"strings"
)

// ScanFunc scans the next token and returns its kind.
// By convention, zero is reserved to signal end-of-input
// and negative values signal an invalid token.
type ScanFunc func(*Scanner) (kind int)

// AcceptFunc accepts a rune.
type AcceptFunc func(rune) bool

// Scanner produces a stream of tokens from a string or an io.RuneReader.
type Scanner struct {
	// Position of the last read token.
	// The Filename field will never be modified by Scanner.
	Position

	// Scan scans tokens.
	Scan ScanFunc

	reader  spanReader
	peek    rune
	peekw   int
	cursor  int
	curline int
	curcoln int
}

func (s *Scanner) init(r spanReader) *Scanner {
	s.Offset = 0
	s.Column = 0
	s.Line = 0
	s.reader = r
	s.peek = 0
	s.peekw = 0
	s.cursor = 0
	s.curline = 1
	s.curcoln = 1
	s.Advance()
	return s
}

// Init initializes a Scanner with an input string and returns it.
// When initialized this way, the Scanner tokenizes without allocations.
func (s *Scanner) InitWithString(source string) *Scanner {
	// reuse existing stringSpanner if possible
	if ss, ok := s.reader.(*stringSpanner); ok {
		ss.source = source
		ss.cursor = 0
		ss.size = 0
		return s.init(ss)
	}

	return s.init(&stringSpanner{source: source})
}

// Init initializes a Scanner with an input reader and returns it.
// When initialized this way, the Scanner tokenizes unbounded streams but must allocate.
func (s *Scanner) InitWithReader(r io.RuneReader) *Scanner {
	// reuse existing runeReaderSpanner if possible
	if rrs, ok := s.reader.(*runeReaderSpanner); ok {
		rrs.rr = r
		rrs.buf = rrs.buf[:0]
		rrs.r = 0
		return s.init(rrs)
	}

	return s.init(&runeReaderSpanner{rr: r})
}

// Text returns the token string that has been scanned so far.
func (s *Scanner) Text() string {
	return s.reader.Span()
}

// Next implements Iterator.
func (s *Scanner) Next() (Token, bool) {
	var tok Token
	tok.Kind = s.Scan(s)
	tok.Text = s.Text()
	tok.Position = s.Position
	s.Skip()
	return tok, tok.Kind != 0
}

// Skip swallows the next token.
func (s *Scanner) Skip() {
	s.Offset = s.cursor
	s.Line = s.curline
	s.Column = s.curcoln
	s.reader.NextSpan()
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
	if s.peek == '\n' {
		s.curline++
		s.curcoln = 1
	} else if s.peekw > 0 {
		s.curcoln++
	}

	s.cursor += s.peekw
	s.peek, s.peekw, _ = s.reader.ReadRune()
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
