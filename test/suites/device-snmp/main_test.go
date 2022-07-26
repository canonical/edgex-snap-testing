package test

import (
	"edgex-snap-testing/test/utils"
	"log"
	"os"
	"testing"
	"time"
)

const (
	deviceSnmpSnap        = "edgex-device-snmp"
	deviceSnmpApp         = "device-snmp"
	deviceSnmpServicePort = "59993"
)

func TestMain(m *testing.M) {
	start := time.Now()

	log.Println("[SETUP]")

	// start clean
	utils.SnapRemove(nil,
		deviceSnmpSnap,
		"edgexfoundry",
	)

	// install the device snap before edgexfoundry
	// to catch build error sooner and stop
	if utils.LocalSnap != "" {
		utils.SnapInstallFromFile(nil, utils.LocalSnap)

		// for local build, the interface isn't auto-connected.
		// connect manually
		utils.SnapConnect(nil,
			"edgexfoundry:edgex-secretstore-token",
			deviceSnmpSnap+":edgex-secretstore-token",
		)
	} else {
		utils.SnapInstallFromStore(nil, deviceSnmpSnap, utils.ServiceChannel)
	}
	utils.SnapInstallFromStore(nil, "edgexfoundry", utils.PlatformChannel)

	// make sure all services are online before starting the tests
	utils.WaitPlatformOnline(nil)

	exitCode := m.Run()

	log.Println("[TEARDOWN]")

	utils.SnapDumpLogs(nil, start, deviceSnmpSnap)

	utils.SnapRemove(nil,
		deviceSnmpSnap,
		"edgexfoundry",
	)

	os.Exit(exitCode)
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
