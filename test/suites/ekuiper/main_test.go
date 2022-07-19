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
	start := time.Now()

	log.Println("[SETUP]")

	// start clean
	utils.SnapRemove(nil,
		ekuiperSnap,
		"edgexfoundry",
		deviceVirtualSnap,
		ascSnap,
	)

	// install the ekuiper snap before edgexfoundry
	// to catch build error sooner and stop
	if utils.LocalSnap != "" {
		utils.SnapInstallFromFile(nil, utils.LocalSnap)

		// for local build, the interface isn't auto-connected.
		// connect manually
		utils.SnapConnect(nil,
			"edgexfoundry:edgex-secretstore-token",
			ekuiperSnap+":edgex-secretstore-token",
		)
	} else {
		utils.SnapInstallFromStore(nil, ekuiperSnap, utils.ServiceChannel)
	}
	utils.SnapInstallFromStore(nil, "edgexfoundry", utils.PlatformChannel)
	utils.SnapInstallFromStore(nil, deviceVirtualSnap, "latest/edge")
	utils.SnapInstallFromStore(nil, ascSnap, "latest/edge")

	// make sure all services are online before starting the tests
	utils.WaitPlatformOnline(nil)

	// security on (default)
	testSecretsInterface = true
	exitCode := m.Run()
	if exitCode != 0 {
		goto TEARDOWN
	}

	// security off
	utils.SnapStop(nil, "edgex-ekuiper")
	utils.SnapSet(nil, "edgexfoundry", "security-secret-store", "off")
	utils.SnapSet(nil, "edgex-ekuiper", "edgex-security", "off")
	utils.SnapSet(nil, "edgex-device-virtual", "config.edgex-security-secret-store", "false")
	utils.Exec(nil, "sudo rm /var/snap/edgex-ekuiper/current/edgex-ekuiper/secrets-token.json")

	utils.SnapStart(nil,
		ekuiperService,
		deviceVirtualSnap,
		ascSnap,
	)

	testSecretsInterface = false
	exitCode = m.Run()

TEARDOWN:
	log.Println("[TEARDOWN]")

	utils.SnapDumpLogs(nil, start, ekuiperSnap)

	utils.SnapRemove(nil,
		ekuiperSnap,
		"edgexfoundry",
		deviceVirtualSnap,
		ascSnap,
	)

	os.Exit(exitCode)
}

func TestCommon(t *testing.T) {
	utils.TestSecret(t, ekuiperSnap, ekuiperSnap, ekuiperSnap, utils.Secret{
		TestSecretsInterface: testSecretsInterface,
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
