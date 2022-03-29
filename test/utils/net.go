package utils

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// WaitServiceOnline dials port(s)to check if the service comes online until it reaches the maximum retry
func WaitServiceOnline(t *testing.T, ports ...string) {
	const dialTimeout = 2 * time.Second
	const maxRetry = 120

	for _, port := range ports {
		serviceIsOnline := false
		var returnErr error

		for i := 0; !serviceIsOnline && i < maxRetry; i++ {
			t.Logf("Waiting for service. Dialing port %s. Retry %d/%d", port, i+1, maxRetry)
			conn, err := net.DialTimeout("tcp", net.JoinHostPort("127.0.0.1", port), dialTimeout)
			if conn != nil {
				serviceIsOnline = true
				t.Logf("Service online now. Port %s is listening", port)
			}
			returnErr = err
			time.Sleep(1 * time.Second)
		}

		require.Equal(t, true, serviceIsOnline,
			"Service timed out, reached max %d retries. Error:\n%v", maxRetry, returnErr)
	}
}

// PortConnection checks if the port(s) are in use
func PortConnection(host string, ports ...string) (bool, error) {
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
func PortConnectionAllInterface(t *testing.T, ports ...string) bool {
	var isListening bool

	for _, port := range ports {

		stdout, _ := Exec(t, "sudo lsof -nPi :"+port+" | { grep \\* || true; }")
		if stdout == "" {
			isListening = false
		} else {
			isListening = true
		}
	}
	return isListening
}

// PortConnectionLocalhost checks if the port(s) are in use for localhost
func PortConnectionLocalhost(t *testing.T, ports ...string) bool {
	var isOpen bool

	for _, port := range ports {

		stdout, _ := Exec(t, "sudo lsof -nPi :"+port+" | { grep 127.0.0.1  || true; }")
		if stdout == "" {
			isOpen = false
		} else {
			isListening, err := PortConnection("127.0.0.1", port)
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

func CheckPortAvailable(t *testing.T, port string) {
	stdout, _ := Exec(t, fmt.Sprintf("sudo lsof -nPi :%s || true", port))
	if stdout != "" {
		t.Fatalf("Port %s is not available", port)
	}
	t.Logf("Port %s is available.", port)
}
