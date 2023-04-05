package test

import (
	"edgex-snap-testing/test/utils"
	"log"
	"os"
	"testing"
	"time"
)

const (
	uiSnap = "edgex-ui"
	uiApp  = "ui"
)

func TestMain(m *testing.M) {
	teardown, err := setup()
	if err != nil {
		log.Fatalf("Failed to setup tests: %s", err)
	}

	code := m.Run()
	teardown()

	os.Exit(code)
}

func TestCommon(t *testing.T) {
	utils.TestNet(t, uiSnap, utils.Net{
		StartSnap:        true,
		TestOpenPorts:    []string{utils.ServicePort(uiApp)},
		TestBindLoopback: []string{utils.ServicePort(uiApp)},
	})

	utils.TestPackaging(t, uiSnap, utils.Packaging{
		TestSemanticSnapVersion: true,
	})
}

func setup() (teardown func(), err error) {
	log.Println("[CLEAN]")
	utils.SnapRemove(nil, uiSnap)

	log.Println("[SETUP]")
	start := time.Now()

	teardown = func() {
		log.Println("[TEARDOWN]")
		utils.SnapDumpLogs(nil, start, uiSnap)

		log.Println("Removing installed snap:", !utils.SkipTeardownRemoval)
		if !utils.SkipTeardownRemoval {
			utils.SnapRemove(nil, uiSnap)
		}
	}

	if utils.LocalSnap() {
		err = utils.SnapInstallFromFile(nil, utils.LocalSnapPath)
	} else {
		err = utils.SnapInstallFromStore(nil, uiSnap, utils.ServiceChannel)
	}
	if err != nil {
		teardown()
		return
	}

	return
}
