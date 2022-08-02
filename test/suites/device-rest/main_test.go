package test

import (
	"edgex-snap-testing/test/utils"
	"log"
	"os"
	"testing"
)

const (
	deviceRestSnap        = "edgex-device-rest"
	deviceRestApp         = "device-rest"
	deviceRestServicePort = "59986"
)

func TestMain(m *testing.M) {
	code, err := utils.RunDeviceTests(m, deviceRestSnap)
	if err != nil {
		log.Fatalf("Failed to run tests: %s", err)
	}
	os.Exit(code)
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
