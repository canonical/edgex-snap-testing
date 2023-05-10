package test

import (
	"edgex-snap-testing/test/utils"
	"log"
	"os"
	"testing"
	"time"
)

const platformSnap = "edgexfoundry"

func platformPortsNoSec() []string {
	return []string{
		utils.ServicePort("core-data"),
		utils.ServicePort("core-metadata"),
		utils.ServicePort("core-command"),
		utils.ServicePort("consul"),
		utils.ServicePort("redis"),
	}
}

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
	utils.TestNet(t, platformSnap, utils.Net{
		StartSnap:        false, // the service are started by default
		TestOpenPorts:    platformPortsNoSec(),
		TestBindLoopback: platformPortsNoSec(),
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

	// disable security - this sets autostart of security services to false
	utils.SnapSet(nil, platformSnap, "security", "false")

	// enable autostart globally to start all services apart from security services that are explicitly disabled:
	utils.SnapSet(nil, platformSnap, "autostart", "true")
	// The above is equivalent to starting non-security services manually:
	// utils.SnapStart(nil, func() (names []string) {
	// 	nonSecServices := []string{
	// 		"consul", "redis",
	// 		"core-common-config-bootstrapper",
	// 		"core-data", "core-metadata", "core-command",
	// 		"support-scheduler", "support-notifications",
	// 	}
	// 	for _, s := range nonSecServices {
	// 		names = append(names, "edgexfoundry."+s)
	// 	}
	// 	return
	// }()...)

	// make sure all services are online before starting the tests
	if err = utils.WaitServiceOnline(nil, 180, platformPortsNoSec()...); err != nil {
		teardown()
		return
	}

	return
}
