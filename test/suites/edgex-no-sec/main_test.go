package test

import (
	"edgex-snap-testing/test/utils"
	"log"
	"os"
	"testing"
	"time"
)

const (
	platformSnap = "edgexfoundry"

	deviceVirtualSnap = "edgex-device-virtual"
	deviceVirtualPort = "59900"

	ekuiperSnap           = "edgex-ekuiper"
	ekuiperApp            = "kuiper"
	ekuiperService        = ekuiperSnap + "." + ekuiperApp
	ekuiperServerPort     = "20498"
	ekuiperRestfulApiPort = "59720"

	ascSnap             = "edgex-app-service-configurable"
	ascServiceRulesPort = "59701"
)

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

	utils.TestDeviceVirtualReading(t)
}

func setup() (teardown func(), err error) {
	log.Println("[CLEAN]")
	utils.SnapRemove(nil, platformSnap, deviceVirtualSnap, ekuiperSnap, ascSnap)

	log.Println("[SETUP]")
	start := time.Now()

	teardown = func() {
		log.Println("[TEARDOWN]")

		utils.SnapDumpLogs(nil, start, platformSnap)
		utils.SnapDumpLogs(nil, start, deviceVirtualSnap)
		utils.SnapDumpLogs(nil, start, ekuiperSnap)
		utils.SnapDumpLogs(nil, start, ascSnap)

		utils.SnapRemove(nil, platformSnap)
		utils.SnapRemove(nil, deviceVirtualSnap)
		utils.SnapRemove(nil, ekuiperSnap)
		utils.SnapRemove(nil, ascSnap)
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

	if err = utils.SnapInstallFromStore(nil, ekuiperSnap, utils.ServiceChannel); err != nil {
		teardown()
		return
	}

	if err = utils.SnapInstallFromStore(nil, ascSnap, utils.ServiceChannel); err != nil {
		teardown()
		return
	}

	// turn security off
	utils.SnapSet(nil, platformSnap, "security-secret-store", "off")
	utils.SnapSet(nil, deviceVirtualSnap, "config.edgex-security-secret-store", "false")
	utils.SnapSet(nil, ascSnap, "app-options", "true")
	utils.SnapSet(nil, ascSnap, "config.edgex-security-secret-store", "false")
	utils.SnapSet(nil, ekuiperSnap, "config.edgex-security-secret-store", "false")

	// make sure all services are online before starting the tests
	utils.SnapStart(nil, deviceVirtualSnap)
	if err = utils.WaitServiceOnline(nil, 60, deviceVirtualPort); err != nil {
		teardown()
		return
	}

	// subscribe to ASC events
	utils.SnapSet(nil, ekuiperSnap, "config.edgex.default.topic", "rules-events")
	utils.SnapSet(nil, ekuiperSnap, "config.edgex.default.messagetype", "event")
	utils.SnapStart(nil, ekuiperSnap)
	if err = utils.WaitServiceOnline(nil, 60, ekuiperServerPort, ekuiperRestfulApiPort); err != nil {
		teardown()
		return
	}

	// set profile to rules engine
	utils.SnapSet(nil, ascSnap, "profile", "rules-engine")
	utils.SnapStart(nil, ascSnap)
	if err = utils.WaitServiceOnline(nil, 60, ascServiceRulesPort); err != nil {
		teardown()
		return
	}

	if err = utils.WaitServiceOnline(nil, 180, platformPortsNoSec()...); err != nil {
		teardown()
		return
	}

	return
}
