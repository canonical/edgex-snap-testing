package test

import (
	"edgex-snap-testing/test/utils"
	"log"
	"os"
	"testing"
	"time"
)

const (
	ekuiperSnap           = "edgex-ekuiper"
	ekuiperApp            = "kuiper"
	ekuiperService        = ekuiperSnap + "." + ekuiperApp
	ekuiperServerPort     = "20498"
	ekuiperRestfulApiPort = "59720"

	deviceVirtualSnap = "edgex-device-virtual"
	deviceVirtualPort = "59900"

	ascSnap = "edgex-app-service-configurable"
)

var testSecretsInterface bool

func TestMain(m *testing.M) {
	teardown, err := setup()
	if err != nil {
		log.Fatalf("Failed to setup tests: %s", err)
	}

	// set profile to rules engine
	utils.SnapSet(nil, ascSnap, "profile", "rules-engine")

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

	utils.TestNet(t, ekuiperSnap, utils.Net{
		StartSnap:        true,
		TestOpenPorts:    []string{ekuiperServerPort, ekuiperRestfulApiPort},
		TestBindLoopback: []string{ekuiperServerPort, ekuiperRestfulApiPort},
	})

	utils.TestPackaging(t, ekuiperSnap, utils.Packaging{
		TestSemanticSnapVersion: true,
	})
}

func setup() (teardown func(), err error) {
	log.Println("[CLEAN]")
	utils.SnapRemove(nil,
		ekuiperSnap,
		"edgexfoundry",
		deviceVirtualSnap,
		ascSnap,
	)

	log.Println("[SETUP]")
	start := time.Now()

	teardown = func() {
		log.Println("[TEARDOWN]")

		utils.SnapDumpLogs(nil, start, ekuiperSnap)
		utils.SnapDumpLogs(nil, start, "edgexfoundry")
		utils.SnapDumpLogs(nil, start, deviceVirtualSnap)
		utils.SnapDumpLogs(nil, start, ascSnap)

		utils.SnapRemove(nil,
			ekuiperSnap,
			"edgexfoundry",
			deviceVirtualSnap,
			ascSnap,
		)
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

	if err = utils.SnapInstallFromStore(nil, "edgexfoundry", utils.PlatformChannel); err != nil {
		teardown()
		return
	}

	if err = utils.SnapInstallFromStore(nil, deviceVirtualSnap, utils.ServiceChannel); err != nil {
		teardown()
		return
	}

	if err = utils.SnapInstallFromStore(nil, ascSnap, utils.ServiceChannel); err != nil {
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

	// make sure all services are online before starting the tests
	if err = utils.WaitPlatformOnline(nil); err != nil {
		teardown()
		return
	}

	return
}
