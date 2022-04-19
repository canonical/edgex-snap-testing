package test

import (
	"edgex-snap-testing/test/utils"
	"log"
	"os"
	"testing"
	"time"
)

const (
	platformSnap     = "edgexfoundry"
	deviceVirtualApp = "device-virtual"
	snapAppName      = platformSnap + "." + deviceVirtualApp

	deviceVirtualDefaultServicePort = "59900"
)

var start = time.Now()

func TestMain(m *testing.M) {

	log.Println("[SETUP]")

	// start clean
	utils.SnapRemove(nil,
		platformSnap,
	)

	utils.SnapInstallFromStore(nil, platformSnap, utils.PlatformChannel)

	// make sure all services are online before starting the tests
	utils.WaitPlatformOnline(nil)

	// make sure device-virtual service starts and comes online before starting the tests
	utils.SnapSet(nil, platformSnap, deviceVirtualApp, "on")
	utils.WaitServiceOnline(nil, 60, deviceVirtualDefaultServicePort)

	exitCode := m.Run()

	log.Println("[TEARDOWN]")

	utils.SnapDumpLogs(nil, start, platformSnap)

	utils.SnapRemove(nil,
		platformSnap,
	)

	os.Exit(exitCode)
}
