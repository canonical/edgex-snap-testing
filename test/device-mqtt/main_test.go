package test

import (
	"edgex-snap-testing/env"
	"edgex-snap-testing/utils"
	"log"
	"os"
	"testing"
)

const (
	thisSnap = "edgex-device-mqtt"
)

func TestMain(m *testing.M) {

	log.Println("[GLOBAL SETUP]")

	// start clean
	utils.SnapRemove(nil,
		thisSnap,
		"edgexfoundry",
	)

	// install the device snap before edgexfoundry
	// to catch build error sooner and stop
	if env.Snap == "" {
		utils.SnapInstall(nil, thisSnap, env.Channel)
	} else {
		utils.SnapInstallLocal(nil, env.Snap)
	}

	utils.SnapInstall(nil, "edgexfoundry", env.Channel)

	// for local build, the interface isn't auto-connected.
	// connect manually regardless
	utils.SnapConnect(nil,
		"edgexfoundry:edgex-secretstore-token",
		thisSnap+":edgex-secretstore-token",
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
		thisSnap,
		"edgexfoundry",
	)

	os.Exit(exitCode)
}
