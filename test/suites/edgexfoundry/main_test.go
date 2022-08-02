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

func main(m *testing.M) (int, error) {
	log.Println("[CLEAN]")
	utils.SnapRemove(nil,
		platformSnap,
	)

	log.Println("[SETUP]")

	// add this to the bottom of the stack to remove after collecting logs
	defer utils.SnapRemove(nil, platformSnap)

	start := time.Now()
	defer utils.SnapDumpLogs(nil, start, platformSnap)

	var err error

	if utils.LocalSnap != "" {
		err = utils.SnapInstallFromFile(nil, utils.LocalSnap)
	} else {
		err = utils.SnapInstallFromStore(nil, platformSnap, utils.ServiceChannel)
	}
	if err != nil {
		return 0, err
	}

	// make sure all services are online before starting the tests
	err = utils.WaitPlatformOnline(nil)
	if err != nil {
		return 0, err
	}

	// support-scheduler is disabled by default.
	// Start it to have the default configurations registered in the EdgeX Registry
	//	in preparation for the local config tests.
	utils.SnapStart(nil, supportSchedulerService)
	utils.WaitServiceOnline(nil, 60, supportSchedulerServicePort)

	log.Println("[START]")
	return m.Run(), nil
}

func TestMain(m *testing.M) {
	code, err := main(m)
	if err != nil {
		log.Fatalf("Failed to run tests: %s", err)
	}
	os.Exit(code)
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
			TestGlobalConfig:         false, // multiple servers, test setting startup message instead
			TestMixedGlobalAppConfig: false, // multiple servers, test setting startup message instead
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
