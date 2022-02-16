package utils

import (
	"bufio"
	"io"
	"log"
	"os/exec"
	"sync"
	"testing"
)

var testingFatal = false

// Exec executes one or more commands
func Exec(t *testing.T, commands ...string) (stdout, stderr string) {

	for _, command := range commands {
		logf(t, "[exec] %s", command)

		cmd := exec.Command("/bin/sh", "-c", command)

		outStream, err := cmd.StdoutPipe()
		if err != nil {
			fatalf(t, err.Error())
			return
		}

		errStream, err := cmd.StderrPipe()
		if err != nil {
			fatalf(t, err.Error())
			return
		}

		var wg sync.WaitGroup
		// wait for all standard output processing
		defer wg.Wait()

		// stdout reader
		wg.Add(1)
		go func(stream io.ReadCloser) {
			scanner := bufio.NewScanner(stream)
			for scanner.Scan() {
				line := scanner.Text()
				logf(t, "[stdout] %s", line)
				stdout += line + "\n"
			}
			wg.Done()
		}(outStream)

		// stderr reader
		wg.Add(1)
		go func(stream io.ReadCloser) {
			scanner := bufio.NewScanner(stream)
			for scanner.Scan() {
				line := scanner.Text()
				logf(t, "[stderr] %s", line)
				stderr += line + "\n"
			}
			wg.Done()
		}(errStream)

		// start execution
		err = cmd.Start()
		if err != nil {
			fatalf(t, err.Error())
			return
		}

		// wait until it exits
		err = cmd.Wait()
		if err != nil {
			fatalf(t, err.Error())
			return
		}

	}
	return
}

func logf(t *testing.T, format string, args ...interface{}) {
	if t != nil {
		t.Logf(format, args...)
	} else {
		log.Printf(format, args...)
	}
}

func fatalf(t *testing.T, format string, args ...interface{}) {
	// reduce the severity to a log message to not exit when testing
	if testingFatal {
		logf(t, "fatal error: "+format, args...)
		return
	}

	if t != nil {
		t.Fatalf(format, args...)
	} else {
		log.Fatalf(format, args...)
	}
}
