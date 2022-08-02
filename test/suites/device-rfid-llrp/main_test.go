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
	// start clean
	utils.SnapRemove(nil,
		deviceRfidLlrpSnap,
		"edgexfoundry",
	)

	log.Println("[SETUP]")
	start := time.Now()

	// install the device snap before edgexfoundry
	// to catch build error sooner and stop
	if utils.LocalSnap() {
		utils.SnapInstallFromFile(nil, utils.LocalSnapPath)
	} else {
		utils.SnapInstallFromStore(nil, deviceRfidLlrpSnap, utils.ServiceChannel)
	}
	utils.SnapInstallFromStore(nil, "edgexfoundry", utils.PlatformChannel)

	// make sure all services are online before starting the tests
	utils.WaitPlatformOnline(nil)

	// for local build, the interface isn't auto-connected.
	// connect manually
	if utils.LocalSnap() {
		utils.SnapConnect(nil,
			"edgexfoundry:edgex-secretstore-token",
			deviceRfidLlrpSnap+":edgex-secretstore-token",
		)
	}

	exitCode := m.Run()

	log.Println("[TEARDOWN]")

	utils.SnapDumpLogs(nil, start, deviceRfidLlrpSnap)
	utils.SnapDumpLogs(nil, start, "edgexfoundry")

	utils.SnapRemove(nil,
		deviceRfidLlrpSnap,
		"edgexfoundry",
	)

	os.Exit(exitCode)
}

func TestCommon(t *testing.T) {
	utils.TestContentInterfaces(t, utils.ContentInterfaces{
		TestSecretstoreToken: true,
		Snap:                 deviceRfidLlrpSnap,
		App:                  deviceRfidLlrpApp,
	})

	utils.TestConfig(t, deviceRfidLlrpSnap, utils.Config{
		TestChangePort: utils.ConfigChangePort{
			App:                      deviceRfidLlrpApp,
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
