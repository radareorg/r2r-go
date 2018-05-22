package main

import (
	"github.com/radare/r2pipe-go"
	"encoding/json"
	"io/ioutil"
	"fmt"
)

type R2Test struct {
	Name string `json:"name"`
	File string `json:"file"`
	Commands []string `json:"commands"`
	Expected string `json:"expected"`
	Broken bool `json:"broken"`
}

func load(fpath string) []R2Test {
	raw, err := ioutil.ReadFile(fpath)
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}
	var tests []R2Test
	json.Unmarshal(raw, &tests)
	fmt.Printf("Loaded", len(tests), "from", fpath, "\n")
	return tests
}

func runtest(test R2Test) {
	r2p, err := r2pipe.NewPipe(test.File)
	if err != nil {
		fmt.FPrintln("Error:", err.Error())
		return false;
	}
	defer r2p.Close()
	if test.Commands != nil {
		for index, command := range test.Commands {
			if command != nil {				
			buf, err = r2p.Cmd(command)
			if err != nil {
				fmt.Println(string(buf))
				fmt.Println(command)
				fmt.FPrintln("Error:", err.Error())
				return false;
			}
			}
		}
	}
	return true;
}
