package count

import (
	"bufio"
	"io"
	"os"
	"regexp"
	"sync"
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

func TokensInFiles(filePaths []string, tokenRE string, tokensMap TokensMap) {
	var wg sync.WaitGroup
	for i := 0; i < len(filePaths); i++ {
		TokensInFile(filePaths[i], tokenRE, tokensMap)
		/*
			wg.Add(1)
			go func(fp string) {
				TokensInFile(fp, tokenRE, tokensMap)
				wg.Done()
			}(filePaths[i])
		*/
	}
	wg.Wait()
}
