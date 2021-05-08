package prattle

import (
	"strings"
	"unicode/utf8"
)

// ScanFunc is a function that lexically analyses the Scanner input and emits tokens.
type ScanFunc func(*Scanner) ScanFunc

// AcceptFunc accepts a rune.
type AcceptFunc func(rune) bool

// Scanner produces a sequence of tokens from an input string.
// Scanner is UTF-8 aware and consumes a single codepoint at a time.
//
// The design is inspired by Rob Pike's lexer presentation.
type Scanner struct {
	name string

	// input is the string to be scanned.
	input string

	scan    ScanFunc
	peek    rune
	peekw   int
	offset  int
	ofsline int
	ofscoln int
	cursor  int
	curline int
	curcoln int
	tokens  chan Token
}

// NewScanner creates a new scanner.
func NewScanner(name, input string, scan ScanFunc) *Scanner {
	s := &Scanner{
		name:   name,
		input:  input,
		scan:   scan,
		tokens: make(chan Token),
	}
	_ = s.Advance()
	go s.pump()
	return s
}

func (s *Scanner) pump() {
	for scan := s.scan; scan != nil; scan = scan(s) {
	}
	close(s.tokens)
}

// Flush flushes any remaining tokens.
func (s *Scanner) Flush() {
	for range s.tokens {
	}
}

// Next returns the next token in the token stream
// or the zero token when the stream is exhausted.
func (s *Scanner) Next() Token {
	token, _ := <-s.tokens
	return token
}

// Skip swallows the next token.
func (s *Scanner) Skip() {
	s.offset = s.cursor
	s.ofsline = s.curline
	s.ofscoln = s.curcoln
}

// Emit produces the next token.
func (s *Scanner) Emit(kind Kind) {
	text := s.input[s.offset:s.cursor]
	t := Token{
		Kind: kind,
		Text: text,
		Position: Position{
			Filename: s.name,
			Offset:   s.offset,
			Line:     s.ofsline,
			Column:   s.ofscoln,
		},
	}
	s.Skip()
	s.tokens <- t
}

// Peek returns the current rune.
func (s *Scanner) Peek() rune {
	return s.peek
}

// Text returns the token text scanned so far.
func (s *Scanner) Text() string {
	return s.input[s.offset:s.cursor]
}

// Done returns true if the input has been consumed.
func (s *Scanner) Done() bool {
	return s.offset >= len(s.input)
}

// Advance advances the cursor by one rune.
func (s *Scanner) Advance() rune {
	if s.curline == 0 || s.peek == '\n' {
		s.curline++
		s.curcoln = 1
	} else if s.peekw > 0 {
		s.curcoln++
	}

	s.cursor += s.peekw
	s.peek, s.peekw = utf8.DecodeRuneInString(s.input[s.cursor:])
	return s.peek
}

// OneRune advances the cursor if the current rune matches.
func (s *Scanner) OneRune(r rune) bool {
	if s.Peek() == r {
		_ = s.Advance()
		return true
	}
	return false
}

// One advances the cursor if the current rune matches the accept function.
func (s *Scanner) One(acceptFunc AcceptFunc) bool {
	if acceptFunc(s.Peek()) {
		_ = s.Advance()
		return true
	}
	return false
}

// Any advances the cursor for as long as it matches the accept function.
func (s *Scanner) Any(acceptFunc AcceptFunc) {
	for acceptFunc(s.Peek()) {
		_ = s.Advance()
	}
}

// OneOf returns an AcceptFunc that accepts runes that appear in s.
func OneOf(s string) AcceptFunc {
	return func(r rune) bool {
		return strings.IndexRune(s, r) >= 0
	}
}
