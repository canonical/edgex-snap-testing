package utils

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// RequireServiceOnline dials port(s)to check if the service comes online until it reaches the maximum retry
func RequireServiceOnline(t *testing.T, ports ...string) {
	const dialTimeout = 2 * time.Second
	const maxRetry = 60

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

// RequirePortOpen checks if a port accepts connections
func RequirePortOpen(t *testing.T, host, port string) {
	conn, err := net.Dial("tcp", net.JoinHostPort(host, port))
	if err != nil {
		conn.Close()
		t.Fatalf("Port %s is not open: %s", port, err)
	}

	if conn == nil {
		t.Fatalf("Port %s is not open", port)
	}

	if conn != nil {
		t.Logf("Port %v is open.", port)
		conn.Close()
	}
}

// RequirePortOpenLocalhost checks if a local port accepts connections
func RequirePortOpenLocalhost(t *testing.T, port string) {
	RequirePortOpen(t, "localhost", port)
}

// RequirePortAvailable checks if a port is available, i.e. not used by a server
func RequirePortAvailable(t *testing.T, port string) {
	stdout, _ := Exec(t, fmt.Sprintf("sudo lsof -nPi :%s || true", port))
	if stdout != "" {
		t.Fatalf("Port %s is not available", port)
	}
	t.Logf("Port %s is available.", port)
}

// checkListenAllInterfaces checks if a port listens on all interfaces
func checkListenAllInterfaces(t *testing.T, port string, mustListen bool) {
	stdout, _ := Exec(t, "sudo lsof -nPi :"+port+" | { grep \\* || true; }")
	isListening := (stdout != "")

	if mustListen && !isListening {
		t.Fatalf("Port %v not listening to all interfaces", port)
	} else if !mustListen && isListening {
		t.Fatalf("Port %v is listening to all interfaces", port)
	}
}

// RequireNotListenAllInterfaces checks that a port is NOT listening on all interfaces
func RequireNotListenAllInterfaces(t *testing.T, port string) {
	checkListenAllInterfaces(t, port, false)
}

// RequireListenLoopback checks if a port listens on the loopback interface
func RequireListenLoopback(t *testing.T, port string) {
	stdout, _ := Exec(t, "sudo lsof -nPi :"+port+" | { grep 127.0.0.1  || true; }")
	isListening := stdout != ""

	if !isListening {
		t.Fatalf("Port %v not listening on loopback interface", port)
	}
}
