package test

import (
	"edgex-snap-testing/test/utils"
	"log"
	"os"
	"testing"
)

const ascSnap = "edgex-app-service-configurable"
const ascService = "edgex-app-service-configurable.app-service-configurable"

func TestMain(m *testing.M) {

	log.Println("[SETUP]")

	// start clean
	utils.SnapRemove(nil,
		ascSnap,
		"edgexfoundry",
	)

	// install the app-service-configurable snap before edgexfoundry
	// to catch build error sooner and stop
	if utils.LocalSnap != "" {
		utils.SnapInstallFromFile(nil, utils.LocalSnap)
	} else {
		utils.SnapInstallFromStore(nil, ascSnap, utils.ServiceChannel)
	}
	utils.SnapInstallFromStore(nil, "edgexfoundry", utils.PlatformChannel)

	// for local build, the interface isn't auto-connected.
	// connect manually regardless
	utils.SnapConnect(nil,
		"edgexfoundry:edgex-secretstore-token",
		ascSnap+":edgex-secretstore-token",
	)

	// set profile to rules engine
	utils.SnapSet(nil, ascSnap, "profile", "rules-engine")

	exitCode := m.Run()

	log.Println("[TEARDOWN]")

	utils.SnapDumpLogs(nil, ascSnap)

	utils.SnapRemove(nil,
		ascSnap,
		"edgexfoundry",
	)

	os.Exit(exitCode)
}
