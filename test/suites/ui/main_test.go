package test

import (
	"edgex-snap-testing/test/utils"
	"log"
	"os"
	"testing"
)

const uiSnap = "edgex-ui"
const uiService = "edgex-ui.edgex-ui"

func TestMain(m *testing.M) {

	log.Println("[SETUP]")

	// start clean
	utils.SnapRemove(nil,
		uiSnap,
	)

	if utils.LocalSnap != "" {
		utils.SnapInstallFromFile(nil, utils.LocalSnap)
	} else {
		utils.SnapInstallFromStore(nil, uiSnap, utils.ServiceChannel)
	}

	exitCode := m.Run()

	log.Println("[TEARDOWN]")

	utils.SnapDumpLogs(nil, uiSnap)

	utils.SnapRemove(nil,
		uiSnap,
	)

	os.Exit(exitCode)
}
