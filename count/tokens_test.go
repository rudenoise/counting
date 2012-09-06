package count

import (
	"testing"
)

func TestTokens(t *testing.T) {
	tm := make(TokensMap)
	tm.Add("one")
	if tm["one"] != 1 {
		t.Errorf("A failing test")
	}
	tm.Add("a.test")
	tm.Add("a.test")
	if tm["a.test"] != 2 {
		t.Errorf("Expected 2, got %d", tm["a.test"])
	}
}
