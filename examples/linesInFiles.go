package main

import (
	"flag"
	"fmt"
	"github.com/rudenoise/counting/count"
	"github.com/rudenoise/counting/dir"
	"path/filepath"
	"runtime"
	"sort"
	"sync"
)

var exclude = flag.String("exclude", "^$", "regexp pattern to exclude in file path")
var include = flag.String("include", "", "regexp pattern to include file path")
var lmt = flag.Int("limit", 0, "limit number of results")
var ignoreCommentsEmptyLines = flag.Bool("icel", false, "ignore comment and empty lines")

type File struct {
	FilePath  string
	LineCount int
}

type Files []File

func (f Files) Len() int           { return len(f) }
func (f Files) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
func (f Files) Less(i, j int) bool { return f[i].LineCount < f[j].LineCount }

type ByLengthReverse struct{ Files }

func (f ByLengthReverse) Less(i, j int) bool {
	return f.Files[i].LineCount > f.Files[j].LineCount
}

func OpenParallel(paths []string) Files {
	var files Files
	var wg sync.WaitGroup
	for i := 0; i < len(paths); i++ {
		wg.Add(1)
		go func(fp string) {
			lines, err := count.LinesInFile(fp, *ignoreCommentsEmptyLines)
			if err != nil {
				panic(err)
			}
			files = append(files, File{fp, lines})
			wg.Done()
		}(paths[i])
	}
	wg.Wait()
	return files
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	dirStr := filepath.Dir(flag.Arg(0))
	pths, err := dir.AllPaths(dirStr, *exclude, *include)
	total := 0
	if err != nil {
		panic(err)
	}
	limit := 5000
	if len(pths) < limit {
		files := OpenParallel(pths)
		sort.Sort(ByLengthReverse{files})
		// decide length of output
		length := len(files)
		if *lmt != 0 && *lmt <= length {
			length = *lmt
		}
		for i := 0; i < length; i++ {
			total += files[i].LineCount
			fmt.Printf("%7d\t%s\n", files[i].LineCount, files[i].FilePath)
		}
		fmt.Printf("\n%7d\ttotal\n\n", total)
	} else {
		fmt.Printf("\nDidn't bother, you tried to meaure %d files, limit set to %d\n\n", len(pths), limit)
	}
}
