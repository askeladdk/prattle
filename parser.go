package prattle

import (
	"errors"
)

// Sequence is the interface that wraps a sequence of tokens and is the input to Parser.
type Sequence interface {
	// Next returns the next token in the sequence.
	// It must return the zero token when the sequence is depleted.
	Next() Token
}

// ParseFunc parses an expression.
type ParseFunc func(*Parser, Token) error

// Precedence is the token precedence.
// The higher the precedence, the tighter the token binds.
// The lowest precedence is zero, which is reserved for the final token that signals end of input.
type Precedence int

// Context is passed to Parser and drives the parsing algorithm.
type Context interface {
	// Infix returns the ParseFunc associated with an infix operator.
	Infix(Kind) ParseFunc

	// Infix returns the ParseFunc associated with a prefix operator.
	Prefix(Kind) ParseFunc

	// Statement returns the ParseFunc associated with a statement.
	Statement(Kind) ParseFunc

	// Precedence associates a token with an operator precedence.
	// The lowest precedence is zero and is reserved for the EndOf token.
	Precedence(Kind) Precedence

	// ParseError is called by the Parser when it encounters a token it cannot parse.
	ParseError(Token) error
}

// Parser implements the Pratt parsing algorithm,
// also known as the top down operator precedence (TDOP) algorithm.
// This a recursive descent algorithm that is able to handle operator precedence
// in a simple and flexible manner.
//
// Parser operates on a Context and a Sequence.
// Context drives the algorithm by providing the precedence values and parsing functions.
// Sequence provides the token stream.
type Parser struct {
	// Context drives the Parser algorithm.
	Context

	sequence Sequence
	token    Token
}

// ErrNonAssoc is returned by infix ParseFuncs to indicate that an operator is non-associative.
var ErrNonAssoc = errors.New("non-associative operator")

// Init initializes the Parser with a Sequence and returns it.
func (p *Parser) Init(sequence Sequence) *Parser {
	p.sequence = sequence
	p.Advance()
	return p
}

// Peek returns the last read token.
func (p *Parser) Peek() Token {
	return p.token
}

// Advance reads the next token.
func (p *Parser) Advance() {
	p.token = p.sequence.Next()
}

// Expect advances to the next token if the current token kind matches.
func (p *Parser) Expect(kind Kind) bool {
	if p.token.Kind != kind {
		return false
	}
	p.Advance()
	return true
}

// ParseExpression parses until a token with an equal or lower precedence is encountered.
// It is called in a mutual recursive manner by the parsing functions provided by Context.
func (p *Parser) ParseExpression(least Precedence) error {
	t := p.Peek()
	p.Advance()

	if prefix := p.Prefix(t.Kind); prefix == nil {
		return p.ParseError(t)
	} else if err := prefix(p, t); err != nil {
		return err
	}

	for t = p.Peek(); least < p.Precedence(t.Kind); t = p.Peek() {
		p.Advance()

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
	p.Advance()

	stmt := p.Statement(t.Kind)
	if stmt == nil {
		return p.ParseError(t)
	}

	return stmt(p, t)
}

// ParseStatements parses zero or more statements.
func (p *Parser) ParseStatements(accept func(Kind) bool) error {
	for t := p.Peek(); accept(t.Kind); t = p.Peek() {
		p.Advance()

		if stmt := p.Statement(t.Kind); stmt == nil {
			return p.ParseError(t)
		} else if err := stmt(p, t); err != nil {
			return err
		}
	}

	return nil
}
