package prattle

import "testing"

func TestPositionString(t *testing.T) {
	for _, testCase := range []struct {
		Name   string
		Pos    Position
		Expect string
	}{
		{
			Name:   "ZeroValue",
			Pos:    Position{},
			Expect: "<input>",
		},
		{
			Name: "LineOnly",
			Pos: Position{
				Filename: "hello.txt",
				Line:     13,
			},
			Expect: "hello.txt:13",
		},
		{
			Name: "LineAndColumn",
			Pos: Position{
				Filename: "hello.txt",
				Line:     13,
				Column:   37,
			},
			Expect: "hello.txt:13,37",
		},
	} {
		t.Run(testCase.Name, func(t *testing.T) {
			if testCase.Pos.String() != testCase.Expect {
				t.Error()
			}
		})
	}
}

func TestTokenString(t *testing.T) {
	tok := Token{
		Position: Position{
			Filename: "hello.txt",
			Line:     13,
			Column:   37,
		},
		Text: "123",
		Kind: 1,
	}

	if tok.String() != "hello.txt:13,37: '123'(1)" {
		t.Error()
	}
}
