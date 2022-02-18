package test

import (
	"edgex-snap-testing/env"
	"edgex-snap-testing/utils"
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
	if env.SnapcraftProjectDir == "" {
		utils.SnapInstall(nil, "edgex-device-mqtt")
	} else {
		utils.SnapBuild(nil, env.SnapcraftProjectDir)
		utils.SnapInstallLocal(nil, env.SnapcraftProjectDir)
	}

	utils.SnapInstall(nil, "edgexfoundry")

	// for local build, the interface isn't auto-connected.
	// connect manually regardless
	utils.SnapConnect(nil,
		"edgexfoundry:edgex-secretstore-token",
		"edgex-device-mqtt:edgex-secretstore-token",
	)

	exitCode := m.Run()

	log.Println("[GLOBAL TEARDOWN]")

	// TODO: should the logs be fetched in each test?
	// for that, need to use journalctl instead with --since
	if exitCode != 0 {
		stdout, _ := utils.Exec(nil,
			"sudo snap logs -n=all edgex-device-mqtt")
		log.Printf("Snap logs:\n%s\n", stdout)
	}

	utils.SnapRemove(nil,
		"edgex-device-mqtt",
		"edgexfoundry",
	)

	os.Exit(exitCode)
}
