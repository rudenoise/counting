package count

import (
	"testing"
)

func TestLinesInFileWithoutCommentsAndEmptyLines(t *testing.T) {
	expectedTotal := 4
	actualTotal, err := LinesInFile("tstFile.txt", true)
	if err != nil {
		t.Error(err)
	}
	if expectedTotal != actualTotal {
		t.Errorf("Expected: %d, actual: %d", expectedTotal, actualTotal)
	}
}
func TestLinesInFileWithCommentsAndEmptyLines(t *testing.T) {
	expectedTotal := 9
	actualTotal, err := LinesInFile("tstFile.txt", false)
	if err != nil {
		t.Error(err)
	}
	if expectedTotal != actualTotal {
		t.Errorf("Expected: %d, actual: %d", expectedTotal, actualTotal)
	}
}
