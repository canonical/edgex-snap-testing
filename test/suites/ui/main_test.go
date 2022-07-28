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
	// start clean
	utils.SnapRemove(nil,
		uiSnap,
	)

	log.Println("[SETUP]")
	start := time.Now()

	if utils.LocalSnap != "" {
		utils.SnapInstallFromFile(nil, utils.LocalSnap)
	} else {
		utils.SnapInstallFromStore(nil, uiSnap, utils.ServiceChannel)
	}

	exitCode := m.Run()

	log.Println("[TEARDOWN]")

	utils.SnapDumpLogs(nil, start, uiSnap)
	utils.SnapDumpLogs(nil, start, "edgexfoundry")

	utils.SnapRemove(nil,
		uiSnap,
	)

	os.Exit(exitCode)
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
