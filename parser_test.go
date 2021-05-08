package prattle

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"unicode"
)

const (
	kident Kind = 1 + iota
	knumber
	kassign
	kplus
	ksemicolon
)

func testScan(s *Scanner) Kind {
	s.ExpectAny(unicode.IsSpace)
	s.Skip()

	switch {
	case s.Done():
		return 0
	case s.ExpectOne(unicode.IsLetter):
		s.ExpectAny(unicode.IsLetter)
		return kident
	case s.ExpectOne(unicode.IsDigit):
		s.ExpectAny(unicode.IsDigit)
		return knumber
	case s.Expect('='):
		return kassign
	case s.Expect('+'):
		return kplus
	case s.Expect(';'):
		return ksemicolon
	}
	s.Advance()
	return -1
}

type testContext struct {
	stack  []int
	idents map[string]int
}

func (c *testContext) pop() (v int) {
	if n := len(c.stack); n > 0 {
		v, c.stack = c.stack[n-1], c.stack[:n-1]
	}
	return v
}

func (c *testContext) push(v int) {
	c.stack = append(c.stack, v)
}

func (c *testContext) Precedence(kind Kind) Precedence {
	switch kind {
	default:
		return 0
	case kassign:
		return 1
	case kplus:
		return 2
	case kident:
		return 3
	}
}

func (c *testContext) ident(p *Parser, t Token) error {
	v := c.idents[t.Text]
	c.push(v)
	return nil
}

func (c *testContext) number(p *Parser, t Token) error {
	v, _ := strconv.Atoi(t.Text)
	c.push(v)
	return nil
}

func (c *testContext) plus(p *Parser, t Token) error {
	if err := p.ParseExpression(c.Precedence(t.Kind)); err != nil {
		return err
	}
	right := c.pop()
	left := c.pop()
	c.push(right + left)
	return nil
}

func (c *testContext) assign(p *Parser, t Token) error {
	if !p.Expect(kassign) {
		return c.ParseError(p.Peek())
	}

	if err := p.ParseExpression(c.Precedence(kassign)); err != nil {
		return err
	}

	c.idents[t.Text] = c.pop()

	if !p.Expect(ksemicolon) {
		return c.ParseError(p.Peek())
	}

	return nil
}

func (c *testContext) Prefix(kind Kind) ParseFunc {
	switch kind {
	default:
		return nil
	case kident:
		return c.ident
	case knumber:
		return c.number
	}
}

func (c *testContext) Infix(kind Kind) ParseFunc {
	switch kind {
	default:
		return nil
	case kplus:
		return c.plus
	}
}

func (c *testContext) Statement(kind Kind) ParseFunc {
	switch kind {
	case kident:
		return c.assign
	default:
		return nil
	}
}

func (c *testContext) ParseError(t Token) error {
	return fmt.Errorf("%s: error '%s'", t.Position, t.Text)
}

func TestParse(t *testing.T) {
	var c testContext
	c.idents = map[string]int{}

	source := "a = 1;\nb = 2;\nc = a+b+b+a;\n"
	s := (&Scanner{Scan: testScan}).Init(strings.NewReader(source))
	p := (&Parser{Context: &c}).Init(s)

	if err := p.ParseStatement(); err != nil {
		t.Fatal(err)
	}

	if err := p.ParseStatements(func(k Kind) bool {
		return k > 0
	}); err != nil {
		t.Fatal(err)
	}

	if c.idents["c"] != 6 {
		t.Fatal()
	}
}
