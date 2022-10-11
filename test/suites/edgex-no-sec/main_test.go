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

	deviceVirtualSnap = "edgex-device-virtual"
	deviceVirtualPort = "59900"
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
			DefaultPort:              supportSchedulerServicePort,
			TestLegacyEnvConfig:      false, // schemes differ, run specific test instead
			TestAppConfig:            true,
			TestGlobalConfig:         false, // multiple servers, test setting startup message instead
			TestMixedGlobalAppConfig: false, // multiple servers, test setting startup message instead
		},
	})

	utils.TestNet(t, platformSnap, utils.Net{
		StartSnap:        false, // the service are started by default
		TestOpenPorts:    utils.PlatformPortsNoSecurity(true),
		TestBindLoopback: utils.PlatformPortsNoSecurity(false), // exclude public ports
	})

	utils.TestDeviceVirtualReading(t)
}

func setup() (teardown func(), err error) {
	log.Println("[CLEAN]")
	utils.SnapRemove(nil, platformSnap, deviceVirtualSnap)

	log.Println("[SETUP]")
	start := time.Now()

	teardown = func() {
		log.Println("[TEARDOWN]")

		utils.SnapDumpLogs(nil, start, platformSnap)
		utils.SnapDumpLogs(nil, start, deviceVirtualSnap)

		utils.SnapRemove(nil, platformSnap)
		utils.SnapRemove(nil, deviceVirtualSnap)
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

	if err = utils.SnapInstallFromStore(nil, deviceVirtualSnap, utils.ServiceChannel); err != nil {
		teardown()
		return
	}

	// turn security off
	utils.SnapSet(nil, platformSnap, "security-secret-store", "off")
	utils.SnapSet(nil, deviceVirtualSnap, "config.edgex-security-secret-store", "false")

	// make sure all services are online before starting the tests
	utils.SnapStart(nil, deviceVirtualSnap)
	if err = utils.WaitServiceOnline(nil, 60, deviceVirtualPort); err != nil {
		teardown()
		return
	}

	if err = utils.WaitPlatformOnlineNoSecurity(nil); err != nil {
		teardown()
		return
	}

	// support-scheduler is disabled by default.
	// Start it to have the default configurations registered in the EdgeX Registry
	//	in preparation for the local config tests.
	utils.SnapStart(nil, supportSchedulerService)
	if err = utils.WaitServiceOnline(nil, 60, supportSchedulerServicePort); err != nil {
		teardown()
		return
	}

	return
}
