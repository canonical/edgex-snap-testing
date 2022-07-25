package test

import (
	"edgex-snap-testing/test/utils"
	"log"
	"os"
	"testing"
	"time"
)

const (
	deviceOnvifCameraSnap        = "edgex-device-onvif-camera"
	deviceOnvifCameraApp         = "device-onvif-camera"
	deviceOnvifCameraServicePort = "59984"
)

func TestMain(m *testing.M) {

	log.Println("[SETUP]")
	start := time.Now()
	// start clean
	utils.SnapRemove(nil,
		deviceOnvifCameraSnap,
		"edgexfoundry",
	)

	// install the device snap before edgexfoundry
	// to catch build error sooner and stop
	if utils.LocalSnap != "" {
		utils.SnapInstallFromFile(nil, utils.LocalSnap)

		// for local build, the interface isn't auto-connected.
		// connect manually
		utils.SnapConnect(nil,
			"edgexfoundry:edgex-secretstore-token",
			deviceOnvifCameraSnap+":edgex-secretstore-token",
		)
	} else {
		utils.SnapInstallFromStore(nil, deviceOnvifCameraSnap, utils.ServiceChannel)
	}
	utils.SnapInstallFromStore(nil, "edgexfoundry", utils.PlatformChannel)

	// make sure all services are online before starting the tests
	utils.WaitPlatformOnline(nil)

	exitCode := m.Run()

	log.Println("[TEARDOWN]")

	utils.SnapDumpLogs(nil, start, deviceOnvifCameraSnap)

	utils.SnapRemove(nil,
		deviceOnvifCameraSnap,
		"edgexfoundry",
	)

	os.Exit(exitCode)
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
			TestLegacyEnvConfig:      utils.FullConfigTest,
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
