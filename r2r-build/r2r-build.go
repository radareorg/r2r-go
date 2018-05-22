package main

import (
	"fmt"
	"os"
)

func exists(name string) bool {
	_, err := os.Stat(name)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println(os.Args[0], "<path/regression/test> <file.json>")
		os.Exit(1)
	}
	filepath := os.Args[1]
	outputpath := os.Args[2]
	if !exists(filepath) {
		fmt.Println(filepath, "doesn't exists!")
		os.Exit(1)
	}
	fmt.Println("TESTS: ", filepath)
	fmt.Println("OUTPUT:", outputpath)
	build(filepath, outputpath)
}