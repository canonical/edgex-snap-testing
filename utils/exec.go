package utils

import (
	"bufio"
	"io"
	"log"
	"os/exec"
	"sync"
	"testing"
)

// RunCommand executes one or more commands
func RunCommand(t *testing.T, commands ...string) (string, string) {
	var stdout string
	var stderr string

	for _, command := range commands {
		if t != nil {
			t.Logf("Running command: %s", command)
		} else {
			log.Printf("Running command: %s", command)
		}

		cmd := exec.Command("/bin/sh", "-c", command)

		outStream, err := cmd.StdoutPipe()
		if err != nil {
			if t != nil {
				t.Fatal(err)
			} else {
				log.Fatal(err)
			}
			return stdout, stderr
		}

		errStream, err := cmd.StderrPipe()
		if err != nil {
			if t != nil {
				t.Fatal(err)
			} else {
				log.Fatal(err)
			}
			return stdout, stderr
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
				if t != nil {
					t.Logf("stdout: %s", line)
				} else {
					log.Printf("stdout: %s", line)
				}
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
				if t != nil {
					t.Logf("stderr: %s", line)
				} else {
					log.Printf("stderr: %s", line)
				}
				stderr += line + "\n"
			}
			wg.Done()
		}(errStream)

		// start execution
		err = cmd.Start()
		if err != nil {
			if t != nil {
				t.Fatal(err)
			} else {
				log.Fatal(err)
			}
			return stdout, stderr
		}

		// wait until it exits
		err = cmd.Wait()
		if err != nil {
			if t != nil {
				t.Fatal(err)
			} else {
				log.Fatal(err)
			}
			return stdout, stderr
		}

	}
	return stdout, stderr
}
