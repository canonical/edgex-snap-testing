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
	// edgex-ekuiper's latest/edge channel is currently broken:
	// https://forum.snapcraft.io/t/snapcraft-release-has-no-effects-for-channel-latest-edge/29069
	if utils.ServiceChannel == "latest/edge" {
		utils.ServiceChannel = "1/edge"
	}

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

	// make sure all services are online before starting the tests
	utils.WaitPlatformOnline(nil)

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
