package test

import (
	"edgex-snap-testing/test/utils"
	"log"
	"os"
	"testing"
	"time"
)

const (
	deviceMqttSnap        = "edgex-device-mqtt"
	deviceMqttApp         = "device-mqtt"
	deviceMqttServicePort = "59982"
)

func TestMain(m *testing.M) {
	// start clean
	utils.SnapRemove(nil,
		deviceMqttSnap,
		"edgexfoundry",
	)

	log.Println("[SETUP]")
	start := time.Now()

	// install the device snap before edgexfoundry
	// to catch build error sooner and stop
	if utils.LocalSnap != "" {
		utils.SnapInstallFromFile(nil, utils.LocalSnap)

		// for local build, the interface isn't auto-connected.
		// connect manually
		utils.SnapConnect(nil,
			"edgexfoundry:edgex-secretstore-token",
			deviceMqttSnap+":edgex-secretstore-token",
		)
	} else {
		utils.SnapInstallFromStore(nil, deviceMqttSnap, utils.ServiceChannel)
	}
	utils.SnapInstallFromStore(nil, "edgexfoundry", utils.PlatformChannel)

	// make sure all services are online before starting the tests
	utils.WaitPlatformOnline(nil)

	exitCode := m.Run()

	log.Println("[TEARDOWN]")

	utils.SnapDumpLogs(nil, start, deviceMqttSnap)
	utils.SnapDumpLogs(nil, start, "edgexfoundry")

	utils.SnapRemove(nil,
		deviceMqttSnap,
		"edgexfoundry",
	)

	os.Exit(exitCode)
}

func TestCommon(t *testing.T) {
	utils.TestContentInterfaces(t, utils.ContentInterfaces{
		TestSecretstoreToken: true,
		Snap:                 deviceMqttSnap,
		App:                  deviceMqttApp,
	})

	utils.TestConfig(t, deviceMqttSnap, utils.Config{
		TestChangePort: utils.ConfigChangePort{
			App:                      deviceMqttApp,
			DefaultPort:              deviceMqttServicePort,
			TestLegacyEnvConfig:      utils.FullConfigTest,
			TestAppConfig:            true,
			TestGlobalConfig:         true,
			TestMixedGlobalAppConfig: utils.FullConfigTest,
		},
	})

	utils.TestNet(t, deviceMqttSnap, utils.Net{
		StartSnap:        true,
		TestOpenPorts:    []string{deviceMqttServicePort},
		TestBindLoopback: []string{deviceMqttServicePort},
	})

	utils.TestPackaging(t, deviceMqttSnap, utils.Packaging{
		TestSemanticSnapVersion: true,
	})
}
