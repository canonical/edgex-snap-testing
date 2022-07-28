package test

import (
	"edgex-snap-testing/test/utils"
	"log"
	"os"
	"testing"
	"time"
)

const (
	deviceVirtualSnap        = "edgex-device-virtual"
	deviceVirtualApp         = "device-virtual"
	deviceVirtualServicePort = "59900"
)

func TestMain(m *testing.M) {
	start := time.Now()

	log.Println("[SETUP]")

	// start clean
	utils.SnapRemove(nil,
		deviceVirtualSnap,
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
			deviceVirtualSnap+":edgex-secretstore-token",
		)
	} else {
		utils.SnapInstallFromStore(nil, deviceVirtualSnap, utils.ServiceChannel)
	}
	utils.SnapInstallFromStore(nil, "edgexfoundry", utils.PlatformChannel)

	// make sure all services are online before starting the tests
	utils.WaitPlatformOnline(nil)

	exitCode := m.Run()

	log.Println("[TEARDOWN]")

	utils.SnapDumpLogs(nil, start, deviceVirtualSnap)
	utils.SnapDumpLogs(nil, start, "edgexfoundry")

	utils.SnapRemove(nil,
		deviceVirtualSnap,
		"edgexfoundry",
	)

	os.Exit(exitCode)
}

func TestCommon(t *testing.T) {
	utils.TestContentInterfaces(t, utils.ContentInterfaces{
		TestSecretstoreToken: true,
		Snap:                 deviceVirtualSnap,
		App:                  deviceVirtualApp,
	})

	utils.TestConfig(t, deviceVirtualSnap, utils.Config{
		TestChangePort: utils.ConfigChangePort{
			App:                      deviceVirtualApp,
			DefaultPort:              deviceVirtualServicePort,
			TestAppConfig:            true,
			TestGlobalConfig:         true,
			TestMixedGlobalAppConfig: utils.FullConfigTest,
		},
	})

	utils.TestNet(t, deviceVirtualSnap, utils.Net{
		StartSnap:        true,
		TestOpenPorts:    []string{deviceVirtualServicePort},
		TestBindLoopback: []string{deviceVirtualServicePort},
	})

	utils.TestPackaging(t, deviceVirtualSnap, utils.Packaging{
		TestSemanticSnapVersion: true,
	})
}
