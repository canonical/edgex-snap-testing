package test

import (
	"edgex-snap-testing/test/utils"
	"log"
	"os"
	"testing"
	"time"
)

const (
	platformSnap            = "edgexfoundry"
	supportSchedulerApp     = "support-scheduler"
	supportSchedulerService = platformSnap + "." + supportSchedulerApp

	supportSchedulerServicePort = "59861"
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

	utils.SnapStart(nil, supportSchedulerService)

	exitCode := m.Run()

	log.Println("[TEARDOWN]")

	utils.SnapDumpLogs(nil, start, platformSnap)

	utils.SnapRemove(nil,
		platformSnap,
	)

	os.Exit(exitCode)
}

func TestCommon(t *testing.T) {
	// check network interface status for all platform ports except for:
	// Kong’s port: 8000
	// Kong-db's port: 5432
	// Redis's port: 6379
	var localPlatformPorts []string
	for _, port := range utils.PlatformPorts {
		if port != "8000" && port != "5432" && port != "6379" {
			localPlatformPorts = append(localPlatformPorts, port)
		}
	}

	utils.TestConfig(t, platformSnap, utils.Config{
		TestChangePort: utils.ConfigChangePort{
			App:                      supportSchedulerApp,
			DefaultPort:              supportSchedulerServicePort,
			TestLegacyEnvConfig:      false, // schemes differ, run specific test instead
			TestAppConfig:            true,
			TestGlobalConfig:         true,
			TestMixedGlobalAppConfig: utils.FullConfigTest,
		},
	})

	utils.TestNet(t, platformSnap, utils.Net{
		StartSnap:        false, // the service are started by default
		TestOpenPorts:    utils.PlatformPorts,
		TestBindLoopback: localPlatformPorts,
	})

	utils.TestPackaging(t, platformSnap, utils.Packaging{
		TestSemanticSnapVersion: true,
	})

	utils.TestRefresh(t, platformSnap)
}
