package test

import (
	"edgex-snap-testing/test/utils"
	"log"
	"os"
	"testing"
)

const (
	deviceMqttSnap = "edgex-device-mqtt"
	deviceMqttApp  = "device-mqtt"
)

func TestMain(m *testing.M) {
	teardown, err := utils.SetupServiceTests(deviceMqttSnap)
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
			DefaultPort:              utils.ServicePort(deviceMqttApp),
			TestAppConfig:            true,
			TestGlobalConfig:         true,
			TestMixedGlobalAppConfig: utils.FullConfigTest,
		},
		TestAutoStart: true,
	})

	utils.TestNet(t, deviceMqttSnap, utils.Net{
		StartSnap:        true,
		TestOpenPorts:    []string{utils.ServicePort(deviceMqttApp)},
		TestBindLoopback: []string{utils.ServicePort(deviceMqttApp)},
	})

	utils.TestPackaging(t, deviceMqttSnap, utils.Packaging{
		TestSemanticSnapVersion: true,
	})
}
