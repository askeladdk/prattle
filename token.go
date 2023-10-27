package prattle

import "fmt"

// Position represents the position of a Token in the input string.
type Position struct {
	// Filename is the filename of the input, if any.
	Filename string

	// Offset is the byte offset, starting at 0.
	Offset int

	// Line is the line number, starting at 1.
	Line int

	// Column is the column number, starting at 1 (character count per line).
	Column int
}

// IsValid reports whether a Position is valid.
// A Position is valid if its Line field is greater than zero.
func (p Position) IsValid() bool {
	return p.Line > 0
}

// String implements fmt.Stringer.
func (p Position) String() string {
	filename := p.Filename
	if filename == "" {
		filename = "<input>"
	}
	if p.IsValid() {
		if p.Column > 0 {
			return fmt.Sprintf("%s:%d,%d", filename, p.Line, p.Column)
		}
		return fmt.Sprintf("%s:%d", filename, p.Line)
	}
	return filename
}

// Token is a fragment of tokenised text.
type Token struct {
	Position

	// Kind identifies the kind of token.
	Kind int

	// Text is the token value.
	Text string
}

// String implements fmt.Stringer.
func (t Token) String() string {
	return fmt.Sprintf("%s: '%s'(%d)", t.Position, t.Text, t.Kind)
}
