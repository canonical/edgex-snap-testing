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
)

func TestMain(m *testing.M) {
	teardown, err := setup()
	if err != nil {
		log.Fatalf("Failed to setup tests: %s", err)
	}

	code := m.Run()
	teardown()

	os.Exit(code)
}

func TestCommon(t *testing.T) {

	utils.TestConfig(t, platformSnap, utils.Config{
		TestChangePort: utils.ConfigChangePort{
			App:                      supportSchedulerApp,
			DefaultPort:              utils.ServicePort(supportSchedulerApp),
			TestAppConfig:            false, // covered in local startup message testing
			TestGlobalConfig:         false, // multiple servers, test setting startup message instead
			TestMixedGlobalAppConfig: false, // multiple servers, test setting startup message instead
		},
	})

	utils.TestNet(t, platformSnap, utils.Net{
		StartSnap:        false, // the service are started by default
		TestOpenPorts:    utils.PlatformPorts(true),
		TestBindLoopback: utils.PlatformPorts(false), // exclude public ports
	})

	utils.TestPackaging(t, platformSnap, utils.Packaging{
		TestSemanticSnapVersion: true,
	})
}

func setup() (teardown func(), err error) {
	log.Println("[CLEAN]")
	utils.SnapRemove(nil, platformSnap)

	log.Println("[SETUP]")
	start := time.Now()

	teardown = func() {
		log.Println("[TEARDOWN]")
		utils.SnapDumpLogs(nil, start, platformSnap)

		log.Println("Removing installed snap:", !utils.SkipTeardownRemoval)
		if !utils.SkipTeardownRemoval {
			utils.SnapRemove(nil, platformSnap)
		}
	}

	if utils.LocalSnap() {
		err = utils.SnapInstallFromFile(nil, utils.LocalSnapPath)
	} else {
		err = utils.SnapInstallFromStore(nil, platformSnap, utils.PlatformChannel)
	}
	if err != nil {
		teardown()
		return
	}

	// make sure all services are online before starting the tests
	if err = utils.WaitPlatformOnline(nil); err != nil {
		teardown()
		return
	}

	// support-scheduler is disabled by default.
	// Start it to have the default configurations registered in the EdgeX Registry
	//	in preparation for the local config tests.
	utils.SnapStart(nil, supportSchedulerService)
	if err = utils.WaitServiceOnline(nil, 60, utils.ServicePort(supportSchedulerApp)); err != nil {
		teardown()
		return
	}

	return
}
