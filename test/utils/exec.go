package utils

import (
	"bufio"
	"os/exec"
	"sync"
	"testing"
)

// Exec executes a command
func Exec(t *testing.T, command string) (stdout, stderr string) {
	logf(t, "[exec] %s", command)

	cmd := exec.Command("/bin/sh", "-c", command)

	var wg sync.WaitGroup

	// standard output
	outStream, err := cmd.StdoutPipe()
	if err != nil {
		fatalf(t, err.Error())
		return
	}
	outScanner := bufio.NewScanner(outStream)
	// stdout reader
	wg.Add(1)
	go func() {
		for outScanner.Scan() {
			line := outScanner.Text()
			logf(t, "[stdout] %s", line)
			stdout += line + "\n"
		}
		if err := outScanner.Err(); err != nil {
			fatalf(t, err.Error())
		}
		wg.Done()
	}()

	// standard error
	errStream, err := cmd.StderrPipe()
	if err != nil {
		fatalf(t, err.Error())
		return
	}
	errScanner := bufio.NewScanner(errStream)
	// stderr reader
	wg.Add(1)
	go func() {
		for errScanner.Scan() {
			line := errScanner.Text()
			logf(t, "[stderr] %s", line)
			stderr += line + "\n"
		}
		if err := errScanner.Err(); err != nil {
			fatalf(t, err.Error())
		}
		wg.Done()
	}()

	// start execution
	if err := cmd.Start(); err != nil {
		fatalf(t, err.Error())
		return
	}

	// wait for all standard output processing before waiting to exit!
	wg.Wait()

	// wait until command exits
	if err := cmd.Wait(); err != nil {
		fatalf(t, err.Error())
		return
	}

	return
}
