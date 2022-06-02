package test

import (
	"edgex-snap-testing/test/utils"
	"log"
	"os"
	"testing"
	"time"
)

const (
	ekuiperSnap       = "edgex-ekuiper"
	ekuiperApp        = "kuiper"
	ekuiperService    = ekuiperSnap + "." + ekuiperApp
	deviceVirtualSnap = "edgex-device-virtual"
)

var start = time.Now()

func TestMain(m *testing.M) {
	log.Println("[SETUP]")

	// start clean
	utils.SnapRemove(nil,
		ekuiperSnap,
		"edgexfoundry",
		deviceVirtualSnap,
	)

	// install the ekuiper snap before edgexfoundry
	// to catch build error sooner and stop
	if utils.LocalSnap != "" {
		utils.SnapInstallFromFile(nil, utils.LocalSnap)
	} else {
		utils.SnapInstallFromStore(nil, ekuiperSnap, utils.ServiceChannel)
	}
	utils.SnapInstallFromStore(nil, "edgexfoundry", utils.PlatformChannel)
	utils.SnapInstallFromStore(nil, deviceVirtualSnap, "latest/edge")

	// make sure all services are online before starting the tests
	utils.WaitPlatformOnline(nil)

	// for local build, the interface isn't auto-connected.
	// connect manually regardless
	utils.SnapConnect(nil,
		"edgexfoundry:edgex-secretstore-token",
		ekuiperSnap+":edgex-secretstore-token",
	)
	utils.SnapConnect(nil,
		"edgexfoundry:edgex-secretstore-token",
		deviceVirtualSnap+":edgex-secretstore-token",
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
	utils.SnapSet(nil, "edgex-device-virtual", "app-options", "true")
	utils.SnapSet(nil, "edgex-device-virtual", "config.edgex-security-secret-store", "false")
	utils.Exec(nil, "sudo rm /var/snap/edgex-ekuiper/current/edgex-ekuiper/secrets-token.json")

	utils.SnapStart(nil,
		ekuiperService,
		deviceVirtualSnap,
	)

	exitCode = m.Run()

TEARDOWN:
	log.Println("[TEARDOWN]")

	utils.SnapDumpLogs(nil, start, ekuiperSnap)

	utils.SnapRemove(nil,
		ekuiperSnap,
		"edgexfoundry",
		deviceVirtualSnap,
	)

	os.Exit(exitCode)
}
