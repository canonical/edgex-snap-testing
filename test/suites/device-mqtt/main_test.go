package test

import (
	"edgex-snap-testing/test/utils"
	"log"
	"os"
	"testing"
)

const (
	deviceMqttSnap        = "edgex-device-mqtt"
	deviceMqttApp         = "device-mqtt"
	deviceMqttServicePort = "59982"
)

func TestMain(m *testing.M) {
	teardown, err := utils.SetupDeviceTests(deviceMqttSnap)
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
		Snap:                 deviceMqttSnap,
		App:                  deviceMqttApp,
	})

	utils.TestConfig(t, deviceMqttSnap, utils.Config{
		TestChangePort: utils.ConfigChangePort{
			App:                      deviceMqttApp,
			DefaultPort:              deviceMqttServicePort,
			TestLegacyEnvConfig:      utils.FullConfigTest,
			TestAppConfig:            true,
			TestGlobalConfig:         true,
			TestMixedGlobalAppConfig: utils.FullConfigTest,
		},
	})

	utils.TestNet(t, deviceMqttSnap, utils.Net{
		StartSnap:        true,
		TestOpenPorts:    []string{deviceMqttServicePort},
		TestBindLoopback: []string{deviceMqttServicePort},
	})

	utils.TestPackaging(t, deviceMqttSnap, utils.Packaging{
		TestSemanticSnapVersion: true,
	})
}
