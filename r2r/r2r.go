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
	"io/ioutil"
	"runtime"
	"strconv"
	"fmt"
	"os"
)

type ArgOption struct {
	Description string
	Argc int
	Callback func(... string)
}

func loadJSON(fpath string) []R2Test {
	raw, err := ioutil.ReadFile(fpath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err.Error())
		os.Exit(1)
	}
	var tests []R2Test
	json.Unmarshal(raw, &tests)
	return tests
}

var MAX_JOBS int = runtime.NumCPU()

var Options = map[string]ArgOption {
	"--jobs": {
		"defines how many jobs can be spawn. (if n < 1 then will be used the number of CPUs).",
		1,
		func(value... string) {
			if s, err := strconv.Atoi(value[0]); err == nil &&  s > 0 {
				MAX_JOBS = s
			} else {
				MAX_JOBS = runtime.NumCPU()
			}
		},
	},
}

func usage() {
	fmt.Println("Usage: ")
	for k, v := range Options { 
		fmt.Printf("%8s | %s (%d args)\n", k, v.Description, v.Argc)
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
		if arg == "--help" {
			usage()
		}
		pair, ok := Options[arg]
		if ok {
			max := i + pair.Argc
			if max < Argc {
				pair.Callback(os.Args[i:max]...)
			} else {
				badarg(arg)
			}
		} else {
			badarg(arg)
		}
	}
	filepath := string(os.Args[Argc - 1])
	if filepath == "--help" {
		usage()
	}
	tests := loadJSON(filepath)
	pool := NewR2Pool(4)
	if !pool.PerformTests(&tests) {
		os.Exit(1)
	}
}