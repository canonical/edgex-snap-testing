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
	deviceMqttService     = deviceMqttSnap + "." + deviceMqttApp
	deviceMqttServicePort = "59982"
)

var start = time.Now()

func TestMain(m *testing.M) {

	log.Println("[SETUP]")

	// start clean
	utils.SnapRemove(nil,
		deviceMqttSnap,
		"edgexfoundry",
	)

	// install the device snap before edgexfoundry
	// to catch build error sooner and stop
	if utils.LocalSnap != "" {
		utils.SnapInstallFromFile(nil, utils.LocalSnap)
	} else {
		utils.SnapInstallFromStore(nil, deviceMqttSnap, utils.ServiceChannel)
	}
	utils.SnapInstallFromStore(nil, "edgexfoundry", utils.PlatformChannel)

	// make sure all services are online before starting the tests
	utils.WaitPlatformOnline(nil)

	// for local build, the interface isn't auto-connected.
	// connect manually regardless
	utils.SnapConnect(nil,
		"edgexfoundry:edgex-secretstore-token",
		deviceMqttSnap+":edgex-secretstore-token",
	)

	exitCode := m.Run()

	log.Println("[TEARDOWN]")

	utils.SnapDumpLogs(nil, start, deviceMqttSnap)

	utils.SnapRemove(nil,
		deviceMqttSnap,
		"edgexfoundry",
	)

	os.Exit(exitCode)
}

func TestCommon(t *testing.T) {
	utils.TestCommon(t, utils.Params{
		Snap: deviceMqttSnap,
		Config: utils.Config{
			TestChangePort: utils.ConfigChangePort{
				App:                      deviceMqttApp,
				DefaultPort:              deviceMqttServicePort,
				TestLegacyEnvConfig:      utils.FullConfigTest,
				TestAppConfig:            true,
				TestGlobalConfig:         true,
				TestMixedGlobalAppConfig: utils.FullConfigTest,
			},
		},
		Net: utils.Net{
			StartSnap:        true,
			TestOpenPorts:    []string{deviceMqttServicePort},
			TestBindLoopback: []string{deviceMqttServicePort},
		},
		Packaging: utils.Packaging{
			TestSemanticSnapVersion: true,
		},
	})
}
