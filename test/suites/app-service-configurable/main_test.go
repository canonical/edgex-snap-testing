package test

import (
	"edgex-snap-testing/test/utils"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	log.Println("[SETUP]")

	// start clean
	utils.SnapRemove(nil,
		"edgex-app-service-configurable",
		"edgexfoundry",
	)

	// install the app-service-configurable snap before edgexfoundry
	// to catch build error sooner and stop
	if utils.LocalSnap != "" {
		utils.SnapInstallFromFile(nil, utils.LocalSnap)
	} else {
		utils.SnapInstallFromStore(nil, "edgex-app-service-configurable", utils.ServiceChannel)
	}
	utils.SnapInstallFromStore(nil, "edgexfoundry", utils.PlatformChannel)

	// for local build, the interface isn't auto-connected.
	// connect manually regardless
	utils.SnapConnect(nil,
		"edgexfoundry:edgex-secretstore-token",
		"edgex-app-service-configurable:edgex-secretstore-token",
	)

	// set profile to rules engine
	utils.Exec(nil, "sudo snap set edgex-app-service-configurable profile=rules-engine")

	exitCode := m.Run()

	log.Println("[TEARDOWN]")

	utils.SnapDumpLogs(nil, "edgex-app-service-configurable")

	utils.SnapRemove(nil,
		"edgex-app-service-configurable",
		"edgexfoundry",
	)

	os.Exit(exitCode)
}
