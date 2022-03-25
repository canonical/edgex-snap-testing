package test

import (
	"edgex-snap-testing/test/utils"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	log.Println("[SETUP]")

	// start clean
	utils.SnapRemove(nil,
		"edgex-device-mqtt",
		"edgexfoundry",
	)

	// install the device snap before edgexfoundry
	// to catch build error sooner and stop
	if utils.LocalSnap != "" {
		utils.SnapInstallFromFile(nil, utils.LocalSnap)
	} else {
		utils.SnapInstallFromStore(nil, "edgex-device-mqtt", utils.ServiceChannel)
	}
	utils.SnapInstallFromStore(nil, "edgexfoundry", utils.PlatformChannel)

	// for local build, the interface isn't auto-connected.
	// connect manually regardless
	utils.SnapConnect(nil,
		"edgexfoundry:edgex-secretstore-token",
		"edgex-device-mqtt:edgex-secretstore-token",
	)

	exitCode := m.Run()

	log.Println("[TEARDOWN]")

	utils.SnapDumpLogs(nil, "edgex-device-mqtt")

	utils.SnapRemove(nil,
		"edgex-device-mqtt",
		"edgexfoundry",
	)

	os.Exit(exitCode)
}
