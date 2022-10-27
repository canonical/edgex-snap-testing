package test

import (
	"edgex-snap-testing/test/utils"
	"log"
	"os"
	"testing"
)

const (
	deviceGpioSnap        = "edgex-device-gpio"
	deviceGpioApp         = "device-gpio"
	deviceGpioServicePort = "59910"
)

func TestMain(m *testing.M) {
	teardown, err := utils.SetupServiceTests(deviceGpioSnap)
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
		Snap:                 deviceGpioSnap,
		App:                  deviceGpioApp,
	})

	utils.TestConfig(t, deviceGpioSnap, utils.Config{
		TestChangePort: utils.ConfigChangePort{
			App:                      deviceGpioApp,
			DefaultPort:              deviceGpioServicePort,
			TestLegacyEnvConfig:      utils.FullConfigTest,
			TestAppConfig:            true,
			TestGlobalConfig:         true,
			TestMixedGlobalAppConfig: utils.FullConfigTest,
		},
		TestAutoStart: true,
	})

	utils.TestNet(t, deviceGpioSnap, utils.Net{
		StartSnap:        true,
		TestOpenPorts:    []string{deviceGpioServicePort},
		TestBindLoopback: []string{deviceGpioServicePort},
	})

	utils.TestPackaging(t, deviceGpioSnap, utils.Packaging{
		TestSemanticSnapVersion: true,
	})
}
