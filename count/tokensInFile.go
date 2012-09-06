package count

import (
	"bufio"
	"io"
	"os"
	"regexp"
)

func TokensInFile(path string, tokenRE string, tokensMap TokensMap) {
	contents, err := os.Open(path)
	defer contents.Close()
	if err != nil {
		panic(err)
	}
	buf := bufio.NewReader(contents)
	r, err := regexp.Compile(tokenRE)
	if err != nil {
		panic(err)
	}
	for {
		read, err := buf.ReadString('\n')
		tokens := r.FindAllString(read, -1)
		for i := 0; i < len(tokens); i++ {
			tokensMap[tokens[i]] = tokensMap[tokens[i]] + 1
		}
		if err == io.EOF {
			break
		}
	}
}
