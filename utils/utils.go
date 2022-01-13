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

// Command executes an array of string(s)
func Command(commands ...string) (string, string, error) {
	var stdout string
	var stderr string
	var err error

	for _, command := range commands {
		cmd := exec.Command("/bin/sh", "-c", command)

		outStream, err := cmd.StdoutPipe()
		if err != nil {
			return stdout, stderr, err
		}

		errStream, err := cmd.StderrPipe()
		if err != nil {
			return stdout, stderr, err
		}

		var wg sync.WaitGroup

		// stdout reader
		wg.Add(1)
		go func(stream io.ReadCloser) {
			scanner := bufio.NewScanner(stream)
			for scanner.Scan() {
				line := scanner.Text()
				stdout += line
				stdout += "\n"
			}
			wg.Done()
		}(outStream)

		wg.Add(1)
		go func(stream io.ReadCloser) {
			scanner := bufio.NewScanner(stream)
			for scanner.Scan() {
				line := scanner.Text()
				stderr = line
				stderr += "\n"
			}
			wg.Done()
		}(errStream)

		// stderr reader
		err = cmd.Start()
		if err != nil {
			return stdout, stderr, err
		}
		wg.Wait()
	}
	return stdout, stderr, err
}

// CommandLog logs output and err of commands, it is used together with Command
func CommandLog(t *testing.T, stdout string, stderr string, err error) {

	// caller passes t *testing.T
	if t != nil {
		if stdout != "" {
			t.Log(stdout)
		}
		if stderr != "" {
			if err != nil {
				t.Error(stderr)
			} else {
				t.Log(stdout)
			}

		}
		if err != nil {
			t.Fatal(err)
		}
	} else {
		// caller does not passes t *testing.T
		if stdout != "" {
			log.Println(stdout)
		}
		if stderr != "" {
			if err != nil {
				log.Panicln(stderr)
			} else {
				log.Println(stdout)
			}
			if err != nil {
				log.Fatalln(err)
			}
		}
	}
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

		stdout, stderr, err := Command("sudo lsof -nPi :" + port + " | { grep \\* || true; }")
		CommandLog(t, stdout, stderr, err)
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

		stdout, stderr, err := Command("sudo lsof -nPi :" + port + " | { grep 127.0.0.1  || true; }")
		CommandLog(t, stdout, stderr, err)
		if stdout == "" {
			isOpen = false
		} else {
			isListening, err := PortConnection("127.0.0.1", []string{port})
			if isListening {
				isOpen = true
			} else {
				isOpen = false
				CommandLog(t, "", "", err)
			}
		}
	}
	return isOpen
}
