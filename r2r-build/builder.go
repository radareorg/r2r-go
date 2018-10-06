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
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
)

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

func exists(name string) bool {
	_, err := os.Stat(name)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func decode64(encoded string) string {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Decode error:", err.Error())
		return ""
	}
	return string(decoded)
}

func multilinequote(scanner *bufio.Scanner, instr string) string {
	var s string = instr
	for scanner.Scan() {
		str := scanner.Text()
		if strings.HasPrefix(str, "'") {
			break
		}
		s += str + "\n"
	}
	return s
}

func populate(test *R2Test, str string, scanner *bufio.Scanner) bool {
	if strings.HasPrefix(str, "NAME=") {
		test.Name = str[5:]
		return true
	} else if strings.HasPrefix(str, "ARGS=") {
		test.Args = str[5:]
		return true
	} else if strings.HasPrefix(str, "FILE=") {
		test.File = str[5:]
		for strings.HasPrefix(test.File, "../") {
			test.File = test.File[3:]
		}
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
	} else if strings.HasPrefix(str, "EXPECT='") {
		str = str[8:]
		test.Expected = multilinequote(scanner, str)
		return true
	} else if strings.HasPrefix(str, "EXPECT=") {
		test.Expected = str[7:]
		return true
	} else if strings.HasPrefix(str, "CMDS64=") {
		cmd := decode64(str[7:])
		test.Commands = strings.Split(cmd, "\n")
		return true
	} else if strings.HasPrefix(str, "CMDS='") {
		str = str[6:]
		test.Commands = strings.Split(multilinequote(scanner, str), "\n")
		return true
	} else if strings.HasPrefix(str, "CMDS=") {
		str = str[6:]
		test.Commands = strings.Split(str, "\n")
		return true
	} else {
		fmt.Println("Unknown:", str)
	}
	return false
}

func populate_asm(name string, test *R2Test, str string, scanner *bufio.Scanner) bool {
	if len(str) < 1 {
		return false
	}
	args := strings.Split(name, "_")
	re := regexp.MustCompile("\\w+|\".+\"")
	found := re.FindAllString(str, -1)
	cmds := found[0]
	asm := found[1]
	hex := found[2]
	var skip string = "0x0"
	if len(found) > 3 {
		skip = found[3]
	}
	if len(cmds) < 4 {
		if len(args) == 1 {
			test.Args = "-a " + args[0]
		} else if len(args) == 2 {
			test.Args = "-a " + args[0] + " -b " + args[1]
		} else if len(args) == 3 {
			test.Args = "-a " + args[0] + " -e asm.cpu=" + args[1] + " -b " + args[2]
		} else {
			return false
		}
		test.Broken = false
		if skip != "0x0" {
			test.Commands = append(test.Commands, "s "+skip)
		}
		for i := 0; i < len(cmds); i++ {
			if cmds[i] == 'a' {
				test.Commands = append(test.Commands, "pa "+asm[1:len(asm)-1])
				test.Expected += hex + "\n"
			} else if cmds[i] == 'd' {
				test.Commands = append(test.Commands, "pad "+hex)
				test.Expected += asm[1:len(asm)-1] + "\n"
			} else if cmds[i] == 'B' {
				test.Broken = true
			} else if cmds[i] == 'E' {
				test.Args += " -e cfg.bigendian=true"
			} else {
				fmt.Println("Unknown flag:", cmds[i])
			}
		}
		test.File = "-"
		test.Name = name + ": " + str
		return true
	}
	return false
}

func build(infilepath string, outfilepath string) {
	var skipone bool = false
	var special string
	var str string
	var regr R2RegressionTest
	var e R2Test = R2Test{"", "", "", make([]string, 0), "", false}
	file, err := os.Open(infilepath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err.Error())
		os.Exit(1)
	}
	fmt.Println("Open:", path.Base(infilepath))
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for skipone || scanner.Scan() {
		str = scanner.Text()

		if strings.Compare(str, "RUN") == 0 {
			// fmt.Println(fmt.Sprintf(`Added: "%s"`, e.Name))
			regr.Tests = append(regr.Tests, e)
			e = R2Test{"", "", "", make([]string, 0), "", false}
			skipone = false
		} else if strings.HasPrefix(str, "CMDS=<<EXPECT") {
			special = "CMDS=" + str[13:]
			skipone = true
			for scanner.Scan() {
				str = scanner.Text()
				if strings.HasPrefix(str, "EXPECT=") {
					break
				}
				special += str + "\n"
			}
			populate(&e, special[:len(special)-1], scanner)
		} else if strings.HasPrefix(str, "EXPECT=<<RUN") {
			special = "EXPECT=" + str[12:]
			skipone = true
			for scanner.Scan() {
				str = scanner.Text()
				if strings.HasPrefix(str, "RUN") {
					break
				}
				special += str + "\n"
			}
			populate(&e, special[:len(special)-1], scanner)
		} else {
			if strings.Contains(infilepath, "/asm/") && populate_asm(path.Base(infilepath), &e, str, scanner) {
				regr.Tests = append(regr.Tests, e)
				e = R2Test{"", "", "", make([]string, 0), "", false}
			} else if !strings.Contains(infilepath, "/asm/") && !populate(&e, str, scanner) {
				fmt.Println("Unknown:", str)
			}
			skipone = false
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err.Error())
	}
	if strings.Contains(infilepath, "/asm/") {
		regr.Type = "asm"
		strings.Replace(outfilepath, ".json", ".asm.json", -1)
	} else {
		regr.Type = "cmd"
	}
	bytes, err := json.MarshalIndent(regr, "", "    ")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err.Error())
		os.Exit(1)
	}
	err = ioutil.WriteFile(outfilepath, bytes, 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err.Error())
	}
}
