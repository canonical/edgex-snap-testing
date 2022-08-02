package test

import (
	"edgex-snap-testing/test/utils"
	"log"
	"os"
	"testing"
)

const (
	deviceUSBCamSnap        = "edgex-device-usb-camera"
	deviceUSBCamApp         = "device-usb-camera"
	deviceUSBCamServicePort = "59983"
	rtspServerPort          = "8554"
)

func TestMain(m *testing.M) {
	teardown, err := utils.SetupServiceTests(deviceUSBCamSnap)
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
		Snap:                 deviceUSBCamSnap,
		App:                  deviceUSBCamApp,
	})

	utils.TestConfig(t, deviceUSBCamSnap, utils.Config{
		TestChangePort: utils.ConfigChangePort{
			App:                      deviceUSBCamApp,
			DefaultPort:              deviceUSBCamServicePort,
			TestAppConfig:            true,
			TestGlobalConfig:         true,
			TestMixedGlobalAppConfig: utils.FullConfigTest,
		},
	})

	utils.TestNet(t, deviceUSBCamSnap, utils.Net{
		StartSnap:        true,
		TestOpenPorts:    []string{deviceUSBCamServicePort, rtspServerPort},
		TestBindLoopback: []string{deviceUSBCamServicePort, rtspServerPort},
	})

	utils.TestPackaging(t, deviceUSBCamSnap, utils.Packaging{
		TestSemanticSnapVersion: true,
	})
}
