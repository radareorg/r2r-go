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
	"github.com/pmezard/go-difflib/difflib"
//	"github.com/radare/r2pipe-go"
	"encoding/json"
	"io/ioutil"
	"strings"
	"bytes"
	"fmt"
	"os"
)

func diff(str1, str2 string) string {
	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(str1),
		B:        difflib.SplitLines(str2),
		FromFile: "expected",
		ToFile:   "r2pipe",
		Context:  3,
	}
	text, _ := difflib.GetUnifiedDiffString(diff)
	return text
}

func load(fpath string) []R2Test {
	raw, err := ioutil.ReadFile(fpath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err.Error())
		os.Exit(1)
	}
	var tests []R2Test
	json.Unmarshal(raw, &tests)
	return tests
}

type TestResult struct {
	Message string
	Success bool
	Error bool
	Test *R2Test
}

func (result TestResult) Print(printall bool) bool {
	if result.Error {
		fmt.Println("[XX]", result.Test.Name, "something went really wrong.")
		fmt.Println(result.Message)
	} else if result.Success {
		if result.Test.Broken {
			fmt.Println("[FX]", result.Test.Name)
		} else if printall {
			fmt.Println("[OK]", result.Test.Name)
		}
		return true
	} else if result.Test.Broken {
		fmt.Println("[BR]", result.Test.Name)
		return true
	} else {
		fmt.Println("[XX]", result.Test.Name, result.Test.File, result.Test.Args)
		fmt.Println(result.Message)
	}
	return false
}

type R2Test struct {
	Name string `json:"name"`
	File string `json:"file"`
	Args string `json:"args"`
	Commands []string `json:"commands"`
	Expected string `json:"expected"`
	Broken bool `json:"broken"`
}

func (test R2Test) Exec() *TestResult {
	result := &TestResult{"", false, false, &test}
	result.Success = true
	result.Error = false
	var args []string = strings.Split(test.Args, " ")
	args = append(args, test.File)
	instance, err := NewPipe(args...)
	if err != nil {
		result.Message = fmt.Sprintf("Error: %s\n", err.Error())
		result.Success = false
		result.Error = true
		return result;
	}
	defer instance.Close()
	if test.Commands != nil {
		var buffer bytes.Buffer
		for _, command := range test.Commands {
			output, err := instance.Cmd(command)
			if err != nil {
				result.Message = fmt.Sprintf("Error: %s\n", err.Error())
				result.Success = false
				result.Error = true
				return result;
			}
			t := string(output)
			if len(t) > 0 {
				buffer.WriteString(t)
			}
		}
		// simple workaround for bad endline
		if len(buffer.String()) < len(test.Expected) {
			buffer.WriteString("\n")
		} 
		str := buffer.String()
		if strings.Compare(str, test.Expected) != 0 {
			diffs := diff(test.Expected, str)
			result.Message = diffs
			result.Success = false
		}
	}
	return result
}

