package count

import (
	"testing"
)

func TestLinesInFile(t *testing.T) {
	const expectedTotal = 3
	actualTotal, err := LinesInFile("tstFile.txt")
	if err != nil {
		t.Error(err)
	}
	if expectedTotal != actualTotal {
		t.Errorf("Expected: %d, actual: %d", expectedTotal, actualTotal)
	}
}
