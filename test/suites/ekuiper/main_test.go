package test

import (
	"edgex-snap-testing/test/utils"
	"log"
	"os"
	"testing"
)

const ekuiperSnap = "edgex-ekuiper"
const ekuiperService = "edgex-ekuiper.kuiper"

func TestMain(m *testing.M) {

	log.Println("[SETUP]")

	// start clean
	utils.SnapRemove(nil,
		ekuiperSnap,
		"edgexfoundry",
	)

	// install the ekuiper snap before edgexfoundry
	// to catch build error sooner and stop
	if utils.LocalSnap != "" {
		utils.SnapInstallFromFile(nil, utils.LocalSnap)
	} else {
		utils.SnapInstallFromStore(nil, ekuiperSnap, utils.ServiceChannel)
	}
	utils.SnapInstallFromStore(nil, "edgexfoundry", utils.PlatformChannel)

	// for local build, the interface isn't auto-connected.
	// connect manually regardless
	utils.SnapConnect(nil,
		"edgexfoundry:edgex-secretstore-token",
		ekuiperSnap+":edgex-secretstore-token",
	)
	utils.SnapRestart(nil,
		ekuiperService,
	)

	// security on (default)
	exitCode := m.Run()
	if exitCode != 0 {
		goto TEARDOWN
	}

	// security off
	utils.Exec(nil, "sudo snap set edgexfoundry security-secret-store=off")
	utils.Exec(nil, "sudo rm /var/snap/edgex-ekuiper/current/edgex-ekuiper/secrets-token.json")
	utils.SnapDisconnect(nil,
		"edgexfoundry:edgex-secretstore-token",
		ekuiperSnap+":edgex-secretstore-token",
	)

	exitCode = m.Run()

TEARDOWN:
	log.Println("[TEARDOWN]")

	utils.SnapDumpLogs(nil, ekuiperSnap)

	utils.SnapRemove(nil,
		ekuiperSnap,
		"edgexfoundry",
	)

	os.Exit(exitCode)
}
