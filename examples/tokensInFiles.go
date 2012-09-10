package main

import (
	"flag"
	"fmt"
	"github.com/rudenoise/counting/count"
	"github.com/rudenoise/counting/dir"
	"path/filepath"
	"runtime"
	"sort"
)

var asc = flag.Bool("asc", false, "order ascending/descending")
var exclude = flag.String("exclude", "^$", "regexp pattern to exclude in file path")
var include = flag.String("include", "", "regexp pattern to include file path")
var lmt = flag.Int("limit", 0, "limit number of results")
var tokenRegExp = flag.String("tokenRegExp", "[a-zA-z]+", "regexp pattern to define a 'tokenRegExp'")

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	dirStr := filepath.Dir(flag.Arg(0))

	filePaths, err := dir.AllPaths(dirStr, *exclude, *include)
	if err != nil {
		panic(err)
	}

	limit := 5000
	length := len(filePaths)

	if length < limit {
		tMap := make(count.TokensMap)
		count.TokensInFiles(filePaths, *tokenRegExp, tMap)

		tSlice := tMap.ToSlice()
		if (*asc == true) {
			sort.Sort(count.TokenSliceByCountAsc{tSlice})
		} else {
			sort.Sort(count.TokenSliceByCountDesc{tSlice})
		}
		length = len(tSlice)
		if *lmt != 0 && *lmt <= length {
			length = *lmt
		}
		for i := 0; i < length; i++ {
			fmt.Printf("%7d\t%s\n", tSlice[i].Count, tSlice[i].Token)
		}

	} else {
		fmt.Printf("\nDidn't bother, you tried to meaure %d files, limit set to %d\n\n", length, limit)
	}
}
