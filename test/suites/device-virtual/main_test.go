package test

import (
	"edgex-snap-testing/test/utils"
	"log"
	"os"
	"testing"
)

const (
	deviceVirtualSnap        = "edgex-device-virtual"
	deviceVirtualApp         = "device-virtual"
	deviceVirtualServicePort = "59900"
)

func TestMain(m *testing.M) {
	teardown, err := utils.SetupServiceTests(deviceVirtualSnap)
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
		TestAutoStart: true,
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
