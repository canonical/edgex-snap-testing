package test

import (
	"edgex-snap-testing/test/utils"
	"fmt"
	"log"
	"os"
	"testing"
	"time"
)

const (
	deviceOnvifcameraSnap        = "edgex-device-onvif-camera"
	deviceOnvifCameraApp         = "device-onvif-camera"
	deviceOnvifcameraService     = deviceOnvifcameraSnap + "." + deviceOnvifCameraApp
	deviceOnvifcameraServicePort = "59984"
	testOnvifCredFile            = "/var/snap/edgex-device-onvif-camera/current/config/device-onvif-camera/test-onvif-credentials.json"
)

func TestMain(m *testing.M) {

	log.Println("[SETUP]")
	start := time.Now()
	// start clean
	utils.SnapRemove(nil,
		deviceOnvifcameraSnap,
		"edgexfoundry",
	)

	// install the device snap before edgexfoundry
	// to catch build error sooner and stop
	if utils.LocalSnap != "" {
		utils.SnapInstallFromFile(nil, utils.LocalSnap)
	} else {
		utils.SnapInstallFromStore(nil, deviceOnvifcameraSnap, utils.ServiceChannel)
	}
	utils.SnapInstallFromStore(nil, "edgexfoundry", utils.PlatformChannel)

	// make sure all services are online before starting the tests
	utils.WaitPlatformOnline(nil)

	// for local build, the interface isn't auto-connected.
	// connect manually regardless
	utils.SnapConnect(nil,
		"edgexfoundry:edgex-secretstore-token",
		deviceOnvifcameraSnap+":edgex-secretstore-token",
	)

	log.Printf("Copying test credentials to %s", testOnvifCredFile)
	// seed test onvif credentials
	testData, err := os.ReadFile("test-onvif-credentials.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = os.WriteFile(testOnvifCredFile, testData, 0644)
	err = os.WriteFile(testOnvifCredFile, testData, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	utils.SnapSet(nil, deviceOnvifcameraSnap, "config.secretstore-secretsfile", testOnvifCredFile)

	exitCode := m.Run()

	log.Println("[TEARDOWN]")

	utils.SnapDumpLogs(nil, start, deviceOnvifcameraSnap)

	utils.SnapRemove(nil,
		deviceOnvifcameraSnap,
		"edgexfoundry",
	)

	os.Exit(exitCode)
}

func TestCommon(t *testing.T) {
	utils.TestConfig(t, deviceOnvifcameraSnap, utils.Config{
		TestChangePort: utils.ConfigChangePort{
			App:                      deviceOnvifCameraApp,
			DefaultPort:              deviceOnvifcameraServicePort,
			TestLegacyEnvConfig:      utils.FullConfigTest,
			TestAppConfig:            true,
			TestGlobalConfig:         true,
			TestMixedGlobalAppConfig: utils.FullConfigTest,
		},
	})

	utils.TestNet(t, deviceOnvifcameraSnap, utils.Net{
		StartSnap:        true,
		TestOpenPorts:    []string{deviceOnvifcameraServicePort},
		TestBindLoopback: []string{deviceOnvifcameraServicePort},
	})

	utils.TestPackaging(t, deviceOnvifcameraSnap, utils.Packaging{
		TestSemanticSnapVersion: true,
	})
}
