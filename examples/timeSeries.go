package main

import(
	"encoding/json"
	"fmt"
	"flag"
	"path/filepath"
	"github.com/rudenoise/counting/dir"
	"github.com/rudenoise/counting/count"
	"os/exec"
	"strings"
	"sort"
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
	o, err := json.Marshal(pathsSlice)
	if err != nil {
		panic(err)
	}
	jsonStr := fmt.Sprintf("%s", o)
	jsonStr = strings.Replace(jsonStr, "Data", "data", -1)
	jsonStr = strings.Replace(jsonStr, "Name", "name", -1)
	fmt.Printf(jsonStr)
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
	Name string
	Data []int
}

type Paths []Path
func (p Paths) Len() int           { return len(p) }
func (p Paths) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p Paths) Less(i, j int) bool {
	l := len(p[i].Data) - 1
	return p[i].Data[l] < p[j].Data[l]
}

type PathsReverse struct { Paths }

func (p PathsReverse) Less (i, j int) bool {
	l := len(p.Paths[i].Data) - 1
	return p.Paths[i].Data[l] > p.Paths[j].Data[l]
}

func mapToSlice(cMap CountMap) []Path {
	paths := make(Paths, 0)
	for k, v := range cMap {
		paths = append(paths, Path{k, v})
	}
	sort.Sort(PathsReverse{ paths })
	return  paths
}
