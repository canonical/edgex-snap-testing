package test

import (
	"edgex-snap-testing/test/utils"
	"log"
	"os"
	"testing"
)

const (
	deviceOnvifCameraSnap = "edgex-device-onvif-camera"
	deviceOnvifCameraApp  = "device-onvif-camera"
)

func TestMain(m *testing.M) {
	teardown, err := utils.SetupServiceTests(deviceOnvifCameraSnap)
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
		Snap:                 deviceOnvifCameraSnap,
		App:                  deviceOnvifCameraApp,
	})

	utils.TestConfig(t, deviceOnvifCameraSnap, utils.Config{
		TestChangePort: utils.ConfigChangePort{
			App:                      deviceOnvifCameraApp,
			DefaultPort:              utils.ServicePort(deviceOnvifCameraApp),
			TestAppConfig:            true,
			TestGlobalConfig:         true,
			TestMixedGlobalAppConfig: utils.FullConfigTest,
		},
		TestAutoStart: true,
	})

	utils.TestNet(t, deviceOnvifCameraSnap, utils.Net{
		StartSnap:        true,
		TestOpenPorts:    []string{utils.ServicePort(deviceOnvifCameraApp)},
		TestBindLoopback: []string{utils.ServicePort(deviceOnvifCameraApp)},
	})

	utils.TestPackaging(t, deviceOnvifCameraSnap, utils.Packaging{
		TestSemanticSnapVersion: true,
	})
}
