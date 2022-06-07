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

func TestCommon(t *testing.T) {
	params := &utils.TestParams{
		Snap:               "edgexfoundry",
		App:                "core-metadata",
		DefaultServicePort: "59881",
		TestConfigs: utils.TestConfigs{
			TestEnvConfig:    true,
			TestAppConfig:    true,
			TestGlobalConfig: true,
			TestMixedConfig:  true,
		},
		TestNetworking: utils.TestNetworking{
			TestOpenPorts:        []string{"59881"},
			TestBindAddrLoopback: true,
		},
		TestVersion: utils.TestVersion{
			TestSemanticSnapVersion: true,
		},
	}
	utils.TestCommon(t, params)
}

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


	os.Exit(exitCode)
}
