package test

import (
	"edgex-snap-testing/test/utils"
	"log"
	"os"
	"testing"
	"time"
)

const (
	uiSnap    = "edgex-ui"
	uiApp     = "edgex-ui"
	uiService = uiSnap + "." + uiApp
)

var start = time.Now()

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

	utils.SnapDumpLogs(nil, start, uiSnap)

	utils.SnapRemove(nil,
		uiSnap,
	)

	os.Exit(exitCode)
}
