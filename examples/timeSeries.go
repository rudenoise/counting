package main

import(
	"fmt"
	"flag"
	"path/filepath"
	"github.com/rudenoise/counting/dir"
	"github.com/rudenoise/counting/count"
	"os/exec"
)

// map for collecting counts
var countMap = map[string] []int {}

// flags
var re = flag.String("regExp", ".*", "regexp pattern to match file paths")
var steps = flag.Int("steps", 5, "number of git history commits to look back into")

func main() {
	flag.Parse()
	dirStr := filepath.Dir(flag.Arg(0))
	// loop over the previous x commits via git
	for i := *steps; i > 0; i-- {
		arg := fmt.Sprintf("master~%d", i)
		out, err := exec.Command("git", "checkout", arg).Output()
		if err != nil {
			panic(err)
		}
		fmt.Println(out)
		countAll(getPaths(dirStr), i)
	}
	// reset repo to master
	out, err := exec.Command("git", "checkout", "master").Output()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", out)
	countAll(getPaths(dirStr), i)
	fmt.Println(countMap)
}

func getPaths(dirStr string) []string {
	paths, err := dir.AllPaths(dirStr, "^$", *re)
	if err != nil {
		panic(err)
	}
	return paths
}

func countAll(paths []string, position int) {
	for i := 0; i < len(paths); i++ {
		lines, err := count.LinesInFile(paths[i], true)
		if err != nil {
			panic(err)
		}
		file, ok := countMap[paths[i]]
		if ok == false {
			countMap[paths[i]] = make([]int, *steps + 1)
			file = countMap[paths[i]]
		}
		file[*steps - position] = lines
	}
}
