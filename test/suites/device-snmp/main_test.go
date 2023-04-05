package test

import (
	"edgex-snap-testing/test/utils"
	"log"
	"os"
	"testing"
)

const (
	deviceSnmpSnap = "edgex-device-snmp"
	deviceSnmpApp  = "device-snmp"
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
			DefaultPort:              utils.ServicePort(deviceSnmpApp),
			TestAppConfig:            true,
			TestGlobalConfig:         true,
			TestMixedGlobalAppConfig: utils.FullConfigTest,
		},
		TestAutoStart: true,
	})

	utils.TestNet(t, deviceSnmpSnap, utils.Net{
		StartSnap:        true,
		TestOpenPorts:    []string{utils.ServicePort(deviceSnmpApp)},
		TestBindLoopback: []string{utils.ServicePort(deviceSnmpApp)},
	})

	utils.TestPackaging(t, deviceSnmpSnap, utils.Packaging{
		TestSemanticSnapVersion: true,
	})
}
