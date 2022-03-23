package test

import (
	"edgex-snap-testing/test/utils"
	"edgex-snap-testing/test/utils/env"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	log.Println("[GLOBAL SETUP]")

	// start clean
	utils.SnapRemove(nil,
		"edgex-device-mqtt",
		"edgexfoundry",
	)

	// install the device snap before edgexfoundry
	// to catch build error sooner and stop
	if env.Snap == "" {
		utils.SnapInstall(nil, "edgex-device-mqtt", env.Channel)
	} else {
		utils.SnapInstallLocal(nil, env.Snap)
	}

	utils.SnapInstall(nil, "edgexfoundry", env.Channel)

	// for local build, the interface isn't auto-connected.
	// connect manually regardless
	utils.SnapConnect(nil,
		"edgexfoundry:edgex-secretstore-token",
		"edgex-device-mqtt:edgex-secretstore-token",
	)

	exitCode := m.Run()

	log.Println("[GLOBAL TEARDOWN]")

	if exitCode != 0 {
		log.Printf("Snap logs:\n%s\n",
			utils.SnapLogs(nil, "edgex-device-mqtt"))
	}

	utils.SnapRemove(nil,
		"edgex-device-mqtt",
		"edgexfoundry",
	)

	os.Exit(exitCode)
}
