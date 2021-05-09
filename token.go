package prattle

import "fmt"

// Kind identifies the type of token.
// Zero is reserved to signal end-of-input.
type Kind int

// Position represents the position of a Token in the input string.
type Position struct {
	// Filename is the filename of the input.
	Filename string

	// Offset is the byte offset in the input.
	Offset int

	// Line is the line number.
	Line int

	// Column is the character from the line.
	Column int
}

// IsValid reports whether the position is valid.
func (p Position) IsValid() bool {
	return p.Line > 0 && p.Column > 0
}

func (p Position) String() string {
	filename := p.Filename
	if filename == "" {
		filename = "<input>"
	}
	if p.IsValid() {
		return fmt.Sprintf("%s:%d:%d", filename, p.Line, p.Column)
	}
	return filename
}

// Token is a fragment of tokenised text.
type Token struct {
	Position

	// Kind identifies the kind of token.
	Kind Kind

	// Text is the token value.
	Text string
}

func (t Token) String() string {
	return fmt.Sprintf("%s: (%d) '%s'", t.Position, t.Kind, t.Text)
}
