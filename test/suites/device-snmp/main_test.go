package test

import (
	"edgex-snap-testing/test/utils"
	"log"
	"os"
	"testing"
)

const (
	deviceSnmpSnap        = "edgex-device-snmp"
	deviceSnmpApp         = "device-snmp"
	deviceSnmpServicePort = "59993"
)

func TestMain(m *testing.M) {
	teardown, err := utils.SetupServiceTests(deviceSnmpSnap)
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
		Snap:                 deviceSnmpSnap,
		App:                  deviceSnmpApp,
	})

	utils.TestConfig(t, deviceSnmpSnap, utils.Config{
		TestChangePort: utils.ConfigChangePort{
			App:                      deviceSnmpApp,
			DefaultPort:              deviceSnmpServicePort,
			TestLegacyEnvConfig:      utils.FullConfigTest,
			TestAppConfig:            true,
			TestGlobalConfig:         true,
			TestMixedGlobalAppConfig: utils.FullConfigTest,
		},
	})

	utils.TestNet(t, deviceSnmpSnap, utils.Net{
		StartSnap:        true,
		TestOpenPorts:    []string{deviceSnmpServicePort},
		TestBindLoopback: []string{deviceSnmpServicePort},
	})

	utils.TestPackaging(t, deviceSnmpSnap, utils.Packaging{
		TestSemanticSnapVersion: true,
	})
}
