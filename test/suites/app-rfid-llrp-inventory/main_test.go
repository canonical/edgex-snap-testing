package test

import (
	"edgex-snap-testing/test/utils"
	"log"
	"os"
	"testing"
)

const appRfidLlrpSnap = "edgex-app-rfid-llrp-inventory"
const appRfidLlrpService = "edgex-app-rfid-llrp-inventory.app-rfid-llrp-inventory"

func TestMain(m *testing.M) {

	log.Println("[SETUP]")

	// start clean
	utils.SnapRemove(nil,
		appRfidLlrpSnap,
		"edgexfoundry",
	)

	// install the app snap before edgexfoundry
	// to catch build error sooner and stop
	if utils.LocalSnap != "" {
		utils.SnapInstallFromFile(nil, utils.LocalSnap)
	} else {
		utils.SnapInstallFromStore(nil, appRfidLlrpSnap, utils.ServiceChannel)
	}
	utils.SnapInstallFromStore(nil, "edgexfoundry", utils.PlatformChannel)

	// make sure all services are online before starting the tests
	utils.WaitPlatformOnline(nil)

	// for local build, the interface isn't auto-connected.
	// connect manually regardless
	utils.SnapConnect(nil,
		"edgexfoundry:edgex-secretstore-token",
		appRfidLlrpSnap+":edgex-secretstore-token",
	)

	// Start the service so that the default config gets uploaded to consul.
	// Otherwise, settings that get passed using environment variables on first start get uploaded
	// and become the default. This is expected behavior.
	utils.SnapStart(nil, appRfidLlrpService)
	utils.WaitServiceOnline(nil, 60, defaultServicePort)

	exitCode := m.Run()

	log.Println("[TEARDOWN]")

	utils.SnapDumpLogs(nil, appRfidLlrpSnap)

	utils.SnapRemove(nil,
		appRfidLlrpSnap,
		"edgexfoundry",
	)

	os.Exit(exitCode)
}
