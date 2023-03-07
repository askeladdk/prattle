package prattle

import (
	"errors"
	"fmt"
	"testing"
)

type testDriver struct {
	infix  func(*Parser, Token) error
	prefix func(*Parser, Token) error
}

func (d testDriver) Infix(int) ParseFunc      { return d.infix }
func (d testDriver) Prefix(int) ParseFunc     { return d.prefix }
func (testDriver) Precedence(int) int         { return 1 }
func (d testDriver) ParseError(t Token) error { return fmt.Errorf("kind: %d", t.Kind) }

type tokeniter []Token

func (it *tokeniter) Next() (tok Token, ok bool) {
	if len(*it) != 0 {
		tok, *it = (*it)[0], (*it)[1:]
		return tok, true
	}
	return
}

func requireError(t testing.TB, err error) {
	if err == nil {
		t.Helper()
		t.Error("expected error")
	}
}

func requireNoError(t testing.TB, err error) {
	if err != nil {
		t.Helper()
		t.Error("unexpected error", err)
	}
}

func TestPrefixErrors(t *testing.T) {
	tokens := []Token{{Kind: 1}}

	t.Run("one", func(t *testing.T) {
		p := Parser{Driver: &testDriver{}}
		it := tokeniter(tokens)
		p.Init(&it)
		requireError(t, p.Parse(0))
	})

	t.Run("two", func(t *testing.T) {
		p := Parser{
			Driver: &testDriver{
				prefix: func(p *Parser, t Token) error { return errors.New("") },
			},
		}
		it := tokeniter(tokens)
		p.Init(&it)
		requireError(t, p.Parse(0))
	})
}

func TestInfixErrors(t *testing.T) {
	tokens := []Token{{Kind: 1}}

	t.Run("one", func(t *testing.T) {
		p := Parser{Driver: &testDriver{
			prefix: func(p *Parser, t Token) error { return nil },
		}}
		it := tokeniter(tokens)
		p.Init(&it)
		requireError(t, p.Parse(0))
	})

	t.Run("two", func(t *testing.T) {
		p := Parser{
			Driver: &testDriver{
				prefix: func(p *Parser, t Token) error { return nil },
				infix:  func(p *Parser, t Token) error { return errors.New("") },
			},
		}
		it := tokeniter(tokens)
		p.Init(&it)
		requireError(t, p.Parse(0))
	})

	t.Run("nonassoc", func(t *testing.T) {
		p := Parser{
			Driver: &testDriver{
				prefix: func(p *Parser, t Token) error { return nil },
				infix:  func(p *Parser, t Token) error { return NonAssoc },
			},
		}
		it := tokeniter(tokens)
		p.Init(&it)
		requireNoError(t, p.Parse(0))
	})
}

func TestParserExpect(t *testing.T) {
	tokens := []Token{{Kind: 1}}
	p := Parser{}
	it := tokeniter(tokens)
	p.Init(&it)
	if p.Expect(3) {
		t.Error("expected false")
	}
}

func TestErrorString(t *testing.T) {
	if NonAssoc.Error() == "" {
		t.Fatal()
	}
}
