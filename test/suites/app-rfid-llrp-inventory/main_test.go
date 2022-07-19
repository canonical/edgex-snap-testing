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
	start := time.Now()

	log.Println("[SETUP]")

	// start clean
	utils.SnapRemove(nil,
		appRfidLlrpSnap,
		"edgexfoundry",
	)

	// install the app snap before edgexfoundry
	// to catch build error sooner and stop
	if utils.LocalSnap != "" {
		utils.SnapInstallFromFile(nil, utils.LocalSnap)

		// for local build, the interface isn't auto-connected.
		// connect manually
		utils.SnapConnect(nil,
			"edgexfoundry:edgex-secretstore-token",
			appRfidLlrpSnap+":edgex-secretstore-token",
		)
	} else {
		utils.SnapInstallFromStore(nil, appRfidLlrpSnap, utils.ServiceChannel)
	}
	utils.SnapInstallFromStore(nil, "edgexfoundry", utils.PlatformChannel)

	// make sure all services are online before starting the tests
	utils.WaitPlatformOnline(nil)

	exitCode := m.Run()

	log.Println("[TEARDOWN]")

	utils.SnapDumpLogs(nil, start, appRfidLlrpSnap)

	utils.SnapRemove(nil,
		appRfidLlrpSnap,
		"edgexfoundry",
	)

	os.Exit(exitCode)
}

func TestCommon(t *testing.T) {
	utils.TestSecret(t, utils.Secret{
		TestSecretToken: true,
		Snap:            appRfidLlrpSnap,
		App:             appRfidLlrpApp,
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
