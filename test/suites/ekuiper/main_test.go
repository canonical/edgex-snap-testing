package test

import (
	"edgex-snap-testing/test/utils"
	"log"
	"os"
	"testing"
	"time"
)

const (
	ekuiperSnap    = "edgex-ekuiper"
	ekuiperApp     = "kuiper"
	ekuiperService = ekuiperSnap + "." + ekuiperApp
)

var start = time.Now()

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

	// make sure all services are online before starting the tests
	utils.WaitPlatformOnline(nil)

	// for local build, the interface isn't auto-connected.
	// connect manually regardless
	utils.SnapConnect(nil,
		"edgexfoundry:edgex-secretstore-token",
		ekuiperSnap+":edgex-secretstore-token",
	)

	// security on (default)
	exitCode := m.Run()
	if exitCode != 0 {
		goto TEARDOWN
	}

	// security off
	utils.SnapStop(nil, "edgex-ekuiper")
	utils.SnapSet(nil, "edgexfoundry", "security-secret-store", "off")
	utils.SnapSet(nil, "edgex-ekuiper", "edgex-security", "off")
	utils.Exec(nil, "sudo rm /var/snap/edgex-ekuiper/current/edgex-ekuiper/secrets-token.json")
	utils.SnapDisconnect(nil,
		"edgexfoundry:edgex-secretstore-token",
		ekuiperSnap+":edgex-secretstore-token",
	)
	utils.SnapStart(nil,
		ekuiperService,
	)

	exitCode = m.Run()

TEARDOWN:
	log.Println("[TEARDOWN]")

	utils.SnapDumpLogs(nil, start, ekuiperSnap)

	utils.SnapRemove(nil,
		ekuiperSnap,
		"edgexfoundry",
	)

	os.Exit(exitCode)
}
