package main

import (
	"flag"
	"fmt"
	"github.com/rudenoise/counting/count"
	"github.com/rudenoise/counting/dir"
	"path/filepath"
	"regexp"
	//"runtime"
	//"sort"
)

type File struct {
	publish []string
	subcribe []string
}

// set up flag defraults
var exclude = flag.String("e", "^$", "regexp pattern to exclude in file path")
var include = flag.String("i", "", "regexp pattern to include file path")
// set up pubSub RegExps
var pubRE = "m\\.publish\\([\\'\"]([\\w\\.]|\\w)+"
var subRE = "m\\.subscribe\\([\\'\"]([\\w\\.]|\\w)+"
var tokenRE = "([\\w\\.]|\\w)+$"

func main() {
	// parse cmd flags
	flag.Parse()
	// get the directory from which to start searching
	dirStr := filepath.Dir(flag.Arg(0))
	// traverse all child directories and collect a slice of file paths
	filePaths, err := dir.AllPaths(dirStr, *exclude, *include)
	if err != nil {
		panic(err)
	}
	// create the regexp to isolate the key
	compiledTokenRE, err := regexp.Compile(tokenRE)
	if err != nil {
		panic(err)
	}
	// create a map for files
	files := make(map[string]File)
	// loop all files paths
	length := len(filePaths)
	for i := 0; i < length; i++ {
		// collect published tokens in file
		pubMap := make(count.TokensMap)
		count.TokensInFile(filePaths[i], pubRE, pubMap)
		// collect subscribed tokens in file
		subMap := make(count.TokensMap)
		count.TokensInFile(filePaths[i], subRE, subMap)
		if len(pubMap) > 0 || len(subMap) > 0 {
			file := new(File)
			fmt.Printf("%s\n", filePaths[i])
			for token := range pubMap {
				token = compiledTokenRE.FindString(token)
				file.publish = append(file.publish, token)
			}
			for token := range subMap {
				token = compiledTokenRE.FindString(token)
				file.subcribe = append(file.subcribe, token)
			}
			files[filePaths[i]] = File{file.publish, file.subcribe}
		}
	}
	// now we have a map of all files and their publshed and subscribed keys
	fmt.Println(files)
}

