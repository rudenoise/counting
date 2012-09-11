package count

import (
	"bufio"
	"io"
	"os"
	"regexp"
)

func LinesInFile(filePath string, ignoreCommentsAndEmptyLines bool) (int, error) {
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
		if ignoreCommentsAndEmptyLines == true {
			ignoreLineMatch, err := regexp.MatchString("^\n$|^[\\ \t]+\n$|^[\\ \t]+\\/\\/.*\n$|^\\/\\/.*\n$", read)
			if err != nil {
				return total, err
			}
			if ignoreLineMatch == false {
				total++
			}
		} else {
			total++
		}
	}
	return total, nil
}
