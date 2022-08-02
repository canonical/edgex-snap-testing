package test

import (
	"edgex-snap-testing/test/utils"
	"log"
	"os"
	"testing"
)

const (
	deviceOnvifCameraSnap        = "edgex-device-onvif-camera"
	deviceOnvifCameraApp         = "device-onvif-camera"
	deviceOnvifCameraServicePort = "59984"
)

func TestMain(m *testing.M) {
	code, err := utils.RunDeviceTests(m, deviceOnvifCameraSnap)
	if err != nil {
		log.Fatalf("Failed to run tests: %s", err)
	}
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
			DefaultPort:              deviceOnvifCameraServicePort,
			TestAppConfig:            true,
			TestGlobalConfig:         true,
			TestMixedGlobalAppConfig: utils.FullConfigTest,
		},
	})

	utils.TestNet(t, deviceOnvifCameraSnap, utils.Net{
		StartSnap:        true,
		TestOpenPorts:    []string{deviceOnvifCameraServicePort},
		TestBindLoopback: []string{deviceOnvifCameraServicePort},
	})

	utils.TestPackaging(t, deviceOnvifCameraSnap, utils.Packaging{
		TestSemanticSnapVersion: true,
	})
}
