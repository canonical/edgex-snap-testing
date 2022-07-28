package utils

import (
	"bufio"
	"fmt"
	goexec "os/exec"
	"sync"
	"testing"
)

func Exec(t *testing.T, command string) (stdout, stderr string, err error) {
	return exec(t, command, false)
}

// exec executes a command
func exec(t *testing.T, command string, verbose bool) (stdout, stderr string, err error) {
	logf(t, "[exec] %s", command)

	cmd := goexec.Command("/bin/sh", "-c", command)

	var wg sync.WaitGroup

	// standard output
	outStream, err := cmd.StdoutPipe()
	if err != nil {
		// fatalf(t, err.Error())
		// return
		if t != nil {
			t.Fatal(err)
		} else {
			return "", "", fmt.Errorf(err.Error())
		}
	}
	outScanner := bufio.NewScanner(outStream)
	// stdout reader
	wg.Add(1)
	go func() {
		for outScanner.Scan() {
			line := outScanner.Text()
			if verbose {
				logf(t, "[stdout] %s", line)
			}
			stdout += line + "\n"
		}
		if err := outScanner.Err(); err != nil {
			fatalf(t, err.Error()) // TODO: remove log.Fatal
		}
		wg.Done()
	}()

	// standard error
	errStream, err := cmd.StderrPipe()
	if err != nil {
		// fatalf(t, err.Error())
		// return
		if t != nil {
			t.Fatal(err)
		} else {
			return "", "", fmt.Errorf(err.Error())
		}
	}
	errScanner := bufio.NewScanner(errStream)
	// stderr reader
	wg.Add(1)
	go func() {
		for errScanner.Scan() {
			line := errScanner.Text()
			if verbose {
				logf(t, "[stderr] %s", line)
			}
			stderr += line + "\n"
		}
		if err := errScanner.Err(); err != nil {
			fatalf(t, err.Error()) // TODO: remove log.Fatal
		}
		wg.Done()
	}()

	// start execution
	if err = cmd.Start(); err != nil {
		// fatalf(t, err.Error())
		// return
		if t != nil {
			t.Fatal(err)
		} else {
			return "", "", fmt.Errorf(err.Error())
		}
	}

	// wait for all standard output processing before waiting to exit!
	wg.Wait()

	// wait until command exits
	if err = cmd.Wait(); err != nil {
		// fatalf(t, err.Error())
		// return
		if t != nil {
			t.Fatal(err)
		} else {
			return stdout, stderr, fmt.Errorf(err.Error())
		}
	}

	return
}
