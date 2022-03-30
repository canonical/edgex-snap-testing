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
		"edgex-ekuiper",
		"edgexfoundry",
	)

	// install the ekuiper snap before edgexfoundry
	// to catch build error sooner and stop
	if utils.LocalSnap != "" {
		utils.SnapInstallFromFile(nil, utils.LocalSnap)
	} else {
		utils.SnapInstallFromStore(nil, "edgex-ekuiper", utils.ServiceChannel)
	}
	utils.SnapInstallFromStore(nil, "edgexfoundry", utils.PlatformChannel)

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
	utils.SnapSet(nil, "edgexfoundry", "security-secret-store", "off")
	utils.SnapRemove(nil,
		"edgex-ekuiper")

	if utils.LocalSnap != "" {
		utils.SnapInstallFromFile(nil, utils.LocalSnap)
	} else {
		utils.SnapInstallFromStore(nil, "edgex-ekuiper", utils.ServiceChannel)
	}

	exitCode = m.Run()

TEARDOWN:
	log.Println("[TEARDOWN]")

	utils.SnapDumpLogs(nil, "edgex-ekuiper")

	utils.SnapRemove(nil,
		"edgex-ekuiper",
		"edgexfoundry",
	)

	os.Exit(exitCode)
}
