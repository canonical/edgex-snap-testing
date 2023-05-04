package test

import (
	"edgex-snap-testing/test/utils"
	"log"
	"os"
	"testing"
	"time"
)

const (
	platformSnap = "edgexfoundry"

	ekuiperSnap       = "edgex-ekuiper"
	ekuiperApp        = "ekuiper"
	ekuiperRestfulApi = "ekuiper/rest-api"

	deviceVirtualSnap = "edgex-device-virtual"
	deviceVirtualApp  = "device-virtual"
)

var (
	ekuiperPort        = utils.ServicePort(ekuiperApp)
	ekuiperRestfulPort = utils.ServicePort(ekuiperRestfulApi)
)

var testSecretsInterface bool

func TestMain(m *testing.M) {
	teardown, err := setup()
	if err != nil {
		log.Fatalf("Failed to setup tests: %s", err)
	}

	// set profile to rules engine
	testSecretsInterface = true

	code := m.Run()
	teardown()

	os.Exit(code)
}

func TestCommon(t *testing.T) {
	utils.TestContentInterfaces(t, utils.ContentInterfaces{
		TestSecretstoreToken: testSecretsInterface,
		Snap:                 ekuiperSnap,
		App:                  ekuiperSnap,
	})

	utils.TestConfig(t, ekuiperSnap, utils.Config{
		TestAutoStart: true,
	})

	utils.TestNet(t, ekuiperSnap, utils.Net{
		StartSnap:        true,
		TestOpenPorts:    []string{ekuiperPort, ekuiperRestfulPort},
		TestBindLoopback: []string{ekuiperPort, ekuiperRestfulPort},
	})

	utils.TestPackaging(t, ekuiperSnap, utils.Packaging{
		TestSemanticSnapVersion: true,
	})
}

func setup() (teardown func(), err error) {
	log.Println("[CLEAN]")
	utils.SnapRemove(nil,
		ekuiperSnap,
		platformSnap,
		deviceVirtualSnap,
	)

	log.Println("[SETUP]")
	start := time.Now()

	teardown = func() {
		log.Println("[TEARDOWN]")

		utils.SnapDumpLogs(nil, start, ekuiperSnap)
		utils.SnapDumpLogs(nil, start, platformSnap)
		utils.SnapDumpLogs(nil, start, deviceVirtualSnap)

		log.Println("Removing installed snap:", !utils.SkipTeardownRemoval)
		if !utils.SkipTeardownRemoval {
			utils.SnapRemove(nil,
				ekuiperSnap,
				platformSnap,
				deviceVirtualSnap,
			)
		}
	}

	// install the ekuiper snap before edgexfoundry
	// to catch build error sooner and stop
	if utils.LocalSnap() {
		err = utils.SnapInstallFromFile(nil, utils.LocalSnapPath)
	} else {
		err = utils.SnapInstallFromStore(nil, ekuiperSnap, utils.ServiceChannel)
	}
	if err != nil {
		teardown()
		return
	}

	if err = utils.SnapInstallFromStore(nil, platformSnap, utils.PlatformChannel); err != nil {
		teardown()
		return
	}

	if err = utils.SnapInstallFromStore(nil, deviceVirtualSnap, utils.ServiceChannel); err != nil {
		teardown()
		return
	}

	// for local build, the interface isn't auto-connected.
	// connect manually
	if utils.LocalSnap() {
		if err = utils.SnapConnectSecretstoreToken(nil, ekuiperSnap); err != nil {
			teardown()
			return
		}
	}

	utils.SnapStart(nil, platformSnap)

	// make sure all services are online before starting the tests
	if err = utils.WaitPlatformOnline(nil); err != nil {
		teardown()
		return
	}

	return
}
