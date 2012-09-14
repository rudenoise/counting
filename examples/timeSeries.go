package main

import(
	"fmt"
	"flag"
	"path/filepath"
	"github.com/rudenoise/counting/dir"
	"os/exec"
)

type dataPoint struct {
	fileName string
	data []int
}

type series []dataPoint

// flags
var re = flag.String("regExp", ".*", "regexp pattern to match file paths")

func main() {
	flag.Parse()
	dirStr := filepath.Dir(flag.Arg(0))
	// loop over the previous 5 commits via git
	for i := 5; i > 0; i-- {
		arg := fmt.Sprintf("master~%d", i)
		fmt.Println(arg)
		out, err := exec.Command("git", "checkout", arg).Output()
		if err != nil {
			panic(err)
		}
		fmt.Println(out)
		fmt.Println(getPaths(dirStr))
	}
	// reset repo to master
	out, err := exec.Command("git", "checkout", "master").Output()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", out)
}

func getPaths(dirStr string) []string {
	paths, err := dir.AllPaths(dirStr, "^$", *re)
	if err != nil {
		panic(err)
	}
	return paths
}
