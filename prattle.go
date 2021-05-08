/*
Package prattle implements lexical scanning and Pratt parsing algorithms for
parsing programming or other structured text languages.
*/
package prattle

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
	// It is expected to return an appropriate error depending on the token name.
	// A token kind of 0 means unexpected end of input.
	// A token kind of < 0 means that an invalid token was encountered.
	// A token kind of > 0 means that an otherwise valid token was unexpected at this point.
	ParseError(Token) error
}
