package utils

import (
	"bufio"
	"io"
	"log"
	"net"
	"os/exec"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
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

// WaitServiceOnline dials port(s)to check if the service comes online until it reaches the maximum retry
func WaitServiceOnline(t *testing.T, ports []string) error {
	const dialTimeout = 2 * time.Second

	for _, port := range ports {
		serviceIsOnline := false
		var returnErr error

		for i := 0; !serviceIsOnline && i < 60; i++ {
			conn, err := net.DialTimeout("tcp", net.JoinHostPort("127.0.0.1", port), dialTimeout)
			time.Sleep(1 * time.Second)
			t.Logf("Waiting for the service to come online. Current retry count:  %d /60", i+1)
			if conn != nil {
				serviceIsOnline = true
				t.Logf("Service online now. Port %s is listening", port)
			}
			returnErr = err
		}

		require.Equal(t, true, serviceIsOnline, "Service timed out, reached max retry count of 60\n", returnErr)
		if t.Failed() {
			return returnErr
		}
	}
	return nil
}

// PortConnection checks if the port(s) are in use
func PortConnection(host string, ports []string) (bool, error) {
	var isListening bool
	var err error

	for _, port := range ports {

		conn, err := net.Dial("tcp", net.JoinHostPort(host, port))
		if err != nil {
			isListening = false
			conn.Close()
			return isListening, err
		}

		if conn != nil {
			isListening = true
			conn.Close()
			return isListening, err
		}
	}
	return isListening, err
}

// PortConnectionAllInterface checks if the port(s) are in use for all interface
func PortConnectionAllInterface(t *testing.T, ports []string) bool {
	var isListening bool

	for _, port := range ports {

		stdout, _ := RunCommand(t, "sudo lsof -nPi :"+port+" | { grep \\* || true; }")
		if stdout == "" {
			isListening = false
		} else {
			isListening = true
		}
	}
	return isListening
}

// PortConnectionLocalhost checks if the port(s) are in use for localhost
func PortConnectionLocalhost(t *testing.T, ports []string) bool {
	var isOpen bool

	for _, port := range ports {

		stdout, _ := RunCommand(t, "sudo lsof -nPi :"+port+" | { grep 127.0.0.1  || true; }")
		if stdout == "" {
			isOpen = false
		} else {
			isListening, err := PortConnection("127.0.0.1", []string{port})
			if isListening {
				isOpen = true
			} else {
				isOpen = false
				require.NoError(t, err, "Error in bind-address or bind-port.")
			}
		}
	}
	return isOpen
}
