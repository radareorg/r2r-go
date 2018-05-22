package main

import (
	"fmt"
	"os"
)


func main() {
	if len(os.Args) != 2 {
		fmt.Printf(fmt.Sprintf("%s <file.json>\n", os.Args[0]))
		os.Exit(1)
	}
	filepath := os.Args[1]
	fmt.Printf("Loading tests from", filepath, "\n")
	loadtests(filepath)
}