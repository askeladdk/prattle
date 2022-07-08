package prattle

import (
	"io"
	"unicode/utf8"
)

type spanReader interface {
	io.RuneReader
	Span() string
	NextSpan()
}

type stringSpanner struct {
	source string
	cursor int
	size   int
}

func (s *stringSpanner) ReadRune() (r rune, size int, err error) {
	s.cursor += s.size
	r, size = utf8.DecodeRuneInString(s.source[s.cursor:])
	s.size = size
	return
}

func (s *stringSpanner) Span() string {
	return s.source[:s.cursor]
}

func (s *stringSpanner) NextSpan() {
	s.source = s.source[s.cursor:]
	s.cursor = 0
}

type runeReaderSpanner struct {
	rr  io.RuneReader
	buf []rune
	r   rune
}

func (s *runeReaderSpanner) ReadRune() (r rune, size int, err error) {
	s.buf = append(s.buf, s.r)
	r, size, err = s.rr.ReadRune()
	s.r = r
	return
}

func (s *runeReaderSpanner) Span() string {
	return string(s.buf)
}

func (s *runeReaderSpanner) NextSpan() {
	s.buf = s.buf[:0]
}
