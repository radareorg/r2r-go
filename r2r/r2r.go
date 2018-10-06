/*
 * Copyright (c) 2018, Giovanni Dante Grazioli <deroad@libero.it>
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * * Redistributions of source code must retain the above copyright notice, this
 *   list of conditions and the following disclaimer.
 * * Redistributions in binary form must reproduce the above copyright notice,
 *   this list of conditions and the following disclaimer in the documentation
 *   and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
 * LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 * CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 * SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 * INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 * CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 * POSSIBILITY OF SUCH DAMAGE.
 */

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
)

type ArgOption struct {
	Description string
	Argc        int
	Callback    func(...string)
}

func loadJSON(fpath string) R2RegressionTest {
	raw, err := ioutil.ReadFile(fpath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err.Error())
		os.Exit(1)
	}
	var tests R2RegressionTest

	if err := json.Unmarshal(raw, &tests); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err.Error())
		os.Exit(1)
	}
	return tests
}

var options TestsOptions = TestsOptions{
	false,
	false,
	false,
	runtime.NumCPU(),
}

var ArgsOptions = map[string]ArgOption{
	"--jobs": {
		"defines how many jobs can be spawn. (if n < 1 then will be used the number of CPUs).",
		1,
		func(value ...string) {
			s, err := strconv.Atoi(value[0])
			if err != nil || s < 1 {
				fmt.Println(err)
				os.Exit(1)
			}
			options.Jobs = s
		},
	},
	"--wdir": {
		"changes the current working directory",
		1,
		func(value ...string) {
			if err := os.Chdir(value[0]); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	},
	"--debug": {
		"enables debug output",
		0,
		func(value ...string) {
			options.Debug = true
		},
	},
	"--seq": {
		"enables sequence output",
		0,
		func(value ...string) {
			options.Sequence = true
		},
	},
	"--errors-only": {
		"enables only errors (and fixed) output",
		0,
		func(value ...string) {
			options.ErrorsOnly = true
		},
	},
}

func usage() {
	fmt.Println("Usage: ")
	for k, v := range ArgsOptions {
		fmt.Printf("%15s | %s (%d args)\n", k, v.Description, v.Argc)
	}
	os.Exit(1)
}

func badarg(arg string) {
	fmt.Printf("Invalid argument '%s'\n", arg)
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println(string(os.Args[0]), "[options] <file.json>")
		os.Exit(1)
	}
	Argc := len(os.Args)
	for i := 1; i < (Argc - 1); i++ {
		arg := string(os.Args[i])
		if arg == "--help" || arg == "-h" {
			usage()
		}
		pair, ok := ArgsOptions[arg]
		if ok {
			max := i + pair.Argc + 1
			if max < Argc {
				args := os.Args[i+1 : max]
				pair.Callback(args...)
				i = max - 1
			} else {
				badarg(arg)
			}
		} else {
			badarg(arg)
		}
	}
	filepath := string(os.Args[Argc-1])
	if filepath == "--help" || filepath == "-h" {
		usage()
	}
	fmt.Println("Executing", filepath)

	tests := loadJSON(filepath)
	pool := NewR2Pool(&options)
	if !pool.PerformTests(&tests) {
		os.Exit(1)
	}
}
