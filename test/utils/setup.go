package utils

import (
	"log"
	"time"
)

// SetupServiceTests setup up the environment for testing
// It returns a teardown function to be called at the end of the tests
func SetupServiceTests(snapName string) (teardown func(), err error) {
	log.Println("[CLEAN]")
	SnapRemove(nil,
		snapName,
		"edgexfoundry",
	)

	log.Println("[SETUP]")
	start := time.Now()

	teardown = func() {
		log.Println("[TEARDOWN]")
		SnapDumpLogs(nil, start, snapName)
		SnapDumpLogs(nil, start, "edgexfoundry")

		log.Println("Removing installed snap:", !SkipTeardownRemoval)
		if !SkipTeardownRemoval {
			SnapRemove(nil,
				snapName,
				"edgexfoundry",
			)
		}
	}

	// install the device/app service snap before edgexfoundry
	// to catch build error sooner and stop
	if LocalSnap() {
		err = SnapInstallFromFile(nil, LocalSnapPath)
	} else {
		err = SnapInstallFromStore(nil, snapName, ServiceChannel)
	}
	if err != nil {
		teardown()
		return
	}

	if err = SnapInstallFromStore(nil, "edgexfoundry", PlatformChannel); err != nil {
		teardown()
		return
	}

	// for local build, the interface isn't auto-connected.
	// connect manually
	if LocalSnap() {
		if err = SnapConnectSecretstoreToken(nil, snapName); err != nil {
			teardown()
			return
		}
	}

	// make sure all services are online before starting the tests
	if err = WaitPlatformOnline(nil); err != nil {
		teardown()
		return
	}

	return
}
