package prattle

import (
	"testing"
)

func TestParserPanic(t *testing.T) {
	defer func() {
		if s, _ := recover().(string); s == "" {
			t.Fatal()
		}
	}()
	(&Parser{}).Init(nil)
}
