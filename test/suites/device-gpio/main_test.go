package test

import (
	"edgex-snap-testing/test/utils"
	"log"
	"os"
	"testing"
	"time"
)

const (
	deviceGpioSnap        = "edgex-device-gpio"
	deviceGpioApp         = "device-gpio"
	deviceGpioServicePort = "59910"
)

func TestMain(m *testing.M) {
	// start clean
	utils.SnapRemove(nil,
		deviceGpioSnap,
		"edgexfoundry",
	)

	log.Println("[SETUP]")
	start := time.Now()

	// install the device snap before edgexfoundry
	// to catch build error sooner and stop
	if utils.LocalSnap != "" {
		utils.SnapInstallFromFile(nil, utils.LocalSnap)

		// for local build, the interface isn't auto-connected.
		// connect manually
		utils.SnapConnect(nil,
			"edgexfoundry:edgex-secretstore-token",
			deviceGpioSnap+":edgex-secretstore-token",
		)
	} else {
		utils.SnapInstallFromStore(nil, deviceGpioSnap, utils.ServiceChannel)
	}
	utils.SnapInstallFromStore(nil, "edgexfoundry", utils.PlatformChannel)

	// make sure all services are online before starting the tests
	utils.WaitPlatformOnline(nil)

	exitCode := m.Run()

	log.Println("[TEARDOWN]")

	utils.SnapDumpLogs(nil, start, deviceGpioSnap)
	utils.SnapDumpLogs(nil, start, "edgexfoundry")

	utils.SnapRemove(nil,
		deviceGpioSnap,
		"edgexfoundry",
	)

	os.Exit(exitCode)
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
