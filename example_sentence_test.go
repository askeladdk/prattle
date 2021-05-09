package prattle_test

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/askeladdk/prattle"
)

// This example tokenizes a sentence into words and punctuation.
func Example_sentence() {
	scan := func(s *prattle.Scanner) prattle.Kind {
		// Skip any whitespace.
		s.ExpectAny(unicode.IsSpace)
		s.Skip()

		// Return an appropriate token kind.
		switch {
		case s.Done(): // EOF
			return 0
		case s.ExpectOne(unicode.IsLetter): // Word
			s.ExpectAny(unicode.IsLetter)
			return 1
		case s.Expect('!'): // Punctuation
			return 2
		}

		// Unrecognized character.
		s.Advance()
		return -1
	}

	sentence := "I love it when a plan comes together!"

	s := prattle.Scanner{Scan: scan}
	s.Init(strings.NewReader(sentence))

	for tok := s.Next(); tok.Kind != 0; tok = s.Next() {
		fmt.Printf("[%d] %s\n", tok.Kind, tok.Text)
	}

	// Output:
	// [1] I
	// [1] love
	// [1] it
	// [1] when
	// [1] a
	// [1] plan
	// [1] comes
	// [1] together
	// [2] !
}
