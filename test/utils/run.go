package utils

import (
	"log"
	"testing"
	"time"
)

func RunDeviceTests(m *testing.M, snapName string) (int, error) {
	log.Println("[CLEAN]")
	SnapRemove(nil,
		snapName,
		"edgexfoundry",
	)

	log.Println("[SETUP]")

	// add this to the bottom of the defer stack to remove after collecting logs
	defer SnapRemove(nil,
		snapName,
		"edgexfoundry",
	)

	start := time.Now()
	defer SnapDumpLogs(nil, start, snapName)
	defer SnapDumpLogs(nil, start, "edgexfoundry")

	// install the device snap before edgexfoundry
	// to catch build error sooner and stop
	if LocalSnap() {
		if err := SnapInstallFromFile(nil, LocalSnapPath); err != nil {
			return 0, err
		}
	} else {
		if err := SnapInstallFromStore(nil, snapName, ServiceChannel); err != nil {
			return 0, err
		}
	}

	if err := SnapInstallFromStore(nil, "edgexfoundry", PlatformChannel); err != nil {
		return 0, err
	}

	// make sure all services are online before starting the tests
	if err := WaitPlatformOnline(nil); err != nil {
		return 0, err
	}

	// for local build, the interface isn't auto-connected.
	// connect manually
	if LocalSnap() {
		if err := SnapConnectSecretstoreToken(nil, snapName); err != nil {
			return 0, err
		}
	}

	return m.Run(), nil
}
