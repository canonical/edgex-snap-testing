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
	utils.WaitServiceOnline(t, 60, supportSchedulerServicePort)

	// check network interface status for all platform ports except for:
	// Kong’s port: 8000
	// Kong-db's port: 5432
	// Redis's port: 6379
	var platformPorts []string
	for _, port := range utils.PlatformPorts {
		if port != "8000" && port != "5432" && port != "6379" {
			platformPorts = append(platformPorts, port)
		}
	}

	utils.TestCommon(t, utils.TestParams{
		Snap: platformSnap,
		App:  supportSchedulerApp,
		TestConfigs: utils.TestConfigs{
			TestEnvConfig:      utils.FullConfigTest,
			TestAppConfig:      true,
			TestGlobalConfig:   true,
			TestMixedConfig:    utils.FullConfigTest,
			DefaultServicePort: supportSchedulerServicePort,
		},
		TestNetworking: utils.TestNetworking{
			TestOpenPorts:        platformPorts,
			TestBindAddrLoopback: true,
		},
		TestVersion: utils.TestVersion{
			TestSemanticSnapVersion: true,
		},
	})
}
