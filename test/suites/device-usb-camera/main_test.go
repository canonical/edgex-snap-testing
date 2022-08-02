package test

import (
	"edgex-snap-testing/test/utils"
	"log"
	"os"
	"testing"
	"time"
)

const (
	deviceUSBCamSnap        = "edgex-device-usb-camera"
	deviceUSBCamApp         = "device-usb-camera"
	deviceUSBCamServicePort = "59983"
	rtspServerPort          = "8554"
)

func TestMain(m *testing.M) {
	// start clean
	utils.SnapRemove(nil,
		deviceUSBCamSnap,
		"edgexfoundry",
	)

	log.Println("[SETUP]")
	start := time.Now()

	// install the device snap before edgexfoundry
	// to catch build error sooner and stop
	if utils.LocalSnap() {
		utils.SnapInstallFromFile(nil, utils.LocalSnapPath)
	} else {
		utils.SnapInstallFromStore(nil, deviceUSBCamSnap, utils.ServiceChannel)
	}
	utils.SnapInstallFromStore(nil, "edgexfoundry", utils.PlatformChannel)

	// make sure all services are online before starting the tests
	utils.WaitPlatformOnline(nil)

	// for local build, the interface isn't auto-connected.
	// connect manually
	if utils.LocalSnap() {
		utils.SnapConnect(nil,
			"edgexfoundry:edgex-secretstore-token",
			deviceUSBCamSnap+":edgex-secretstore-token",
		)
	}

	exitCode := m.Run()

	log.Println("[TEARDOWN]")

	utils.SnapDumpLogs(nil, start, deviceUSBCamSnap)
	utils.SnapDumpLogs(nil, start, "edgexfoundry")

	utils.SnapRemove(nil,
		deviceUSBCamSnap,
		"edgexfoundry",
	)

	os.Exit(exitCode)
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
