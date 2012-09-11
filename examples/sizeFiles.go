package main

import (
	"flag"
	"fmt"
	"github.com/rudenoise/counting/count"
	"github.com/rudenoise/counting/dir"
	"os"
	"path/filepath"
	"sort"
)

type File struct {
	FilePath string
	ByteLen  int64
	Lines    int
}

type Files []File

func (p Files) Len() int      { return len(p) }
func (p Files) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

type ByByteLen struct{ Files }

func (p ByByteLen) Less(i, j int) bool { return p.Files[i].ByteLen < p.Files[j].ByteLen }

type ByByteLenReverse struct{ Files }

func (p ByByteLenReverse) Less(i, j int) bool { return p.Files[i].ByteLen > p.Files[j].ByteLen }

type ByLines struct{ Files }

func (p ByLines) Less(i, j int) bool { return p.Files[i].Lines < p.Files[j].Lines }

type ByLinesReverse struct{ Files }

func (p ByLinesReverse) Less(i, j int) bool { return p.Files[i].Lines > p.Files[j].Lines }

// cmd flags
var byBytes = flag.Bool("bytes", false, "order by byte length or line length")
var asc = flag.Bool("asc", false, "order ascending/descending")
var colour = flag.Bool("c", false, "colour output")
var ignoreCommentsEmptyLines = flag.Bool("icel", false, "ignore comments and empty lines")
var lmt = flag.Int("limit", 0, "limit number of results")
var exclude = flag.String("exclude", "^$", "regexp pattern to exclude in file path")
var include = flag.String("include", "", "regexp pattern to include file path")

func out(p Files, limit int) {
	// print headings
	fmt.Printf("%7s\t%7s\t%s\n", "Lines", "Bytes", "Path")
	for i := 0; i < limit; i++ {
		// print out all collected file info
		if p[i].FilePath != "." {
			fmt.Printf("%7d\t%7d\t%s\n", p[i].Lines, p[i].ByteLen, p[i].FilePath)
		}
	}
}

func fullInfo(fps []string, elc bool) Files {
	var files Files
	for i := 1; i < len(fps); i++ {
		f, err := os.Open(fps[i])
		if err != nil {
			panic(err)
		}
		info, err := f.Stat()
		if err != nil {
			panic(err)
		}
		lines, err := count.LinesInFile(fps[i], elc)
		if err != nil {
			panic(err)
		}
		files = append(files, File{fps[i], info.Size(), lines})
	}
	return files
}

func main() {
	flag.Parse()
	// start at current dir or path from args
	dirStr := filepath.Dir(flag.Arg(0))
	fps, err := dir.AllPaths(dirStr, *exclude, *include)
	if err != nil {
		panic(err)
	}
	filePaths := fullInfo(fps, *ignoreCommentsEmptyLines)
	// order by args
	if *byBytes {
		if *asc {
			sort.Sort(ByByteLen{filePaths})
		} else {
			sort.Sort(ByByteLenReverse{filePaths})
		}
	} else {
		if *asc {
			sort.Sort(ByLines{filePaths})
		} else {
			sort.Sort(ByLinesReverse{filePaths})
		}
	}
	// decide length of output
	length := len(filePaths)
	if *lmt != 0 && *lmt <= length {
		length = *lmt
	}
	// print output
	out(filePaths, length)
}
