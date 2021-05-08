package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"unicode"

	"github.com/askeladdk/prattle"
)

const (
	number prattle.Kind = 1 + iota
	plus
	minus
	star
	slash
	caret
	modulo
	leftPar
	rightPar
	pi
	bang
	squareRoot
	answer
)

var rune2name = map[rune]prattle.Kind{
	'+': plus,
	'-': minus,
	'*': star,
	'/': slash,
	'^': caret,
	'ˆ': caret,
	'%': modulo,
	'(': leftPar,
	')': rightPar,
	'π': pi,
	'!': bang,
	'√': squareRoot,
}

func scan(s *prattle.Scanner) prattle.ScanFunc {
	s.Any(unicode.IsSpace)
	s.Skip()

	r := s.Peek()
	switch {
	case s.Done():
		return nil
	case s.OneRune('a'):
		if s.OneRune('n') {
			if s.OneRune('s') {
				s.Emit(answer)
				return scan
			}
		}
	case s.One(unicode.IsDigit):
		s.Any(unicode.IsDigit)
		if s.OneRune('.') {
			s.Any(unicode.IsDigit)
		}
		s.Emit(number)
		return scan
	case s.OneRune('.'):
		if s.One(unicode.IsDigit) {
			s.Any(unicode.IsDigit)
			s.Emit(number)
			return scan
		}
	case s.One(prattle.OneOf("+-*/^ˆ%()π!√")):
		s.Emit(rune2name[r])
		return scan
	}

	_ = s.Advance()
	s.Emit(-1)
	return nil
}

type calculator struct {
	stack  []float64
	answer float64
}

func (c *calculator) pop() (v float64) {
	n := len(c.stack)
	if n == 0 {
		return math.NaN()
	}
	v, c.stack = c.stack[n-1], c.stack[:n-1]
	return
}

func (c *calculator) push(v float64) {
	c.stack = append(c.stack, v)
}

func (c *calculator) number(p *prattle.Parser, t prattle.Token) error {
	v, err := strconv.ParseFloat(t.Text, 64)
	if err != nil {
		return err
	}
	c.push(v)
	return nil
}

func (c *calculator) pi(p *prattle.Parser, t prattle.Token) error {
	c.push(math.Pi)
	return nil
}

func (c *calculator) ans(p *prattle.Parser, t prattle.Token) error {
	c.push(c.answer)
	return nil
}

func (c *calculator) binop(kind prattle.Kind) {
	right := c.pop()
	left := c.pop()

	var result float64

	switch kind {
	case plus:
		result = left + right
	case minus:
		result = left - right
	case star:
		result = left * right
	case slash:
		result = left / right
	case caret:
		result = math.Pow(left, right)
	case modulo:
		result = math.Mod(left, right)
	case squareRoot:
		result = math.Pow(right, 1/left)
	default:
		result = math.NaN()
	}

	c.push(result)
}

func (c *calculator) binopLeftAssoc(p *prattle.Parser, t prattle.Token) error {
	if err := p.ParseExpression(c.Precedence(t.Kind)); err != nil {
		return err
	}
	c.binop(t.Kind)
	return nil
}

func (c *calculator) binopRightAssoc(p *prattle.Parser, t prattle.Token) error {
	if err := p.ParseExpression(c.Precedence(t.Kind) - 1); err != nil {
		return err
	}
	c.binop(t.Kind)
	return nil
}

func (c *calculator) binopNonAssoc(p *prattle.Parser, t prattle.Token) error {
	if err := p.ParseExpression(c.Precedence(t.Kind)); err != nil {
		return err
	}
	c.binop(t.Kind)
	return prattle.ErrNonAssoc
}

func (c *calculator) unary(p *prattle.Parser, t prattle.Token) error {
	if err := p.ParseExpression(c.Precedence(t.Kind)); err != nil {
		return err
	}

	v := c.pop()

	switch t.Kind {
	case minus:
		c.push(-v)
	case squareRoot:
		c.push(math.Sqrt(v))
	default:
		return c.ParseError(t)
	}
	return nil
}

func (c *calculator) factorial(p *prattle.Parser, t prattle.Token) error {
	v := c.pop()
	i, f := math.Modf(v)
	if f != 0 {
		return fmt.Errorf("cannot compute factorial of fractional number '%f'", v)
	}
	a, b := float64(1), int(i)
	for ; b > 0; b-- {
		a *= float64(b)
	}
	c.push(a)
	return nil
}

func (c *calculator) paren(p *prattle.Parser, t prattle.Token) error {
	if err := p.ParseExpression(1); err != nil {
		return err
	} else if !p.Expect(rightPar) {
		return c.ParseError(t)
	}
	return nil
}

func (c *calculator) Prefix(kind prattle.Kind) prattle.ParseFunc {
	switch kind {
	case number:
		return c.number
	case pi:
		return c.pi
	case answer:
		return c.ans
	case minus, squareRoot:
		return c.unary
	case leftPar:
		return c.paren
	default:
		return nil
	}
}

func (c *calculator) Infix(kind prattle.Kind) prattle.ParseFunc {
	switch kind {
	case plus, minus, star, slash, modulo:
		return c.binopLeftAssoc
	case caret:
		return c.binopRightAssoc
	case bang:
		return c.factorial
	case squareRoot:
		return c.binopNonAssoc
	default:
		return nil
	}
}

func (c *calculator) Statement(kind prattle.Kind) prattle.ParseFunc {
	return nil
}

func (c *calculator) Precedence(kind prattle.Kind) prattle.Precedence {
	switch kind {
	default:
		return 0
	case leftPar, rightPar:
		return 1
	case squareRoot:
		return 2
	case plus, minus:
		return 3
	case star, slash, modulo:
		return 4
	case caret:
		return 5
	case number, pi:
		return 6
	case bang:
		return 7
	}
}

func (c *calculator) ParseError(t prattle.Token) error {
	if t.Kind == 0 {
		return fmt.Errorf("incomplete equation")
	}
	return fmt.Errorf("i do not understand '%s'", t.Text)
}

func (c *calculator) calculate(expr string) (v float64, err error) {
	s := prattle.NewScanner("<calc>", expr, scan)
	defer s.Flush()
	p := prattle.Parser{
		Sequence: s,
		Context:  c,
	}
	p.Next()
	err = p.ParseExpression(0)
	if err == nil && p.Peek().Kind != 0 {
		err = p.ParseError(p.Peek())
	}
	v = c.pop()
	return
}

func main() {
	fmt.Println("welcome to calculator")
	fmt.Println("enter an equation or q to quit")
	fmt.Println("enter π for pi, √ for square root")

	var calc calculator
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf("> ")
		scanner.Scan()
		text := scanner.Text()

		if text == "" {
			continue
		} else if text == "q" {
			break
		} else if v, err := calc.calculate(text); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(v)
			calc.answer = v
		}
	}
}
