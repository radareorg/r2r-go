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
	"github.com/radare/r2pipe-go"
	"github.com/sergi/go-diff/diffmatchpatch"
	"encoding/json"
	"io/ioutil"
	"strings"
	"fmt"
	"os"
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
		fmt.Fprintln(os.Stderr, "Error:", err.Error())
		os.Exit(1)
	}
	var tests []R2Test
	json.Unmarshal(raw, &tests)
	fmt.Printf("Loaded", len(tests), "from", fpath, "\n")
	return tests
}

func runtest(test R2Test) bool {
	instance, err := r2pipe.NewPipe(test.File)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err.Error())
		return false;
	}
	defer instance.Close()
	if test.Commands != nil {
		for _, command := range test.Commands {
			buf, err := instance.Cmd(command)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error:", err.Error())
				return false;
			}
			str := string(buf)
			if strings.Compare(str, test.Expected) != 0 {
				fmt.Println(test.File)
				dmp := diffmatchpatch.New()
				diffs := dmp.DiffMain(str, test.Expected, false)
				fmt.Println(dmp.DiffPrettyText(diffs))
			}
		}
	}
	return true;
}

