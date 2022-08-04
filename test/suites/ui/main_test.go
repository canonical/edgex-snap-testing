package test

import (
	"edgex-snap-testing/test/utils"
	"log"
	"os"
	"testing"
	"time"
)

const (
	uiSnap        = "edgex-ui"
	uiServicePort = "4000"
)

func TestMain(m *testing.M) {
	teardown, err := setupServiceTest(uiSnap)
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
		TestOpenPorts:    []string{uiServicePort},
		TestBindLoopback: []string{},
	})

	utils.TestPackaging(t, uiSnap, utils.Packaging{
		TestSemanticSnapVersion: true,
	})
}

func setupServiceTest(snapName string) (teardown func(), err error) {
	log.Println("[CLEAN]")
	utils.SnapRemove(nil,
		snapName,
	)

	log.Println("[SETUP]")
	start := time.Now()

	teardown = func() {
		log.Println("[TEARDOWN]")
		utils.SnapDumpLogs(nil, start, snapName)

		utils.SnapRemove(nil,
			snapName,
		)
	}

	if utils.LocalSnap() {
		err = utils.SnapInstallFromFile(nil, utils.LocalSnapPath)
	} else {
		err = utils.SnapInstallFromStore(nil, snapName, utils.ServiceChannel)
	}
	if err != nil {
		teardown()
		return
	}

	// make sure service is online before starting the tests
	if err = utils.WaitServiceOnline(nil, 60, uiServicePort); err != nil {
		teardown()
		return
	}

	return
}
