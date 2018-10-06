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
	"bytes"
	"fmt"
	"github.com/pmezard/go-difflib/difflib"
	"os"
	"strings"
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

type TestResult struct {
	Message string
	Success bool
	Error   bool
	Test    *R2Test
	Options *TestsOptions
}

func (result TestResult) Print(printall bool) bool {
	if result.Error {
		fmt.Println("[XX]", result.Test.Name, "something went really wrong.")
		result.Options.Println("r2", result.Test.Args, result.Test.File)
		result.Options.Println(strings.Join(result.Test.Commands, "; "))
		fmt.Println(result.Message)
	} else if result.Success {
		if result.Test.Broken {
			fmt.Println("[FX]", result.Test.Name)
		} else if printall && !result.Options.ErrorsOnly {
			fmt.Println("[OK]", result.Test.Name)
		}
		return true
	} else if result.Test.Broken {
		if !result.Options.ErrorsOnly {
			fmt.Println("[BR]", result.Test.Name)
		}
		return true
	} else {
		fmt.Println("[XX]", result.Test.Name)
		result.Options.Println("r2", result.Test.Args, result.Test.File)
		fmt.Println(result.Message)
	}
	return false
}

type R2Test struct {
	Name     string   `json:"name"`
	File     string   `json:"file"`
	Args     string   `json:"args"`
	Commands []string `json:"commands"`
	Expected string   `json:"expected"`
	Broken   bool     `json:"broken"`
}

type R2RegressionTest struct {
	Type  string   `json:"type"`
	Tests []R2Test `json:"tests"`
}

func (test R2Test) Exec(options *TestsOptions) *TestResult {
	result := &TestResult{"", false, false, &test, options}
	result.Success = true
	result.Error = false
	var args []string = strings.Split(test.Args, " ")
	args = append(args, test.File)
	instance, err := NewPipe(args...)
	if err != nil {
		result.Message = fmt.Sprintf("Error: %s", err.Error())
		result.Success = false
		result.Error = true
		if _, err := os.Stat(test.File); os.IsNotExist(err) {
			result.Message = fmt.Sprintf("Error: File %s doesn't exists", test.File)
			result.Success = false
			result.Error = true
			return result
		}
		return result
	}
	defer instance.Close()
	if test.Commands != nil {
		var buffer bytes.Buffer
		for _, command := range test.Commands {
			if command == "q" {
				continue
			}
			output, err := instance.Cmd(command)
			if err != nil {
				result.Message = fmt.Sprintf("Error: %s", err.Error())
				result.Success = false
				result.Error = true
				return result
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
	if options.Sequence {
		result.Print(true)
	}

	return result
}
