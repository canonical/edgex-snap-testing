package test

import (
	"edgex-snap-testing/test/utils"
	"log"
	"os"
	"testing"
	"time"
)

const (
	deviceRestSnap        = "edgex-device-rest"
	deviceRestApp         = "device-rest"
	deviceRestServicePort = "59986"
)

func TestMain(m *testing.M) {
	// start clean
	utils.SnapRemove(nil,
		deviceRestSnap,
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
			deviceRestSnap+":edgex-secretstore-token",
		)
	} else {
		utils.SnapInstallFromStore(nil, deviceRestSnap, utils.ServiceChannel)
	}
	utils.SnapInstallFromStore(nil, "edgexfoundry", utils.PlatformChannel)

	// make sure all services are online before starting the tests
	utils.WaitPlatformOnline(nil)

	exitCode := m.Run()

	log.Println("[TEARDOWN]")

	utils.SnapDumpLogs(nil, start, deviceRestSnap)
	utils.SnapDumpLogs(nil, start, "edgexfoundry")

	utils.SnapRemove(nil,
		deviceRestSnap,
		"edgexfoundry",
	)

	os.Exit(exitCode)
}

func TestCommon(t *testing.T) {
	utils.TestContentInterfaces(t, utils.ContentInterfaces{
		TestSecretstoreToken: true,
		Snap:                 deviceRestSnap,
		App:                  deviceRestApp,
	})

	utils.TestConfig(t, deviceRestSnap, utils.Config{
		TestChangePort: utils.ConfigChangePort{
			App:                      deviceRestApp,
			DefaultPort:              deviceRestServicePort,
			TestLegacyEnvConfig:      utils.FullConfigTest,
			TestAppConfig:            true,
			TestGlobalConfig:         true,
			TestMixedGlobalAppConfig: utils.FullConfigTest,
		},
	})

	utils.TestNet(t, deviceRestSnap, utils.Net{
		StartSnap:        true,
		TestOpenPorts:    []string{deviceRestServicePort},
		TestBindLoopback: []string{deviceRestServicePort},
	})

	utils.TestPackaging(t, deviceRestSnap, utils.Packaging{
		TestSemanticSnapVersion: true,
	})
}
