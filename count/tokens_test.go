package count

import (
	"testing"
)

func TestTokens(t *testing.T) {
	var expected int
	expected = 1
	tm := make(TokensMap)
	tm.Add("one")
	if tm["one"] != expected {
		t.Errorf("Expected: %d, got: %d", expected, tm["one"])
	}
	expected = 2
	tm.Add("a.test")
	tm.Add("a.test")
	if tm["a.test"] != expected {
		t.Errorf("Expected %d, got %d", expected, tm["a.test"])
	}
	expected = 0
	if tm["nonExistant"] != expected {
		t.Errorf("Expected %d, got %d", expected, tm["nonExistant"])
	}
}
