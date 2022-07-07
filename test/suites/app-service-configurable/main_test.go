package test

import (
	"edgex-snap-testing/test/utils"
	"log"
	"os"
	"testing"
	"time"
)

const (
	ascSnap                       = "edgex-app-service-configurable"
	ascApp                        = "app-service-configurable"
	ascService                    = ascSnap + "." + ascApp
	defaultTestProfile            = "rules-engine"
	defaultTestProfileServicePort = "59701"
)

func TestMain(m *testing.M) {
	start := time.Now()

	log.Println("[SETUP]")

	// start clean
	utils.SnapRemove(nil,
		ascSnap,
		"edgexfoundry",
	)

	// install the app-service-configurable snap before edgexfoundry
	// to catch build error sooner and stop
	if utils.LocalSnap != "" {
		utils.SnapInstallFromFile(nil, utils.LocalSnap)
	} else {
		utils.SnapInstallFromStore(nil, ascSnap, utils.ServiceChannel)
	}
	utils.SnapInstallFromStore(nil, "edgexfoundry", utils.PlatformChannel)

	// make sure all services are online before starting the tests
	utils.WaitPlatformOnline(nil)

	// for local build, the interface isn't auto-connected.
	// connect manually regardless
	utils.SnapConnect(nil,
		"edgexfoundry:edgex-secretstore-token",
		ascSnap+":edgex-secretstore-token",
	)

	// set profile to rules engine
	utils.SnapSet(nil, ascSnap, "profile", defaultTestProfile)

	// Start the service so that the default config gets uploaded to consul.
	// Otherwise, settings that get passed using environment variables on first start get uploaded
	// and become the default.
	utils.SnapStart(nil, ascService)
	utils.WaitServiceOnline(nil, 60, defaultTestProfileServicePort)

	exitCode := m.Run()

	log.Println("[TEARDOWN]")

	utils.SnapDumpLogs(nil, start, ascSnap)

	utils.SnapRemove(nil,
		ascSnap,
		"edgexfoundry",
	)

	os.Exit(exitCode)
}

func TestCommon(t *testing.T) {
	utils.TestConfig(t, ascSnap, utils.Config{
		TestChangePort: utils.ConfigChangePort{
			App:                      ascApp,
			DefaultPort:              defaultTestProfileServicePort,
			TestLegacyEnvConfig:      utils.FullConfigTest,
			TestAppConfig:            true,
			TestGlobalConfig:         true,
			TestMixedGlobalAppConfig: utils.FullConfigTest,
		},
	})

	utils.TestNet(t, ascSnap, utils.Net{
		StartSnap:        true,
		TestOpenPorts:    []string{defaultTestProfileServicePort},
		TestBindLoopback: []string{defaultTestProfileServicePort},
	})

	utils.TestPackaging(t, ascSnap, utils.Packaging{
		TestSemanticSnapVersion: true,
	})
}
