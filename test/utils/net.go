package utils

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// RequireServiceOnline checks if a service comes online by dialing port(s) the ports for a limited period
func RequireServiceOnline(t *testing.T, ports ...string) {
	const timeout = 2 * time.Second
	const maxRetry = 60

	for _, port := range ports {
		serviceIsOnline := false
		var returnErr error

		for i := 0; !serviceIsOnline && i < maxRetry; i++ {
			t.Logf("Waiting for service. Dialing port %s. Retry %d/%d", port, i+1, maxRetry)
			conn, err := net.DialTimeout("tcp", net.JoinHostPort("127.0.0.1", port), timeout)
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

// RequirePortOpen checks if the local port(s) accepts connections
func RequirePortOpen(t *testing.T, host string, ports ...string) {
	const timeout = 2 * time.Second

	for _, port := range ports {
		conn, err := net.DialTimeout("tcp", net.JoinHostPort("localhost", port), timeout)
		if err != nil {
			conn.Close()
			t.Errorf("Port %s is not open: %s", port, err)
		}

		if conn == nil {
			t.Errorf("Port %s is not open", port)
		}

		if conn != nil {
			t.Logf("Port %v is open.", port)
			conn.Close()
		}
	}
	if t.Failed() {
		t.FailNow()
	}
}

// checkListenAllInterfaces checks if the port(s) listen on all interfaces
func RequireListenAllInterfaces(t *testing.T, mustListen bool, ports ...string) {
	for _, port := range ports {
		stdout, _ := Exec(t, "sudo lsof -nPi :"+port+" | { grep \\* || true; }")
		isListening := (stdout != "")

		if mustListen && !isListening {
			t.Errorf("Port %v not listening to all interfaces", port)
		} else if !mustListen && isListening {
			t.Errorf("Port %v is listening to all interfaces", port)
		}
	}
	if t.Failed() {
		t.FailNow()
	}
}

// RequireListenLoopback checks if the port(s) listen on the loopback interface
func RequireListenLoopback(t *testing.T, ports ...string) {
	for _, port := range ports {
		stdout, _ := Exec(t, "sudo lsof -nPi :"+port+" | { grep 127.0.0.1  || true; }")
		isListening := stdout != ""

		if !isListening {
			t.Errorf("Port %v not listening on loopback interface", port)
		}
	}
	if t.Failed() {
		t.FailNow()
	}
}

// RequirePortAvailable checks if a port is available (not open) locally
func RequirePortAvailable(t *testing.T, port string) {
	stdout, _ := Exec(t, fmt.Sprintf("sudo lsof -nPi :%s || true", port))
	if stdout != "" {
		t.Fatalf("Port %s is not available", port)
	}
	t.Logf("Port %s is available.", port)
}
