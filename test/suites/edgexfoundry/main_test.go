package test

import (
	"edgex-snap-testing/test/utils"
	"log"
	"os"
	"testing"
)

const platformSnap = "edgexfoundry"
const deviceVirtualApp = "device-virtual"
const deviceVirtualDefaultServicePort = "59900"

var snapAppName = platformSnap + "." + deviceVirtualApp

func TestMain(m *testing.M) {

	log.Println("[SETUP]")

	// start clean
	utils.SnapRemove(nil,
		platformSnap,
	)

	utils.SnapInstallFromStore(nil, platformSnap, utils.PlatformChannel)

	// make sure all services are online before starting the tests
	utils.WaitServiceOnline(nil, utils.PlatformPorts...)
	utils.SnapSet(nil, platformSnap, deviceVirtualApp, "on")
	utils.WaitServiceOnline(nil, deviceVirtualDefaultServicePort)

	exitCode := m.Run()

	log.Println("[TEARDOWN]")

	utils.SnapDumpLogs(nil, platformSnap)

	utils.SnapRemove(nil,
		platformSnap,
	)

	os.Exit(exitCode)
}
