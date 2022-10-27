package test

import (
	"edgex-snap-testing/test/utils"
	"log"
	"os"
	"testing"
)

const (
	appRfidLlrpSnap               = "edgex-app-rfid-llrp-inventory"
	appRfidLlrpApp                = "app-rfid-llrp-inventory"
	appRfidLlrpServiceServicePort = "59711"
)

func TestMain(m *testing.M) {
	teardown, err := utils.SetupServiceTests(appRfidLlrpSnap)
	if err != nil {
		log.Fatalf("Failed to setup tests: %s", err)
	}

	code := m.Run()
	teardown()

	os.Exit(code)
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
		TestAutoStart: true,
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
