package test

import (
	"edgex-snap-testing/test/utils"
	"log"
	"os"
	"testing"
	"time"
)

const (
	platformSnap      = "edgexfoundry"
	deviceVirtualSnap = "edgex-device-virtual"
	deviceVirtualPort = "59900"
)

func PlatformPortsNoSecurity(includePublicPorts bool) (ports []string) {
	ports = append(ports,
		utils.ServicePort("core-data"),
		utils.ServicePort("core-metadata"),
		utils.ServicePort("core-command"),
		utils.ServicePort("consul"),
		utils.ServicePort("redis"),
	)
	return
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
		TestOpenPorts:    PlatformPortsNoSecurity(true),
		TestBindLoopback: PlatformPortsNoSecurity(false), // exclude public ports
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

	if err = utils.WaitServiceOnline(nil, 180, PlatformPortsNoSecurity(false)...); err != nil {
		teardown()
		return
	}

	return
}
