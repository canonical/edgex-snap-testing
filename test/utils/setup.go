package utils

import (
	"log"
	"time"
)

const platformSnap = "edgexfoundry"

// SetupServiceTests setup up the environment for testing
// It returns a teardown function to be called at the end of the tests
func SetupServiceTests(snapName string) (teardown func(), err error) {
	log.Println("[CLEAN]")
	SnapRemove(nil,
		snapName,
		platformSnap,
	)

	log.Println("[SETUP]")
	start := time.Now()

	teardown = func() {
		log.Println("[TEARDOWN]")
		SnapDumpLogs(nil, start, snapName)
		SnapDumpLogs(nil, start, platformSnap)

		log.Println("Removing installed snap:", !SkipTeardownRemoval)
		if !SkipTeardownRemoval {
			SnapRemove(nil,
				snapName,
				platformSnap,
			)
		}
	}

	// install the device/app service snap before edgexfoundry
	// to catch build error sooner and stop
	if LocalServiceSnap() {
		err = SnapInstallFromFile(nil, LocalServiceSnapPath)
	} else {
		err = SnapInstallFromStore(nil, snapName, ServiceChannel)
	}
	if err != nil {
		teardown()
		return
	}

	if err = SnapInstallFromStore(nil, platformSnap, PlatformChannel); err != nil {
		teardown()
		return
	}

	// install the edgexfoundry platform snap
	if LocalPlatformSnap() {
		err = SnapInstallFromFile(nil, LocalPlatformSnapPath)
	} else {
		err = SnapInstallFromStore(nil, snapName, PlatformChannel)
	}
	if err != nil {
		teardown()
		return
	}

	// for local build, the interface isn't auto-connected.
	// connect manually
	if LocalServiceSnap() || LocalPlatformSnap() {
		if err = SnapConnectSecretstoreToken(nil, snapName); err != nil {
			teardown()
			return
		}
	}

	SnapStart(nil, platformSnap)

	// make sure all services are online before starting the tests
	if err = WaitPlatformOnline(nil); err != nil {
		teardown()
		return
	}

	return
}
