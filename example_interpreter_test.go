package prattle_test

import (
	"fmt"
	"strconv"
	"unicode"

	"github.com/askeladdk/prattle"
)

const (
	ksemicolon = 1 + iota
	kassign
	kplus
	kident
	knumber
)

func testScan(s *prattle.Scanner) int {
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
	stack  []prattle.Token
	idents map[string]int
}

func (d *testDriver) pop() (v prattle.Token) {
	if n := len(d.stack); n > 0 {
		v, d.stack = d.stack[n-1], d.stack[:n-1]
	}
	return v
}

func (d *testDriver) push(v prattle.Token) {
	d.stack = append(d.stack, v)
}

func (d *testDriver) Precedence(kind int) int {
	return kind
}

func (d *testDriver) tonumber(t prattle.Token) int {
	if t.Kind == kident {
		return d.idents[t.Text]
	}
	num, _ := strconv.Atoi(t.Text)
	return num
}

func (d *testDriver) plus(p *prattle.Parser, t prattle.Token) error {
	if err := p.Parse(d.Precedence(t.Kind)); err != nil {
		return err
	}
	right := d.tonumber(d.pop())
	left := d.tonumber(d.pop())
	d.push(prattle.Token{
		Kind: knumber,
		Text: strconv.Itoa(left + right),
	})
	return nil
}

func (d *testDriver) assign(p *prattle.Parser, t prattle.Token) error {
	if err := p.Parse(d.Precedence(kassign)); err != nil {
		return err
	}

	right := d.pop()
	left := d.pop()
	d.idents[left.Text] = d.tonumber(right)
	return nil
}

func (d *testDriver) primitive(p *prattle.Parser, t prattle.Token) error {
	d.push(t)
	return nil
}

func (d *testDriver) Prefix(kind int) prattle.ParseFunc {
	switch kind {
	default:
		return nil
	case kident, knumber:
		return d.primitive
	}
}

func (d *testDriver) Infix(kind int) prattle.ParseFunc {
	switch kind {
	default:
		return nil
	case kplus:
		return d.plus
	case kassign:
		return d.assign
	}
}

func (d *testDriver) ParseError(t prattle.Token) error {
	return fmt.Errorf("%s: unexpected '%s'", t.Position, t.Text)
}

// This example demonstrates parsing a simple programming language that consists of a sequence of statements.
func Example_interpreter() {
	c := testDriver{
		idents: make(map[string]int),
	}

	source := "a = 1;\nb = 2;\nc = a+b+b+a;\n"

	s := prattle.Scanner{Scan: testScan}
	p := prattle.Parser{Driver: &c}
	p.Init(s.InitWithString(source))

	// Parse expressions separated by semicolons.
	for p.Peek().Kind != 0 {
		if err := p.Parse(ksemicolon); err != nil {
			fmt.Println(err)
			return
		} else if !p.Expect(ksemicolon) {
			fmt.Println("expected semicolon")
			return
		}
	}

	fmt.Printf("c = %d\n", c.idents["c"])

	// Output:
	// c = 6
}
