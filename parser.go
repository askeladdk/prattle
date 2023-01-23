package prattle

import (
	"errors"
)

// ErrNonAssoc is returned by infix ParseFuncs to indicate that an operator is non-associative.
var ErrNonAssoc = errors.New("non-associative operator")

// Sequence represents a stream of tokens.
type Sequence interface {
	// NextToken returns the next token in the stream.
	NextToken() Token
}

// ParseFunc parses an expression or statement.
type ParseFunc func(*Parser, Token) error

// Driver drives the parsing algorithm by associating tokens to parser functions.
// It is expected to hold the parse state and results, such as the syntax tree.
type Driver interface {
	// Infix associates an infix ParseFunc with a token.
	// Returning nil is a parse error.
	Infix(kind int) ParseFunc

	// Prefix associates a prefix ParseFunc with a token.
	// Returning nil is a parse error.
	Prefix(kind int) ParseFunc

	// Statement associates a statement ParseFunc with a token.
	// Returning nil is a parse error.
	Statement(kind int) ParseFunc

	// Precedence associates an operator precedence with a token.
	// The higher the precedence, the tighter the token binds.
	Precedence(kind int) (precedence int)

	// ParseError is called by the Parser when it encounters a token that it cannot parse.
	ParseError(Token) error
}

// Parser implements the Pratt parsing algorithm,
// also known as the top down operator precedence (TDOP) algorithm.
// This a recursive descent algorithm that is able to handle operator precedence
// in a simple and flexible manner.
//
// Parser consumes tokens from a Sequence and uses a Driver
// to determine precedence and executing parsing functions.
type Parser struct {
	// Driver drives the Parser.
	Driver

	sequence Sequence
	token    Token
}

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

// Advance reads the next token from the Sequence.
func (p *Parser) Advance() {
	p.token = p.sequence.NextToken()
}

// Expect advances to the next token if the current token kind matches.
func (p *Parser) Expect(kind int) bool {
	if p.token.Kind != kind {
		return false
	}
	p.Advance()
	return true
}

// ParseExpression parses until a token with an equal or lower precedence than least is encountered.
// It is called in a mutual recursive manner by the parsing functions provided by the Driver.
func (p *Parser) ParseExpression(least int) error {
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

// ParseStatements parses zero or more statements until a token
// with an equal or lower precedence than least is encountered.
func (p *Parser) ParseStatements(least int) error {
	for t := p.Peek(); least < p.Precedence(t.Kind); t = p.Peek() {
		p.Advance()

		if stmt := p.Statement(t.Kind); stmt == nil {
			return p.ParseError(t)
		} else if err := stmt(p, t); err != nil {
			return err
		}
	}

	return nil
}
