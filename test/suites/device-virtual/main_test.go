package test

import (
	"edgex-snap-testing/test/utils"
	"log"
	"os"
	"testing"
	"time"
)

const (
	deviceVirtualSnap        = "edgex-device-virtual"
	deviceVirtualApp         = "device-virtual"
	deviceVirtualServicePort = "59900"
)

func main(m *testing.M) (int, error) {
	log.Println("[CLEAN]")
	utils.SnapRemove(nil,
		deviceVirtualSnap,
		"edgexfoundry",
	)

	log.Println("[SETUP]")

	// add this to the bottom of the defer stack to remove after collecting logs
	defer utils.SnapRemove(nil,
		deviceVirtualSnap,
		"edgexfoundry",
	)

	start := time.Now()
	defer utils.SnapDumpLogs(nil, start, deviceVirtualSnap)
	defer utils.SnapDumpLogs(nil, start, "edgexfoundry")

	// install the device snap before edgexfoundry
	// to catch build error sooner and stop
	if utils.LocalSnap() {
		if err := utils.SnapInstallFromFile(nil, utils.LocalSnapPath); err != nil {
			return 0, err
		}
	} else {
		if err := utils.SnapInstallFromStore(nil, deviceVirtualSnap, utils.ServiceChannel); err != nil {
			return 0, err
		}
	}

	if err := utils.SnapInstallFromStore(nil, "edgexfoundry", utils.PlatformChannel); err != nil {
		return 0, err
	}

	// make sure all services are online before starting the tests
	if err := utils.WaitPlatformOnline(nil); err != nil {
		return 0, err
	}

	// for local build, the interface isn't auto-connected.
	// connect manually
	if utils.LocalSnap() {
		if err := utils.SnapConnectSecretstoreToken(nil, deviceVirtualSnap); err != nil {
			return 0, err
		}
	}

	log.Println("[START]")
	return m.Run(), nil
}

func TestMain(m *testing.M) {
	code, err := main(m)
	if err != nil {
		log.Fatalf("Failed to run tests: %s", err)
	}
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
