package test

import (
	"edgex-snap-testing/env"
	"edgex-snap-testing/utils"
	"log"
	"os"
	"testing"
)

// This variable should be removed once 1.4.3 edegx-ekuiper can be promoted to the latest/edge channel
var temporaryChannel = "latest/beta"

func TestMain(m *testing.M) {

	log.Println("[GLOBAL SETUP]")

	// start clean
	utils.SnapRemove(nil,
		"edgex-ekuiper",
		"edgexfoundry",
	)

	// install the ekuiper snap before edgexfoundry
	// to catch build error sooner and stop
	if env.Snap == "" {

		// utils.SnapInstall(nil, "edgex-ekuiper", env.Channel)
		utils.SnapInstall(nil, "edgex-ekuiper", temporaryChannel)
	} else {
		utils.SnapInstallLocal(nil, env.Snap)
	}

	utils.SnapInstall(nil, "edgexfoundry", env.Channel)

	// for local build, the interface isn't auto-connected.
	// connect manually regardless
	utils.SnapConnect(nil,
		"edgexfoundry:edgex-secretstore-token",
		"edgex-ekuiper:edgex-secretstore-token",
	)
	utils.Exec(nil,
		"sudo snap restart edgex-ekuiper.kuiper",
	)

	// security on (default)
	exitCode := m.Run()
	if exitCode != 0 {
		goto TEARDOWN
	}

	// security off
	utils.Exec(nil, "sudo snap set edgexfoundry security-secret-store=off")
	utils.SnapRemove(nil,
		"edgex-ekuiper")
	if env.Snap == "" {
		utils.SnapInstall(nil, "edgex-ekuiper", env.ekuiperChannel)
	} else {
		utils.SnapInstallLocal(nil, env.Snap)
	}

	exitCode = m.Run()

TEARDOWN:
	log.Println("[GLOBAL TEARDOWN]")

	// TODO: should the logs be fetched in each test?
	// for that, need to use journalctl instead with --since
	if exitCode != 0 {
		stdout, _ := utils.Exec(nil,
			"sudo snap logs -n=all edgex-ekuiper")
		log.Printf("Snap logs:\n%s\n", stdout)
	}

	utils.SnapRemove(nil,
		"edgex-ekuiper",
		"edgexfoundry",
	)

	os.Exit(exitCode)
}
