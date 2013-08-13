package main

import (
	"flag"
	"fmt"
	"github.com/rudenoise/counting/count"
	"github.com/rudenoise/counting/dir"
	"path/filepath"
	"regexp"
)

type File struct {
	publish  []string
	subcribe []string
}

// set up flag defraults
var exclude = flag.String("e", "^$", "regexp pattern to exclude in file path")
var include = flag.String("i", "", "regexp pattern to include file path")

// set up pubSub RegExps
var pubRE = "m\\.publish\\([\\'\"]([\\w\\.]|\\w)+"
var subRE = "m\\.subscribe\\([\\'\"]([\\w\\.]|\\w)+"
var tokenRE = "([\\w\\.]|\\w)+$"
var fileNameRE = "[\\w\\.]+$"

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
	compiledFileNameRE, err := regexp.Compile(fileNameRE)
	if err != nil {
		panic(err)
	}
	// create a difinitive map of tokens
	allTokens := make(map[string]int)
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
			// loop all publishes
			for token := range pubMap {
				token = compiledTokenRE.FindString(token)
				file.publish = append(file.publish, token)
				// add to deduped list
				allTokens[token] = 1
			}
			// loop all subscribes
			for token := range subMap {
				token = compiledTokenRE.FindString(token)
				file.subcribe = append(file.subcribe, token)
				// add to deduped list
				allTokens[token] = 1
			}
			files[filePaths[i]] = File{file.publish, file.subcribe}
		}
	}
	// now we have a map of all files and their publshed and subscribed keys
	fmt.Println("digraph PubSub{")
	// create dot file output
	// all token nodes
	for tkn := range allTokens {
		fmt.Printf("\t\"%s\" [shape=circle]", tkn)
	}
	// all file nodes:
	for fn, rel := range files {
		baseFn := compiledFileNameRE.FindString(fn)
		fmt.Printf("\t\"%s\" [shape=box];\n", baseFn)
		// all publish relationships
		for i := 0; i < len(rel.publish); i++ {
			fmt.Printf("\t\"%s\"->\"%s\" [color=blue];", baseFn, rel.publish[i])
		}
		// all subscribe relationships
		for i := 0; i < len(rel.subcribe); i++ {
			fmt.Printf("\t\"%s\"->\"%s\" [color=orange];", rel.subcribe[i], baseFn)
		}
	}
	fmt.Println("}")
}
