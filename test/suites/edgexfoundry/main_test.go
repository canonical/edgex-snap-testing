package test

import (
	"edgex-snap-testing/test/utils"
	"log"
	"os"
	"testing"
	"time"
)

const (
	platformSnap        = "edgexfoundry"
	coreMetadataApp     = "core-metadata"
	coreMetadataService = platformSnap + "." + coreMetadataApp

	coreMetadataDefaultServicePort = "59881"
)

var start = time.Now()

func TestMain(m *testing.M) {

	log.Println("[SETUP]")

	// start clean
	utils.SnapRemove(nil,
		platformSnap,
	)

	if utils.LocalSnap != "" {
		utils.SnapInstallFromFile(nil, utils.LocalSnap)
	} else {
		utils.SnapInstallFromStore(nil, platformSnap, utils.ServiceChannel)
	}

	// make sure all services are online before starting the tests
	utils.WaitPlatformOnline(nil)

	// make sure core-metadata service starts and comes online before starting the tests
	utils.SnapStart(nil, coreMetadataService)
	utils.WaitServiceOnline(nil, 60, coreMetadataDefaultServicePort)

	exitCode := m.Run()

	log.Println("[TEARDOWN]")

	utils.SnapDumpLogs(nil, start, platformSnap)

	utils.SnapRemove(nil,
		platformSnap,
	)

	FullConfigTest = false

	os.Exit(exitCode)
}
