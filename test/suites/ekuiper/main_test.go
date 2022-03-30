package test

import (
	"edgex-snap-testing/test/utils"
	"log"
	"os"
	"testing"
)

const snap = "edgex-ekuiper"
const snapService = "edgex-ekuiper.kuiper"

func TestMain(m *testing.M) {

	log.Println("[SETUP]")

	// start clean
	utils.SnapRemove(nil,
		snap,
		"edgexfoundry",
	)

	// install the ekuiper snap before edgexfoundry
	// to catch build error sooner and stop
	if utils.LocalSnap != "" {
		utils.SnapInstallFromFile(nil, utils.LocalSnap)
	} else {
		utils.SnapInstallFromStore(nil, snap, utils.ServiceChannel)
	}
	utils.SnapInstallFromStore(nil, "edgexfoundry", utils.PlatformChannel)

	// for local build, the interface isn't auto-connected.
	// connect manually regardless
	utils.SnapConnect(nil,
		"edgexfoundry:edgex-secretstore-token",
		snap+":edgex-secretstore-token",
	)
	utils.SnapRestart(nil,
		snapService,
	)

	// security on (default)
	exitCode := m.Run()
	if exitCode != 0 {
		goto TEARDOWN
	}

	// security off
	utils.SnapSet(nil, "edgexfoundry", "security-secret-store", "off")
	utils.SnapRemove(nil,
		snap)

	if utils.LocalSnap != "" {
		utils.SnapInstallFromFile(nil, utils.LocalSnap)
	} else {
		utils.SnapInstallFromStore(nil, snap, utils.ServiceChannel)
	}

	exitCode = m.Run()

TEARDOWN:
	log.Println("[TEARDOWN]")

	utils.SnapDumpLogs(nil, snap)

	utils.SnapRemove(nil,
		snap,
		"edgexfoundry",
	)

	os.Exit(exitCode)
}
