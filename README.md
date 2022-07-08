# Prattle

[![GoDoc](https://godoc.org/github.com/askeladdk/prattle?status.png)](https://godoc.org/github.com/askeladdk/prattle)
[![Go Report Card](https://goreportcard.com/badge/github.com/askeladdk/prattle)](https://goreportcard.com/report/github.com/askeladdk/prattle)

## Overview

Package prattle implements a general purpose, unicode-aware lexical scanner and top down operator precedence parser suitable for parsing LL(1) grammars. The scanner and parser can be used independently from each other if desired.

## Install

```
go get -u github.com/askeladdk/prattle
```

## Quickstart

Use `Scanner` to produce a sequence of `Token`s by scanning a source text using a `ScanFunc`. Use the `Expect*`, `Skip` and `Advance` methods to scan tokens.

```go
// Prepare the scanner.
scanner := prattle.Scanner{
	// Scan scans the next token and returns its kind.
	Scan: func(s *prattle.Scanner) (kind int) {
		// Skip any whitespace.
		s.ExpectAny(unicode.IsSpace)
		s.Skip()

		// Scan the next token.
		switch {
		case s.Done(): // Stop when the entire input has been consumed.
			return 0
		case s.Expect('+'): // Scan the addition operator.
			return 1
		case s.ExpectOne(unicode.IsDigit): // Scan a number consisting of one or more digits.
			s.ExpectAny(unicode.IsDigit)
			return 2
		}

		// Invalid token.
		s.Advance()
		return -1
	},
}
```

Use `Parser` and `Driver` to associate tokens produced by `Scanner` with `ParseFunc`s. Define the `Driver` first.

```go
// Define the parsing Driver.
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
	// Parse the right hand operator.
	_ = p.ParseExpression(d.Precedence(t.Kind))

	right := d.pop()
	left := d.pop()
	sum := left + right
	fmt.Printf("%d + %d = %d\n", left, right, sum)
	d.push(sum)
	return nil
}

func (d *driver) Prefix(kind int) prattle.ParseFunc {
	return d.number
}

func (d *driver) Infix(kind int) prattle.ParseFunc {
	return d.add
}

func (d *driver) Statement(kind int) prattle.ParseFunc {
	return nil
}

func (d *driver) Precedence(kind int) int {
	return kind
}

func (d *driver) ParseError(t prattle.Token) error {
	return fmt.Errorf("%s", t)
}

// Prepare the parser.
parser := prattle.Parser{
	Driver: &driver{},
}
```

Finally, `Init` the scanner and parser, and parse an expression.

```go
source := "1 + 23 + 456 + 7890"
scanner.InitWithString(source)
parser.Init(&scanner)
_ = parser.ParseExpression(0)

// Output:
// 1 + 23 = 24
// 24 + 456 = 480
// 480 + 7890 = 8370
```

See the [calculator](_examples/calculator/main.go) for a complete example featuring prefix, postfix, infix, left, right and non-associative operators. Also read the [documentation on pkg.go.dev](https://pkg.go.dev/github.com/askeladdk/prattle).

## License

Package prattle is released under the terms of the ISC license.
