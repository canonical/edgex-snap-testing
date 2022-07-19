package test

import (
	"edgex-snap-testing/test/utils"
	"log"
	"os"
	"testing"
	"time"
)

const (
	deviceRfidLlrpSnap        = "edgex-device-rfid-llrp"
	deviceRfidLlrpApp         = "device-rfid-llrp"
	deviceRfidLlrpServicePort = "59989"
)

func TestMain(m *testing.M) {
	start := time.Now()

	log.Println("[SETUP]")

	// start clean
	utils.SnapRemove(nil,
		deviceRfidLlrpSnap,
		"edgexfoundry",
	)

	// install the device snap before edgexfoundry
	// to catch build error sooner and stop
	if utils.LocalSnap != "" {
		utils.SnapInstallFromFile(nil, utils.LocalSnap)
	} else {
		utils.SnapInstallFromStore(nil, deviceRfidLlrpSnap, utils.ServiceChannel)
	}
	utils.SnapInstallFromStore(nil, "edgexfoundry", utils.PlatformChannel)

	// make sure all services are online before starting the tests
	utils.WaitPlatformOnline(nil)

	// for local build, the interface isn't auto-connected.
	// connect manually regardless
	utils.SnapConnect(nil,
		"edgexfoundry:edgex-secretstore-token",
		deviceRfidLlrpSnap+":edgex-secretstore-token",
	)

	// Start the service so that the default config gets uploaded to consul.
	// Otherwise, settings that get passed using environment variables on first start get uploaded
	// and become the default.
	utils.SnapStart(nil, deviceRfidLlrpService)
	utils.WaitServiceOnline(nil, 60, deviceRfidLlrpServicePort)

	exitCode := m.Run()

	log.Println("[TEARDOWN]")

	utils.SnapDumpLogs(nil, start, deviceRfidLlrpSnap)

	utils.SnapRemove(nil,
		deviceRfidLlrpSnap,
		"edgexfoundry",
	)

	os.Exit(exitCode)
}

func TestCommon(t *testing.T) {
	utils.TestConfig(t, deviceRfidLlrpSnap, utils.Config{
		TestChangePort: utils.ConfigChangePort{
			App:                      deviceRfidApp,
			DefaultPort:              deviceRfidLlrpServicePort,
			TestLegacyEnvConfig:      utils.FullConfigTest,
			TestAppConfig:            true,
			TestGlobalConfig:         true,
			TestMixedGlobalAppConfig: utils.FullConfigTest,
		},
	})

	utils.TestNet(t, deviceRfidLlrpSnap, utils.Net{
		StartSnap:        true,
		TestOpenPorts:    []string{deviceRfidLlrpServicePort},
		TestBindLoopback: []string{deviceRfidLlrpServicePort},
	})

	utils.TestPackaging(t, deviceRfidLlrpSnap, utils.Packaging{
		TestSemanticSnapVersion: true,
	})
}
