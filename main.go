package main

import (
	"flag"
	"fmt"
	"github.com/kr/pty"
	"io"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"syscall"
)

var (
	Version string
	Build   string
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	configPath := flag.String("c", "megahorn.json", "config json file path")

	flag.Usage = func() {
		fmt.Printf("megahorn v%s (%s)\n\n", Version, Build)
		fmt.Println("Usage:")
		fmt.Println("  megahorn [options] cmd [args...]")
		fmt.Println("  cmd | megahorn [options]\n")
		fmt.Println("Options:")
		flag.PrintDefaults()
	}

	flag.Parse()

	config := newConfig()
	err := config.LoadFile(*configPath)
	errWithKill(err)

	stdoutPty, fakeStdout, err := pty.Open()
	errWithKill(err)

	stderrPty, fakeStderr, err := pty.Open()
	errWithKill(err)

	cmdMode := flag.NArg() > 0

	oDrivers := config.StdoutDrviers(cmdMode)
	eDrivers := config.StderrDrviers(cmdMode)

	oWriters := []io.Writer{}
	for _, od := range oDrivers {
		oWriters = append(oWriters, od.(io.Writer))
	}
	eWriters := []io.Writer{}
	for _, ed := range eDrivers {
		eWriters = append(eWriters, ed.(io.Writer))
	}

	oWriter := io.MultiWriter(oWriters...)
	eWriter := io.MultiWriter(eWriters...)

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		_, err := io.Copy(oWriter, stdoutPty)
		if err != nil {
			os.Stderr.WriteString(err.Error())
		}
	}()

	go func() {
		defer wg.Done()
		_, err := io.Copy(eWriter, stderrPty)
		if err != nil {
			os.Stderr.WriteString(err.Error())
		}
	}()

	exitStatus := 0

	if cmdMode {
		cmd := flag.Arg(0)
		args := []string{}

		if len(flag.Args()) > 1 {
			args = flag.Args()[1:]
		}

		os.Stdin.Close()
		runner := exec.Command(cmd, args...)

		runner.Dir = config.WorkingDir
		runner.Env = config.Environments()
		runner.Stdout = fakeStdout
		runner.Stderr = fakeStderr

		if err := runner.Run(); err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
					exitStatus = status.ExitStatus()
				}
			} else {
				errWithKill(err)
			}
		}
	} else {
		io.Copy(stdoutPty, os.Stdin)
	}

	fakeStdout.Close()
	fakeStderr.Close()

	wg.Add(2)
	go func() {
		closeAll(oDrivers...)
		wg.Done()
	}()
	go func() {
		closeAll(eDrivers...)
		wg.Done()
	}()

	wg.Wait()

	os.Exit(exitStatus)
}

func closeAll(closers ...io.WriteCloser) {
	wg := sync.WaitGroup{}
	for _, closer := range closers {
		wg.Add(1)
		go func(closer io.WriteCloser) {
			closer.Close()
			wg.Done()
		}(closer)
	}
	wg.Wait()
}

func errWithKill(err error) {
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}
}
