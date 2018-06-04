// radare - LGPL - Copyright 2018 - deroad

package main

import (
	"encoding/json"
	"os/exec"
	"strings"
	"bufio"
	"bytes"
	"time"
	"fmt"
	"io"
)

// A Pipe represents a communication interface with r2 that will be used to
// execute commands and obtain their results.
type Pipe struct {
	File   string
	r2cmd  *exec.Cmd
	stdin  io.WriteCloser
	stdout io.ReadCloser
	Core   *struct{}
	cmd    CmdDelegate
	close  CloseDelegate
}

type CmdDelegate func(*Pipe, string) (string, error)
type CloseDelegate func(*Pipe) error

// NewPipe returns a new r2 pipe and initializes an r2 core that will try to
// load the provided file or URI. If file is an empty string, the env vars
// R2PIPE_{IN,OUT} will be used as file descriptors for input and output, this
// is the case when r2pipe is called within r2.

func NewPipe(args ...string) (*Pipe, error) {
	file := args[len(args) - 1]
	args[len(args) - 1] = "-q0"
	args = append(args, file)
	r2cmd := exec.Command("r2", args...)
	stdin, err := r2cmd.StdinPipe()
	if err != nil {
		fmt.Println("Error Stdin")
		return nil, err
	}
	stdout, err := r2cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Error Stdout")
		return nil, err
	}
	if err := r2cmd.Start(); err != nil {
		fmt.Println("Error Start")
		return nil, err
	}
	// Read initial data
	for i := 0; ; i++ {
		if _, err := bufio.NewReader(stdout).ReadString('\x00'); err != nil {
			if i < 4 {
				fmt.Println("Error Reader")
				return nil, err
			}
			time.Sleep(time.Second)
		} else {
			break
		}
	}

	r2p := &Pipe{
		File:   args[len(args) - 2],
		r2cmd:  r2cmd,
		stdin:  stdin,
		stdout: stdout,
	}
	return r2p, nil
}

// Write implements the standard Write interface: it writes data to the r2
// pipe, blocking until r2 have consumed all the data.
func (r2p *Pipe) Write(p []byte) (n int, err error) {
	return r2p.stdin.Write(p)
}

// Read implements the standard Read interface: it reads data from the r2
// pipe, blocking until the previously issued commands have finished.
func (r2p *Pipe) Read(p []byte) (n int, err error) {
	return r2p.stdout.Read(p)
}

// Cmd is a helper that allows to run r2 commands and receive their output.
func (r2p *Pipe) Cmd(cmd string) (string, error) {
	if r2p.Core != nil {
		if r2p.cmd != nil {
			return r2p.cmd(r2p, cmd)
		}
		return "", nil
	}
	if _, err := fmt.Fprintln(r2p, cmd); err != nil {
		return "", err
	}
	buf, err := bufio.NewReader(r2p).ReadString('\x00')
	if err != nil {
		return "", err
	}
	return strings.TrimRight(buf, "\x00"), nil
}

// Cmdj acts like Cmd but interprets the output of the command as json. It
// returns the parsed json keys and values.
func (r2p *Pipe) Cmdj(cmd string) (interface{}, error) {
	if _, err := fmt.Fprintln(r2p, cmd); err != nil {
		return nil, err
	}
	buf, err := bufio.NewReader(r2p).ReadBytes('\x00')
	if err != nil {
		return nil, err
	}
	buf = bytes.TrimRight(buf, "\x00")
	var output interface{}
	if err := json.Unmarshal(buf, &output); err != nil {
		return nil, err
	}
	return output, nil
}

// Close shuts down r2, closing the created pipe.
func (r2p *Pipe) Close() error {
	if r2p.close != nil {
		return r2p.close(r2p)
	}
	if r2p.File == "" {
		return nil
	}
	if _, err := r2p.Cmd("q!"); err != nil {
		return err
	}
	return r2p.r2cmd.Wait()
}
