package test

import (
	"edgex-snap-testing/test/utils"
	"log"
	"os"
	"testing"
	"time"
)

const (
	appRfidLlrpSnap               = "edgex-app-rfid-llrp-inventory"
	appRfidLlrpApp                = "app-rfid-llrp-inventory"
	appRfidLlrpServiceServicePort = "59711"
)

func TestMain(m *testing.M) {
	// start clean
	utils.SnapRemove(nil,
		appRfidLlrpSnap,
		"edgexfoundry",
	)

	log.Println("[SETUP]")
	start := time.Now()

	// install the app snap before edgexfoundry
	// to catch build error sooner and stop
	if utils.LocalSnap() {
		utils.SnapInstallFromFile(nil, utils.LocalSnapPath)
	} else {
		utils.SnapInstallFromStore(nil, appRfidLlrpSnap, utils.ServiceChannel)
	}
	utils.SnapInstallFromStore(nil, "edgexfoundry", utils.PlatformChannel)

	// make sure all services are online before starting the tests
	utils.WaitPlatformOnline(nil)

	// for local build, the interface isn't auto-connected.
	// connect manually
	if utils.LocalSnap() {
		utils.SnapConnect(nil,
			"edgexfoundry:edgex-secretstore-token",
			appRfidLlrpSnap+":edgex-secretstore-token",
		)
	}

	exitCode := m.Run()

	log.Println("[TEARDOWN]")

	utils.SnapDumpLogs(nil, start, appRfidLlrpSnap)
	utils.SnapDumpLogs(nil, start, "edgexfoundry")

	utils.SnapRemove(nil,
		appRfidLlrpSnap,
		"edgexfoundry",
	)

	os.Exit(exitCode)
}

func TestCommon(t *testing.T) {
	utils.TestContentInterfaces(t, utils.ContentInterfaces{
		TestSecretstoreToken: true,
		Snap:                 appRfidLlrpSnap,
		App:                  appRfidLlrpApp,
	})

	utils.TestConfig(t, appRfidLlrpSnap, utils.Config{
		TestChangePort: utils.ConfigChangePort{
			App:                      appRfidLlrpApp,
			DefaultPort:              appRfidLlrpServiceServicePort,
			TestLegacyEnvConfig:      utils.FullConfigTest,
			TestAppConfig:            true,
			TestGlobalConfig:         true,
			TestMixedGlobalAppConfig: utils.FullConfigTest,
		},
	})

	utils.TestNet(t, appRfidLlrpSnap, utils.Net{
		StartSnap:        true,
		TestOpenPorts:    []string{appRfidLlrpServiceServicePort},
		TestBindLoopback: []string{appRfidLlrpServiceServicePort},
	})

	utils.TestPackaging(t, appRfidLlrpSnap, utils.Packaging{
		TestSemanticSnapVersion: true,
	})
}
