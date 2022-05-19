package prattle_test

import (
	"fmt"
	"strconv"
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

type testDriver struct {
	stack  []int
	idents map[string]int
}

func (d *testDriver) pop() (v int) {
	if n := len(d.stack); n > 0 {
		v, d.stack = d.stack[n-1], d.stack[:n-1]
	}
	return v
}

func (d *testDriver) push(v int) {
	d.stack = append(d.stack, v)
}

func (d *testDriver) Precedence(kind prattle.Kind) prattle.Precedence {
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

func (d *testDriver) ident(p *prattle.Parser, t prattle.Token) error {
	v := d.idents[t.Text]
	d.push(v)
	return nil
}

func (d *testDriver) number(p *prattle.Parser, t prattle.Token) error {
	v, _ := strconv.Atoi(t.Text)
	d.push(v)
	return nil
}

func (d *testDriver) plus(p *prattle.Parser, t prattle.Token) error {
	if err := p.ParseExpression(d.Precedence(t.Kind)); err != nil {
		return err
	}
	right := d.pop()
	left := d.pop()
	d.push(right + left)
	return nil
}

func (d *testDriver) assign(p *prattle.Parser, t prattle.Token) error {
	if !p.Expect(kassign) {
		return d.ParseError(p.Peek())
	}

	if err := p.ParseExpression(d.Precedence(kassign)); err != nil {
		return err
	}

	d.idents[t.Text] = d.pop()

	if !p.Expect(ksemicolon) {
		return d.ParseError(p.Peek())
	}

	return nil
}

func (d *testDriver) Prefix(kind prattle.Kind) prattle.ParseFunc {
	switch kind {
	default:
		return nil
	case kident:
		return d.ident
	case knumber:
		return d.number
	}
}

func (d *testDriver) Infix(kind prattle.Kind) prattle.ParseFunc {
	switch kind {
	default:
		return nil
	case kplus:
		return d.plus
	}
}

func (d *testDriver) Statement(kind prattle.Kind) prattle.ParseFunc {
	switch kind {
	case kident:
		return d.assign
	default:
		return nil
	}
}

func (d *testDriver) ParseError(t prattle.Token) error {
	return fmt.Errorf("%s: unexpected '%s'", t.Position, t.Text)
}

// This example interprets a simple programming language that only has assignment statements and addition.
func Example_interpreter() {
	var c testDriver
	c.idents = map[string]int{}

	source := "a = 1;\nb = 2;\nc = a+b+b+a;\n"

	s := prattle.Scanner{Scan: testScan}
	p := prattle.Parser{Driver: &c}
	p.Init(s.Init(source))

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
