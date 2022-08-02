package test

import (
	"edgex-snap-testing/test/utils"
	"log"
	"os"
	"testing"
)

const cliSnap = "edgex-cli"

func TestMain(m *testing.M) {
	// start clean
	utils.SnapRemove(nil,
		cliSnap,
	)

	log.Println("[SETUP]")

	if utils.LocalSnap() {
		utils.SnapInstallFromFile(nil, utils.LocalSnapPath)
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

func TestCommon(t *testing.T) {

	utils.TestPackaging(t, cliSnap, utils.Packaging{
		TestSemanticSnapVersion: true,
	})
}
