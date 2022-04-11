package test

import (
	"edgex-snap-testing/test/utils"
	"log"
	"os"
	"testing"
)

const (
	platformSnap     = "edgexfoundry"
	deviceVirtualApp = "device-virtual"
	snapAppName      = platformSnap + "." + deviceVirtualApp

	deviceVirtualDefaultServicePort = "59900"
)

func TestMain(m *testing.M) {

	log.Println("[SETUP]")

	// start clean
	utils.SnapRemove(nil,
		platformSnap,
	)

	utils.SnapInstallFromStore(nil, platformSnap, utils.PlatformChannel)

	// make sure all services are online before starting the tests
	utils.WaitPlatformOnline(nil)

	utils.SnapSet(nil, platformSnap, deviceVirtualApp, "on")
	utils.WaitServiceOnline(nil, 60, deviceVirtualDefaultServicePort)

	exitCode := m.Run()

	log.Println("[TEARDOWN]")

	utils.SnapDumpLogs(nil, platformSnap)

	utils.SnapRemove(nil,
		platformSnap,
	)

	os.Exit(exitCode)
}
