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
	stateSet []string
	stateGet []string
}

// set up flag defraults
var exclude = flag.String("e", "^$", "regexp pattern to exclude in file path")
var include = flag.String("i", "", "regexp pattern to include file path")

// set up pubSub RegExps
var pubRE = "m\\.publish\\([\\'\"]([\\w\\.]|\\w)+"
var subRE = "m\\.subscribe\\([\\'\"]([\\w\\.]|\\w)+"
var stateSetRE = "m\\.state\\.set\\([\\'\"]([\\w\\.]|\\w)+"
var stateGetRE = "m\\.(state\\.get|state)\\([\\'\"]([\\w\\.]|\\w)+"
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
	// create a difinitive map of tokens
	allTokens := make(map[string]int)
	// create a map for files
	files := make(map[string]File)
	// loop all files paths
	length := len(filePaths)
	for i := 0; i < length; i++ {
		// collect published tokens from file
		pubMap := make(count.TokensMap)
		count.TokensInFile(filePaths[i], pubRE, pubMap)
		// collect subscribed tokens from file
		subMap := make(count.TokensMap)
		count.TokensInFile(filePaths[i], subRE, subMap)
		// collect state.set tokens from file
		stateSetMap := make(count.TokensMap)
		count.TokensInFile(filePaths[i], stateSetRE, stateSetMap)
		// collect state.get tokens from file
		stateGetMap := make(count.TokensMap)
		count.TokensInFile(filePaths[i], stateGetRE, stateGetMap)
		// make sure there is data to collect
		if len(pubMap) > 0 || len(subMap) > 0 || len(stateSetMap) > 0 || len(stateGetMap) > 0 {
			file := new(File)
			// loop all publishes
			for token := range pubMap {
				token = compiledTokenRE.FindString(token)
				file.publish = append(file.publish, token)
				// add to deduped list
				allTokens[token] = 0
			}
			// loop all subscribes
			for token := range subMap {
				token = compiledTokenRE.FindString(token)
				file.subcribe = append(file.subcribe, token)
				// add to deduped list
				allTokens[token] = 0
			}
			// loop all state set tokens
			for token := range stateSetMap {
				token = compiledTokenRE.FindString(token)
				file.stateSet = append(file.stateSet, token)
				// add to deduped list
				allTokens[token] = 1
			}
			// loop all state get tokens
			for token := range stateGetMap {
				token = compiledTokenRE.FindString(token)
				file.stateGet = append(file.stateGet, token)
				// add to deduped list
				allTokens[token] = 1
			}
			files[filePaths[i]] = File{file.publish, file.subcribe, file.stateSet, file.stateGet}
		}
	}
	// now we have a map of all files and their publshed and subscribed keys
	fmt.Printf("digraph PubSub{\n")
	// create dot file output
	// all token nodes
	for tkn, val := range allTokens {
		style := ""
		if val == 1 {
			style = "[style=filled, color=grey]"
		}/* else {
			//style = "[shape=circle]"
		}*/
		fmt.Printf("\t\"%s\" %s\n", tkn, style)
	}
	// all file nodes:
	for fn, rel := range files {
		fmt.Printf("\t\"%s\" [shape=box];\n", fn)
		// all publish relationships
		for i := 0; i < len(rel.publish); i++ {
			fmt.Printf("\t\"%s\"->\"%s\" [color=blue];\n", fn, rel.publish[i])
		}
		// all subscribe relationships
		for i := 0; i < len(rel.subcribe); i++ {
			fmt.Printf("\t\"%s\"->\"%s\" [color=orange];\n", rel.subcribe[i], fn)
		}
		// all stateSet relationships
		for i := 0; i < len(rel.stateSet); i++ {
			fmt.Printf("\t\"%s\"->\"%s\" [color=black];\n", fn, rel.stateSet[i])
		}
		// all stateGet relationships
		for i := 0; i < len(rel.stateGet); i++ {
			fmt.Printf("\t\"%s\"->\"%s\" [color=grey];\n", rel.stateGet[i], fn)
		}
	}
	fmt.Println("}")
}
