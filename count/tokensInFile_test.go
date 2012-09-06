package count

import (
	"testing"
)

func TestTokensInFile(t *testing.T) {
	testFilePath := "testTokens.txt"
	tm := make(TokensMap)
	TokensInFile(testFilePath, "[a-zA-Z]+", tm)
	expectedHello := 3
	expectedWorld := 1
	if tm["hello"] != expectedHello && tm["world"] != expectedWorld {
		t.Errorf("Expected %d occurances of 'hello', got %d", expectedHello, tm["hello"])
		t.Errorf("Expected %d occurances of 'world', got %d", expectedWorld, tm["world"])
	}
}
