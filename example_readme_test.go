package prattle_test

import (
	"fmt"
	"strconv"
	"unicode"

	"github.com/askeladdk/prattle"
)

type driver struct {
	stack []int
}

func (d *driver) push(i int) {
	d.stack = append(d.stack, i)
}

func (d *driver) pop() (i int) {
	n := len(d.stack)
	i, d.stack = d.stack[n-1], d.stack[:n-1]
	return
}

func (d *driver) number(p *prattle.Parser, t prattle.Token) error {
	n, _ := strconv.Atoi(t.Text)
	d.push(n)
	return nil
}

func (d *driver) add(p *prattle.Parser, t prattle.Token) error {
	// First parse the right hand operator.
	if err := p.Parse(d.Precedence(t.Kind)); err != nil {
		return err
	}

	right := d.pop()
	left := d.pop()
	acc := left + right
	fmt.Printf("%d + %d = %d\n", left, right, acc)
	d.push(acc)
	return nil
}

func (d *driver) Prefix(k int) prattle.ParseFunc {
	if k == 2 {
		return d.number
	}
	return nil
}

func (d *driver) Infix(k int) prattle.ParseFunc {
	if k == 1 {
		return d.add
	}
	return nil
}

func (d *driver) Precedence(k int) int {
	return k
}

func (d *driver) ParseError(t prattle.Token) error {
	return fmt.Errorf("%s", t)
}

// This example is described in the readme.
func Example_readme() {
	scanner := prattle.Scanner{
		Scan: func(s *prattle.Scanner) int {
			// Skip whitespaces.
			s.ExpectAny(unicode.IsSpace)
			s.Skip()
			switch {
			case s.Done(): // Stop if the entire input has been consumed.
				return 0
			case s.Expect('+'):
				return 1
			case s.ExpectOne(unicode.IsDigit): // Read a number consisting of one or more digits.
				s.ExpectAny(unicode.IsDigit)
				return 2
			}

			s.Advance()
			return -1
		},
	}

	parser := prattle.Parser{
		Driver: &driver{},
	}

	source := "1 + 23 + 456 + 7890"
	scanner.InitWithString(source)
	if err := parser.Init(&scanner).Parse(0); err != nil {
		fmt.Println(err)
	}

	// Output:
	// 1 + 23 = 24
	// 24 + 456 = 480
	// 480 + 7890 = 8370
}
