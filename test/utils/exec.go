package utils

import (
	"bufio"
	"io"
	"log"
	goexec "os/exec"
	"sync"
	"testing"
)

func Exec(t *testing.T, command string) (stdout, stderr string, err error) {
	return exec(t, command, false)
}

// exec executes a command
func exec(t *testing.T, command string, verbose bool) (stdout, stderr string, err error) {
	if t != nil {
		t.Logf("[exec] %s", command)
	} else {
		log.Printf("[exec] %s", command)
	}

	cmd := goexec.Command("/bin/sh", "-c", command)

	var wg sync.WaitGroup

	// standard output
	outStream, err := cmd.StdoutPipe()
	if err != nil {
		if t != nil {
			t.Fatal(err)
		} else {
			return "", "", err
		}
	}
	wg.Add(1)
	go scanStdPipe(t, outStream, &stdout, &wg, verbose, false)

	// standard error
	errStream, err := cmd.StderrPipe()
	if err != nil {
		if t != nil {
			t.Fatal(err)
		} else {
			return "", "", err
		}
	}
	wg.Add(1)
	go scanStdPipe(t, errStream, &stderr, &wg, verbose, true)

	// start execution
	if err = cmd.Start(); err != nil {
		if t != nil {
			t.Fatal(err)
		} else {
			return stdout, stderr, err
		}
	}

	// wait for all standard output processing before waiting to exit!
	wg.Wait()

	// wait until command exits
	if err = cmd.Wait(); err != nil {
		if t != nil {
			t.Fatal(err)
		} else {
			return stdout, stderr, err
		}
	}

	return
}

// scan and process the standard output / error streams
func scanStdPipe(t *testing.T, stream io.Reader, streamStr *string, wg *sync.WaitGroup, verbose, stderr bool) {
	defer wg.Done()

	var prefix string
	if stderr {
		prefix = "[stderr]"
	} else {
		prefix = "[stdout]"
	}

	scanner := bufio.NewScanner(stream)

	for scanner.Scan() {
		line := scanner.Text()
		if verbose {
			if t != nil {
				t.Logf("%s %s", prefix, line)
			} else {
				log.Printf("%s %s", prefix, line)
			}
		}
		*streamStr += line + "\n"
	}
	if err := scanner.Err(); err != nil {
		if t != nil {
			t.Error(err)
		} else {
			log.Printf("Scanner error: %s", err)
		}
	}
}
