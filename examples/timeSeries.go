package main

import(
	"encoding/json"
	"fmt"
	"flag"
	"path/filepath"
	"github.com/rudenoise/counting/dir"
	"github.com/rudenoise/counting/count"
	"os/exec"
)

type CountMap map[string] []int

// map for collecting counts
var countMap = CountMap{}

// flags
var re = flag.String("regExp", ".*", "regexp pattern to match file paths")
var steps = flag.Int("steps", 5, "number of git history commits to look back into")

func main() {
	flag.Parse()
	dirStr := filepath.Dir(flag.Arg(0))
	// loop over the previous x commits via git
	for i := *steps; i > 0; i-- {
		arg := fmt.Sprintf("master~%d", i)
		err := exec.Command("git", "checkout", arg).Run()
		if err != nil {
			panic(err)
		}
		countAll(getPaths(dirStr), *steps - i)
	}
	// reset repo to master
	err := exec.Command("git", "checkout", "master").Run()
	if err != nil {
		panic(err)
	}
	countAll(getPaths(dirStr), *steps)
	pathsSlice := mapToSlice(countMap)
	fmt.Println(pathsSlice)
	o, err := json.MarshalIndent(pathsSlice, "\n", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", o)
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
		file[position] = lines
	}
}

type Path struct {
	name string
	data []int
}

func mapToSlice(cMap CountMap) []Path {
	paths := make([]Path, 0)
	for k, v := range cMap {
		paths = append(paths, Path{k, v})
	}
	return paths
}
