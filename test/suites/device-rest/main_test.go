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
	deviceRestService     = deviceRestSnap + "." + deviceRestApp
	deviceRestServicePort = "59986"
)

func TestMain(m *testing.M) {
	start := time.Now()

	log.Println("[SETUP]")

	// start clean
	utils.SnapRemove(nil,
		deviceRestSnap,
		"edgexfoundry",
	)

	// install the device snap before edgexfoundry
	// to catch build error sooner and stop
	if utils.LocalSnap != "" {
		utils.SnapInstallFromFile(nil, utils.LocalSnap)
	} else {
		utils.SnapInstallFromStore(nil, deviceRestSnap, utils.ServiceChannel)
	}
	utils.SnapInstallFromStore(nil, "edgexfoundry", utils.PlatformChannel)

	// make sure all services are online before starting the tests
	utils.WaitPlatformOnline(nil)

	// for local build, the interface isn't auto-connected.
	// connect manually regardless
	utils.SnapConnect(nil,
		"edgexfoundry:edgex-secretstore-token",
		deviceRestSnap+":edgex-secretstore-token",
	)

	// Start the service so that the default config gets uploaded to consul.
	// Otherwise, settings that get passed using environment variables on first start get uploaded
	// and become the default.
	utils.SnapStart(nil, deviceRestService)
	utils.WaitServiceOnline(nil, 60, deviceRestServicePort)

	exitCode := m.Run()

	log.Println("[TEARDOWN]")

	utils.SnapDumpLogs(nil, start, deviceRestSnap)

	utils.SnapRemove(nil,
		deviceRestSnap,
		"edgexfoundry",
	)

	os.Exit(exitCode)
}

func TestCommon(t *testing.T) {
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
