package prattle_test

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/askeladdk/prattle"
)

const (
	kident prattle.Kind = 1 + iota
	knumber
	kassign
	kplus
	ksemicolon
)

func testScan(s *prattle.Scanner) prattle.Kind {
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

func (c *testContext) Precedence(kind prattle.Kind) prattle.Precedence {
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

func (c *testContext) ident(p *prattle.Parser, t prattle.Token) error {
	v := c.idents[t.Text]
	c.push(v)
	return nil
}

func (c *testContext) number(p *prattle.Parser, t prattle.Token) error {
	v, _ := strconv.Atoi(t.Text)
	c.push(v)
	return nil
}

func (c *testContext) plus(p *prattle.Parser, t prattle.Token) error {
	if err := p.ParseExpression(c.Precedence(t.Kind)); err != nil {
		return err
	}
	right := c.pop()
	left := c.pop()
	c.push(right + left)
	return nil
}

func (c *testContext) assign(p *prattle.Parser, t prattle.Token) error {
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

func (c *testContext) Prefix(kind prattle.Kind) prattle.ParseFunc {
	switch kind {
	default:
		return nil
	case kident:
		return c.ident
	case knumber:
		return c.number
	}
}

func (c *testContext) Infix(kind prattle.Kind) prattle.ParseFunc {
	switch kind {
	default:
		return nil
	case kplus:
		return c.plus
	}
}

func (c *testContext) Statement(kind prattle.Kind) prattle.ParseFunc {
	switch kind {
	case kident:
		return c.assign
	default:
		return nil
	}
}

func (c *testContext) ParseError(t prattle.Token) error {
	return fmt.Errorf("%s: unexpected '%s'", t.Position, t.Text)
}

// This example interprets a simple programming language that only has assignment statements and addition.
func Example_interpreter() {
	var c testContext
	c.idents = map[string]int{}

	source := "a = 1;\nb = 2;\nc = a+b+b+a;\n"

	s := prattle.NewScanner(strings.NewReader(source), testScan)
	p := prattle.NewParser(s, &c)

	accept := func(k prattle.Kind) bool {
		return k > 0
	}

	// Parse one or more statements.
	if err := p.ParseStatement(); err != nil {
		fmt.Println(err)
	} else if err := p.ParseStatements(accept); err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("c = %d\n", c.idents["c"])
	}

	// Output:
	// c = 6
}
