package main

import (
	"flag"
	"fmt"
	"github.com/rudenoise/counting/count"
	"github.com/rudenoise/counting/dir"
	"path/filepath"
	"regexp"
	"strings"
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
var channelFilter = flag.String("cf", "", "comma separated list of channels to include")

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
		// collect all state and pubSub interactions for this file
		files[filePaths[i]] = File{
			loopTokens(pubMap, compiledTokenRE, allTokens, 0),
			loopTokens(subMap, compiledTokenRE, allTokens, 0),
			loopTokens(stateSetMap, compiledTokenRE, allTokens, 1),
			loopTokens(stateGetMap, compiledTokenRE, allTokens, 1),
		}
	}
	// now we have a map of all files and their publshed and subscribed keys
	printInDotFormat(
		allTokens,
		filterEmptyFiles(files),
	)
}

func loopTokens(
	spMap count.TokensMap,
	compiledTokenRE *regexp.Regexp,
	allTokens AllTokens, tokenType int,
) []string {
	tokens := []string{}
	filter, total := chopUpChannelFilterFlag()
	for token := range spMap {
		token = compiledTokenRE.FindString(token)
		if total == 0 {
			tokens = append(tokens, token)
			// add to deduped list
			allTokens[token] = tokenType
		} else {
			if filter[token] == true {
				tokens = append(tokens, token)
				// add to deduped list
				allTokens[token] = tokenType
			}
		}
	}
	return tokens
}

func chopUpChannelFilterFlag () (map[string]bool, int) {
	// get channels to filter
	chanFilterMap := make(map[string]bool)
	if *channelFilter != "" {
		channelFilterSlice := strings.Split(*channelFilter, ",")
		total := 0
		for _, v := range channelFilterSlice {
			chanFilterMap[v] = true
			total ++
		}
		return chanFilterMap, total
	} else {
		return chanFilterMap, 0
	}
}

func filterAllTokens(allTokens AllTokens, filter map[string]bool) AllTokens {
	filtered := make(AllTokens)
	for token, val := range allTokens {
		if filter[token] == true {
			filtered[token] = val
		}
	}
	return filtered
}

func filterTokens(
	tokens []string,
	filter map[string]bool,
) []string {
	filtered := []string{}
	for _, token := range tokens {
		if filter[token] == true {
			filtered = append(filtered, token)
		}
	}
	return filtered
}

func filterEmptyFiles(
	files map[string]File,
) map[string]File {
	filtered := make(map[string]File)
	for name, file := range files {
		if (len(file.publish) > 0 || len(file.subcribe) > 0 || len(file.stateSet) > 0 || len(file.stateGet) > 0) {
			filtered[name] = file
		}
	}
	return filtered
}

func printInDotFormat(
	allTokens AllTokens,
	files map[string]File,
) {
	// start printing
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

