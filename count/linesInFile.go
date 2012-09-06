package count

import (
	"bufio"
	"io"
	"os"
)

func LinesInFile(filePath string) (int, error) {
	total := 0
	contents, err := os.Open(filePath)
	defer contents.Close()
	if err != nil {
		return total, err
	}
	buf := bufio.NewReader(contents)
	for {
		read, err := buf.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return total, err
		}
		if read != "" {
			total++
		}
	}
	return total, nil
}
