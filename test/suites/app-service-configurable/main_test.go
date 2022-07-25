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

		// for local build, the interface isn't auto-connected.
		// connect manually
		utils.SnapConnect(nil,
			"edgexfoundry:edgex-secretstore-token",
			ascSnap+":edgex-secretstore-token",
		)

	} else {
		utils.SnapInstallFromStore(nil, ascSnap, utils.ServiceChannel)
	}
	utils.SnapInstallFromStore(nil, "edgexfoundry", utils.PlatformChannel)

	// make sure all services are online before starting the tests
	utils.WaitPlatformOnline(nil)

	// set profile to rules engine
	utils.SnapSet(nil, ascSnap, "profile", defaultTestProfile)

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
	utils.TestContentInterfaces(t, utils.ContentInterfaces{
		TestSecretstoreToken: true,
		Snap:                 ascSnap,
		App:                  "app-" + defaultTestProfile,
	})

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
