package test

import (
	"edgex-snap-testing/test/utils"
	"log"
	"os"
	"testing"
)

const cliSnap = "edgex-cli"

func TestMain(m *testing.M) {

	log.Println("[SETUP]")

	// start clean
	utils.SnapRemove(nil,
		cliSnap,
	)

	if utils.LocalSnap != "" {
		utils.SnapInstallFromFile(nil, utils.LocalSnap)
	} else {
		utils.SnapInstallFromStore(nil, cliSnap, utils.ServiceChannel)
	}

	exitCode := m.Run()

	log.Println("[TEARDOWN]")

	utils.SnapRemove(nil,
		cliSnap,
	)

	os.Exit(exitCode)
}
