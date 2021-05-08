package prattle

import (
	"errors"
)

// Parser implements the Pratt parsing algorithm,
// also known as the top down operator precedence (TDOP) algorithm.
// This a recursive descent algorithm that is able to handle operator precedence
// in a simple and flexible manner.
//
// Parser operates on a Context and Sequence.
// Context drives the algorithm by providing the precedence values and parsing functions.
// Sequence provides the token stream.
type Parser struct {
	Sequence
	Context
	token Token
}

// ErrNonAssoc is returned by infix ParseFuncs to indicate that an operator is non-associative.
var ErrNonAssoc = errors.New("non-associative operator")

// // NewParser creates a new parser.
// func NewParser(ctx Context, seq Sequence) *Parser {
// 	p := &Parser{
// 		ctx: ctx,
// 		seq: seq,
// 	}
// 	_ = p.Next()
// 	return p
// }

// Peek returns the last read token.
func (p *Parser) Peek() Token {
	return p.token
}

// Next advances to the next token.
func (p *Parser) Next() Token {
	p.token = p.Sequence.Next()
	return p.token
}

// Expect advances if the current token is equal to name and returns true if so.
func (p *Parser) Expect(name Kind) bool {
	if p.token.Kind != name {
		return false
	}
	_ = p.Next()
	return true
}

// ParseExpression parses until a token with an equal or lower precedence is encountered.
// It is called in a mutual recursive manner by the parsing functions provided by Context.
// Pass zero precedence to kick off parsing until the end of input.
func (p *Parser) ParseExpression(least Precedence) error {
	t := p.Peek()
	_ = p.Next()

	if prefix := p.Prefix(t.Kind); prefix == nil {
		return p.ParseError(t)
	} else if err := prefix(p, t); err != nil {
		return err
	}

	for t = p.Peek(); least < p.Precedence(t.Kind); t = p.Peek() {
		_ = p.Next()

		if infix := p.Infix(t.Kind); infix == nil {
			return p.ParseError(t)
		} else if err := infix(p, t); err == ErrNonAssoc {
			least = p.Precedence(t.Kind) + 1
		} else if err != nil {
			return err
		}
	}

	return nil
}

// ParseStatement parses one statement.
func (p *Parser) ParseStatement() error {
	t := p.Peek()
	_ = p.Next()

	stmt := p.Statement(t.Kind)
	if stmt == nil {
		return p.ParseError(t)
	}

	return stmt(p, t)
}

// ParseStatements parses zero or more statements.
func (p *Parser) ParseStatements(accept func(Kind) bool) error {
	for t := p.Peek(); accept(t.Kind); t = p.Peek() {
		_ = p.Next()

		if stmt := p.Statement(t.Kind); stmt == nil {
			return p.ParseError(t)
		} else if err := stmt(p, t); err != nil {
			return err
		}
	}

	return nil
}
