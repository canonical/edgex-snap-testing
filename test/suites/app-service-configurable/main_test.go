package test

import (
	"edgex-snap-testing/test/utils"
	"log"
	"os"
	"testing"
)

const (
	ascSnap            = "edgex-app-service-configurable"
	ascApp             = "app-service-configurable"
	ascService         = ascSnap + "." + ascApp
	defaultTestProfile = "rules-engine"
)

func TestMain(m *testing.M) {
	teardown, err := utils.SetupServiceTests(ascSnap)
	if err != nil {
		log.Fatalf("Failed to setup tests: %s", err)
	}

	// set profile to rules engine
	utils.SnapSet(nil, ascSnap, "profile", defaultTestProfile)

	code := m.Run()
	teardown()

	os.Exit(code)
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
			DefaultPort:              utils.ServicePort(ascApp),
			TestAppConfig:            true,
			TestGlobalConfig:         true,
			TestMixedGlobalAppConfig: utils.FullConfigTest,
		},
		TestAutoStart: true,
	})

	utils.TestNet(t, ascSnap, utils.Net{
		StartSnap:        true,
		TestOpenPorts:    []string{utils.ServicePort(ascApp)},
		TestBindLoopback: []string{utils.ServicePort(ascApp)},
	})

	utils.TestPackaging(t, ascSnap, utils.Packaging{
		TestSemanticSnapVersion: true,
	})
}
