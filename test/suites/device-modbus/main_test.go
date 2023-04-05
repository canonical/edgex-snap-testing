package test

import (
	"edgex-snap-testing/test/utils"
	"log"
	"os"
	"testing"
)

const (
	deviceModbusSnap = "edgex-device-modbus"
	deviceModbusApp  = "device-modbus"
)

func TestMain(m *testing.M) {
	teardown, err := utils.SetupServiceTests(deviceModbusSnap)
	if err != nil {
		log.Fatalf("Failed to setup tests: %s", err)
	}

	code := m.Run()
	teardown()

	os.Exit(code)
}

func TestCommon(t *testing.T) {
	utils.TestContentInterfaces(t, utils.ContentInterfaces{
		TestSecretstoreToken: true,
		Snap:                 deviceModbusSnap,
		App:                  deviceModbusApp,
	})

	utils.TestConfig(t, deviceModbusSnap, utils.Config{
		TestChangePort: utils.ConfigChangePort{
			App:                      deviceModbusApp,
			DefaultPort:              utils.ServicePort(deviceModbusApp),
			TestAppConfig:            true,
			TestGlobalConfig:         true,
			TestMixedGlobalAppConfig: utils.FullConfigTest,
		},
		TestAutoStart: true,
	})

	utils.TestNet(t, deviceModbusSnap, utils.Net{
		StartSnap:        true,
		TestOpenPorts:    []string{utils.ServicePort(deviceModbusApp)},
		TestBindLoopback: []string{utils.ServicePort(deviceModbusApp)},
	})

	utils.TestPackaging(t, deviceModbusSnap, utils.Packaging{
		TestSemanticSnapVersion: true,
	})
}
