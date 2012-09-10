package main

import (
	"flag"
	"fmt"
	"path/filepath"
	"github.com/rudenoise/counting/dir"
	"runtime"
	"github.com/rudenoise/counting/count"
	/*
	"sort"
	*/
)

var exclude = flag.String("exclude", "^$", "regexp pattern to exclude in file path")
var include = flag.String("include", "", "regexp pattern to include file path")
var lmt = flag.Int("limit", 0, "limit number of results")
var tokenRegExp = flag.String("tokenRegExp", "[a-zA-z]+", "regexp pattern to define a 'tokenRegExp'")

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	dirStr := filepath.Dir(flag.Arg(0))
	fmt.Println(*exclude, *include, *lmt, *tokenRegExp, dirStr)
	filePaths, err := dir.AllPaths(dirStr, *exclude, *include)
	if err != nil {
		panic(err)
	}
	fmt.Println(filePaths)
	limit := 5000
	length := len(filePaths)

	fmt.Println(limit, length)
	if length < limit {
		tMap := make(count.TokensMap)
		count.TokensInFiles(filePaths, *tokenRegExp, tMap)
		/*
		tSlice := tMap.ToSlice()
		//sort.Sort(count.TokenSliceByCountDesc{tSlice})
		if *lmt != 0 && *lmt <= length {
			length = *lmt
		}
		for i := 0; i < len(tSlice); i++ {
			fmt.Printf("%7d\t%s\n", tSlice[i].Count, tSlice[i].Token)
		}
		*/
	} else {
		fmt.Printf("\nDidn't bother, you tried to meaure %d files, limit set to %d\n\n", length, limit)
	}
}
