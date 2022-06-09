package utils

import (
	"fmt"
	"net"
	"strings"
	"testing"
	"time"
)

const dialTimeout = 2 * time.Second

var PlatformPorts = []string{
	"59880", // core-data
	"59881", // core-metadata
	"59882", // core-command
	"8000",  // kong
	"5432",  // kong-database
	"8200",  // vault
	"8500",  // consul
	"6379",  // redis
}

func TestNet(t *testing.T, snapName string, conf Net) {
	t.Run("net", func(t *testing.T) {
		if conf.StartSnap {
			t.Cleanup(func() {
				SnapStop(t, snapName)
			})
			SnapStart(t, snapName)
		}

		if len(conf.TestOpenPorts) > 0 {
			TestOpenPorts(t, snapName, conf.TestOpenPorts)
		}
		if len(conf.TestBindLoopback) > 0 {
			TestBindLoopback(t, snapName, conf.TestBindLoopback)
		}

	})
}

func TestOpenPorts(t *testing.T, snapName string, ports []string) {
	t.Run("ports open", func(t *testing.T) {
		WaitServiceOnline(t, 60, ports...)
	})
}

func TestBindLoopback(t *testing.T, snapName string, ports []string) {
	WaitServiceOnline(t, 60, ports...)

	t.Run("ports not listening on all interfaces", func(t *testing.T) {
		RequireListenAllInterfaces(t, false, ports...)
	})

	t.Run("ports listening on localhost", func(t *testing.T) {
		RequireListenLoopback(t, ports...)
		// RequirePortOpen(t, params.TestBindAddrLoopback...)
	})
}

// WaitServiceOnline waits for a service to come online by dialing its port(s)
// up to a maximum number
func WaitServiceOnline(t *testing.T, maxRetry int, ports ...string) {
	if len(ports) == 0 {
		panic("No ports given as input")
	}

PORTS:
	for _, port := range ports {
		var returnErr error

		for i := 1; i <= maxRetry; i++ {
			logf(t, "Waiting for service port %s. Retry %d/%d", port, i, maxRetry)

			conn, err := net.DialTimeout("tcp", ":"+port, dialTimeout)
			if conn != nil {
				logf(t, "Service port %s is open.", port)
				continue PORTS
			}
			returnErr = err

			time.Sleep(1 * time.Second)
		}

		if returnErr != nil {
			fatalf(t, "Time out: reached max %d retries. Error: %v", maxRetry, returnErr)
		} else {
			fatalf(t, "Time out: reached max %d retries.", maxRetry)
		}
	}
}

// WaitPlatformOnline waits for all platform ports to come online
// by dialing its port(s) up to a maximum number
func WaitPlatformOnline(t *testing.T) {
	WaitServiceOnline(t, 180, PlatformPorts...)
}

// RequirePortOpen checks if the local port(s) accepts connections
func RequirePortOpen(t *testing.T, ports ...string) {
	if len(ports) == 0 {
		panic("No ports given as input")
	}

	for _, port := range ports {
		conn, err := net.DialTimeout("tcp", ":"+port, dialTimeout)
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
	if len(ports) == 0 {
		panic("No ports given as input")
	}

	for _, port := range ports {
		isListening := isListenInterface(t, "*", port)

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
// It don't not check whether port(s) listen on interfaces other than the loopback
func RequireListenLoopback(t *testing.T, ports ...string) {
	if len(ports) == 0 {
		panic("No ports given as input")
	}

	for _, port := range ports {
		if !isListenInterface(t, "127.0.0.1", port) {
			t.Errorf("Port %s is not restricted to listen on loopback interface", port)
		}
	}
	if t.Failed() {
		t.FailNow()
	}
}

// RequirePortAvailable checks if a port is available (not open) locally
func RequirePortAvailable(t *testing.T, port string) {
	stdout := lsof(t, port)
	if stdout != "" {
		t.Fatalf("Port %s is not available", port)
	}
	t.Logf("Port %s is available.", port)
}

func isListenInterface(t *testing.T, addr string, port string) bool {
	list := filterOpenPorts(t, port)

	// look for LISTEN explicitly to exclude ESTABLISHED connections
	substr := fmt.Sprintf("%s:%s (LISTEN)", addr, port)
	t.Logf("Looking for '%s'", substr)

	return strings.Contains(list, substr)
}

func filterOpenPorts(t *testing.T, port string) string {
	stdout := lsof(t, port)
	if stdout == "" {
		t.Fatalf("Port %s is not open", port)
	}
	return stdout
}

func lsof(t *testing.T, port string) string {
	// The chained true command is to make sure execution succeeds even if
	// 	the first command fails when list is empty
	stdout, _ := Exec(t, fmt.Sprintf("sudo lsof -nPi :%s || true", port))
	return stdout
}
