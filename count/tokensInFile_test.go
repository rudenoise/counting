package count

import (
	"testing"
)

func TestTokensInFile(t *testing.T) {
	testFilePath := "testTokens1.txt"
	tm := make(TokensMap)
	TokensInFile(testFilePath, "[a-zA-Z]+", tm)
	expectedHello := 3
	expectedWorld := 1
	if tm["hello"] != expectedHello && tm["world"] != expectedWorld {
		t.Errorf("Expected %d occurances of 'hello', got %d", expectedHello, tm["hello"])
		t.Errorf("Expected %d occurances of 'world', got %d", expectedWorld, tm["world"])
	}
}

func TestTokensInFiles(t *testing.T) {
	testFilePaths := []string{"testTokens1.txt", "testTokens2.txt"}
	tm := make(TokensMap)
	TokensInFiles(testFilePaths, "[a-zA-z]+", tm)
	expectedHello := 6
	expectedWorld := 2
	if tm["hello"] != expectedHello && tm["world"] != expectedWorld {
		t.Errorf("Expected %d occurances of 'hello', got %d", expectedHello, tm["hello"])
		t.Errorf("Expected %d occurances of 'world', got %d", expectedWorld, tm["world"])
	}
}
