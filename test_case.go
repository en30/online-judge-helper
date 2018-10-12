package main

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
)

type TestCase struct {
	Title  string `yaml:":title"`
	Input  string `yaml:":input"`
	Output string `yaml:":output"`
}

func (tc *TestCase) save(dir string) error {
	path := filepath.Join(dir, tc.Title+".in.txt")
	err := ioutil.WriteFile(path, []byte(tc.Input), 0644)
	if err != nil {
		return err
	}
	path = filepath.Join(dir, tc.Title+".out.txt")
	err = ioutil.WriteFile(path, []byte(tc.Output), 0644)
	if err != nil {
		return err
	}
	return nil
}

func (tc *TestCase) test(command string, args []string, timeLimit time.Duration) *TestCaseResult {
	ctx, cancel := context.WithTimeout(context.Background(), timeLimit)
	defer cancel()
	res := TestCaseResult{TestCase: *tc}

	cmd := exec.CommandContext(ctx, command, args...)
	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	err := cmd.Start()
	if err != nil {
		res.Status = IE
		res.Error = err
		return &res
	}
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		_, err := io.Copy(stdin, strings.NewReader(tc.Input))
		if e, ok := err.(*os.PathError); ok && e.Err == syscall.EPIPE {
			// ignore EPIPE
		} else if err != nil {
			log.Println("failed to write to STDIN", err)
		}
		stdin.Close()
		wg.Done()
	}()
	go func() {
		io.Copy(&res.Stdout, stdout)
		stdout.Close()
		wg.Done()
	}()
	go func() {
		io.Copy(&res.Stderr, stderr)
		stderr.Close()
		wg.Done()
	}()
	wg.Wait()
	err = cmd.Wait()
	if ctx.Err() == context.DeadlineExceeded {
		res.Status = TLE
	} else if err != nil {
		res.Status = RE
		res.Error = err
	} else if res.Stdout.String() == tc.Output {
		res.Status = AC
	} else {
		res.Status = WA
	}
	return &res
}

type TestCaseResult struct {
	TestCase TestCase
	Status   Status
	Stdout   bytes.Buffer
	Stderr   bytes.Buffer
	Error    error
}
