package main

import(
	"fmt"
	"os/exec"
)

type dataPoint struct {
	fileName string
	data []int
}

type series []dataPoint

func main() {
	/*
	a := make([]int, 10)
	fmt.Println("YO", a)
	*/
	// loop over the previous 5 commits via git
	for i := 5; i > 0; i-- {
		arg := fmt.Sprintf("master~%d", i)
		fmt.Println(arg)
		out, err := exec.Command("git", "checkout", arg).Output()
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s", out)
	}
	// reset repo to master
	out, err := exec.Command("git", "checkout", "master").Output()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", out)
}
