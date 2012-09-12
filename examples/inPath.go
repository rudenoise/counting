package main

import(
	"flag"
	"path/filepath"
	"fmt"
	"github.com/rudenoise/counting/dir"
)

var re = flag.String("regExp", ".*", "regexp pattern to match file paths")

func main() {
	flag.Parse()
	dirStr := filepath.Dir(flag.Arg(0))
	paths, err := dir.AllPaths(dirStr, "^$", *re)
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(paths); i++ {
		fmt.Println(paths[i])
	}
}
