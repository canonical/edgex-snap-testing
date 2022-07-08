package test

import (
	"edgex-snap-testing/test/utils"
	"log"
	"os"
	"testing"
	"time"
)

const (
	deviceModbusSnap        = "edgex-device-modbus"
	deviceModbusApp         = "device-modbus"
	deviceModbusService     = deviceModbusSnap + "." + deviceModbusApp
	deviceModbusServicePort = "59901"
)

func TestMain(m *testing.M) {
	start := time.Now()

	log.Println("[SETUP]")

	// start clean
	utils.SnapRemove(nil,
		deviceModbusSnap,
		"edgexfoundry",
	)

	// install the device snap before edgexfoundry
	// to catch build error sooner and stop
	if utils.LocalSnap != "" {
		utils.SnapInstallFromFile(nil, utils.LocalSnap)
	} else {
		utils.SnapInstallFromStore(nil, deviceModbusSnap, utils.ServiceChannel)
	}
	utils.SnapInstallFromStore(nil, "edgexfoundry", utils.PlatformChannel)

	// make sure all services are online before starting the tests
	utils.WaitPlatformOnline(nil)

	// for local build, the interface isn't auto-connected.
	// connect manually regardless
	utils.SnapConnect(nil,
		"edgexfoundry:edgex-secretstore-token",
		deviceModbusSnap+":edgex-secretstore-token",
	)

	exitCode := m.Run()

	log.Println("[TEARDOWN]")

	utils.SnapDumpLogs(nil, start, deviceModbusSnap)

	utils.SnapRemove(nil,
		deviceModbusSnap,
		"edgexfoundry",
	)

	os.Exit(exitCode)
}

func TestCommon(t *testing.T) {
	utils.TestConfig(t, deviceModbusSnap, utils.Config{
		TestChangePort: utils.ConfigChangePort{
			App:                      deviceModbusApp,
			DefaultPort:              deviceModbusServicePort,
			TestLegacyEnvConfig:      utils.FullConfigTest,
			TestAppConfig:            true,
			TestGlobalConfig:         true,
			TestMixedGlobalAppConfig: utils.FullConfigTest,
		},
	})

	utils.TestNet(t, deviceModbusSnap, utils.Net{
		StartSnap:        true,
		TestOpenPorts:    []string{deviceModbusServicePort},
		TestBindLoopback: []string{deviceModbusServicePort},
	})

	utils.TestPackaging(t, deviceModbusSnap, utils.Packaging{
		TestSemanticSnapVersion: true,
	})
}
