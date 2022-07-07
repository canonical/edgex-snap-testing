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
	deviceVirtualService     = deviceVirtualSnap + "." + deviceVirtualApp
	deviceVirtualServicePort = "59900"
)

var start = time.Now()

func TestMain(m *testing.M) {

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
	} else {
		utils.SnapInstallFromStore(nil, deviceVirtualSnap, utils.ServiceChannel)
	}
	utils.SnapInstallFromStore(nil, "edgexfoundry", utils.PlatformChannel)

	// make sure all services are online before starting the tests
	utils.WaitPlatformOnline(nil)

	// for local build, the interface isn't auto-connected.
	// connect manually regardless
	utils.SnapConnect(nil,
		"edgexfoundry:edgex-secretstore-token",
		deviceVirtualSnap+":edgex-secretstore-token",
	)

	// Start the service so that the default config gets uploaded to consul.
	// Otherwise, settings that get passed using environment variables on first start get uploaded
	// and become the default.
	utils.SnapStart(nil, deviceVirtualService)
	utils.WaitServiceOnline(nil, 60, deviceVirtualServicePort)

	exitCode := m.Run()

	log.Println("[TEARDOWN]")

	utils.SnapDumpLogs(nil, start, deviceVirtualSnap)

	utils.SnapRemove(nil,
		deviceVirtualSnap,
		"edgexfoundry",
	)

	os.Exit(exitCode)
}

func TestCommon(t *testing.T) {
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
