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
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"strings"
	"strconv"
	"bufio"
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

func decode64(encoded string) string {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Decode error:", err.Error())
		return ""
	}
	return string(decoded)
}

func populate(test *R2Test, str string) bool {
	if strings.HasPrefix(str, "NAME=") {
		test.Name = str[5:]
		return true
	} else if strings.HasPrefix(str, "FILE=") {
		test.File = str[5:]
		return true
	} else if strings.HasPrefix(str, "BROKEN=") {
		if s, err := strconv.Atoi(str[7:]); err == nil {
			test.Broken = s == 1
		} else {
			test.Broken = false
		}
		return true
	} else if strings.HasPrefix(str, "EXPECT64=") {
		test.Expected = decode64(str[9:])
		return true
	} else if strings.HasPrefix(str, "CMDS64=") {
		test.Commands = strings.Split(decode64(str[7:]), "\n")
		return true
	}
	return false;
}

func build(infilepath string, outfilepath string) {
	var tests []R2Test
	var e R2Test = R2Test{"","", make([]string, 0),"", false}
	file, err := os.Open(infilepath)
    if err != nil {
	fmt.Fprintln(os.Stderr, "Error:", err.Error())
	os.Exit(1)
    }
    defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		str := scanner.Text()
		
		if !populate(&e, str) && strings.Compare(str, "RUN") == 0 {
			fmt.Println(fmt.Sprintf(`Added: "%s"`, e.Name))
			tests = append(tests, e)
			e = R2Test{"","", make([]string, 0),"", false}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err.Error())
	}
	bytes, err := json.Marshal(tests)
    if err != nil {
	fmt.Fprintln(os.Stderr, "Error:", err.Error())
	os.Exit(1)
    }
    err = ioutil.WriteFile(outfilepath, bytes, 0644)
    if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err.Error())
    }
}
