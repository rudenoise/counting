package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/rudenoise/counting/count"
	"github.com/rudenoise/counting/dir"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

type CountMap map[string][]int

// map for collecting counts
var countMap = CountMap{}

// flags
var re = flag.String("regExp", ".*", "regexp pattern to match file paths")
var steps = flag.Int("steps", 5, "number of git history commits to look back into")
var interval = flag.Int("interval", 1, "intervals between steps")
var top = flag.Int("top", 10, "top X largest files")

func main() {
	flag.Parse()
	dirStr := filepath.Dir(flag.Arg(0))
	// loop over the previous x commits via git
	for i := (*steps * *interval); i > 0; i -= *interval {
		arg := fmt.Sprintf("master~%d", i)
		err := exec.Command("git", "checkout", arg).Run()
		if err != nil {
			panic(err)
		}
		countAll(getPaths(dirStr), *steps-i)
	}
	// reset repo to master
	err := exec.Command("git", "checkout", "master").Run()
	if err != nil {
		panic(err)
	}
	countAll(getPaths(dirStr), *steps)
	pathsSlice := mapToSlice(countMap)
	if *top < len(pathsSlice) {
		pathsSlice = pathsSlice[0:*top]
	}
	o, err := json.Marshal(pathsSlice)
	if err != nil {
		panic(err)
	}
	jsonStr := fmt.Sprintf("%s", o)
	jsonStr = strings.Replace(jsonStr, "Data", "data", -1)
	jsonStr = strings.Replace(jsonStr, "Name", "name", -1)
	fmt.Printf(printHTML(jsonStr))
}

func printHTML(json string) string {
	s := "<html><head>"
	s += "<script src=\"http://ajax.googleapis.com/ajax/libs/jquery/1.8.1/jquery.min.js\"></script>"
	s += "<script src=\"http://code.highcharts.com/highcharts.js\"></script>"
	s += "<script src=\"http://code.highcharts.com/modules/exporting.js\"></script>"
	s += "</head><body>"
	s += "<div id=\"container\" style=\"min-width: 400px; height: 600px; margin: 0 auto\"></div>"
	s += "<script>$(function(){var chart;$(document).ready(function(){chart=new Highcharts.Chart({chart:{renderTo:'container',type:'line',marginRight:230,marginBottom:25},title:{text:'File Size (Lines)'},xAxis:{categories:[]},yAxis:{title:{text:'Lines'},plotLines:[{value:0,width:1,color:'#808080'}]},tooltip:{formatter:function(){return'<b>'+this.series.name+'</b><br/>'+this.y;}},legend:{layout:'vertical',align:'right',verticalAlign:'top',x:-10,y:100,borderWidth:0},series:"
	s += json
	s+="});});});</script>"
	s += "</body></html>"
	return s
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
			countMap[paths[i]] = make([]int, *steps+1)
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

func (p Paths) Len() int      { return len(p) }
func (p Paths) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p Paths) Less(i, j int) bool {
	l := len(p[i].Data) - 1
	return p[i].Data[l] < p[j].Data[l]
}

type PathsReverse struct{ Paths }

func (p PathsReverse) Less(i, j int) bool {
	l := len(p.Paths[i].Data) - 1
	return p.Paths[i].Data[l] > p.Paths[j].Data[l]
}

func mapToSlice(cMap CountMap) []Path {
	paths := make(Paths, 0)
	for k, v := range cMap {
		paths = append(paths, Path{k, v})
	}
	sort.Sort(PathsReverse{paths})
	return paths
}
