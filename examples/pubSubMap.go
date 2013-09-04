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

type AllTokens map[string]int

// set up flag defraults
var exclude = flag.String("e", "^$", "regexp pattern to exclude in file path")
var include = flag.String("i", "", "regexp pattern to include file path")
var fileShape = flag.String("f", "box", "a graphviz/dot shape for files")
var channelShape = flag.String("c", "oval", "a graphviz/dot shape for pubSub/state channels")

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
	allTokens := make(AllTokens)
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
			// collect all state and pubSub interactions for this file
			files[filePaths[i]] = File{
				loopTokens(pubMap, compiledTokenRE, allTokens, 0),
				loopTokens(subMap, compiledTokenRE, allTokens, 0),
				loopTokens(stateSetMap, compiledTokenRE, allTokens, 1),
				loopTokens(stateGetMap, compiledTokenRE, allTokens, 1),
			}
		}
	}
	// now we have a map of all files and their publshed and subscribed keys
	printInDotFormat(allTokens, files)
}

func loopTokens(
	spMap count.TokensMap,
	compiledTokenRE *regexp.Regexp,
	allTokens AllTokens, tokenType int,
) []string {
	tokens := []string{}
	for token := range spMap {
		token = compiledTokenRE.FindString(token)
		tokens = append(tokens, token)
		// add to deduped list
		allTokens[token] = tokenType
	}
	return tokens
}

func printInDotFormat(allTokens AllTokens, files map[string]File) {
	fmt.Printf("digraph PubSub{\n")
	// create dot file output
	// all token nodes
	for tkn, val := range allTokens {
		if val == 1 {
			fmt.Printf("\t\"%s\" [shape=%s, style=filled, color=grey]\n", tkn, *channelShape)
		} else {
			fmt.Printf("\t\"%s\" [shape=%s]\n", tkn, *channelShape)
		}
	}
	// all file nodes:
	for fn, rel := range files {
		fmt.Printf("\t\"%s\" [shape=%s];\n", fn, *fileShape)
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
