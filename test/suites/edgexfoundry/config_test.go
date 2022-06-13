package test

import (
	"edgex-snap-testing/test/utils"
	"testing"
)

func TestEnvConfigChangePort(t *testing.T) {
	if !utils.FullConfigTest {
		t.Skip("Full config test is disabled.")
	}

	const newPort = "11111"
	const envServicePort = "env." + supportSchedulerApp + ".service.port"

	// start clean
	utils.SnapStop(t, supportSchedulerService)

	t.Cleanup(func() {
		utils.SnapUnset(t, platformSnap, envServicePort)
		utils.SnapStop(t, supportSchedulerService)
	})

	// make sure the port is available before using it
	utils.RequirePortAvailable(t, newPort)

	// set env. and validate the new port comes online
	utils.SnapSet(t, platformSnap, envServicePort, newPort)
	utils.SnapStart(t, supportSchedulerService)
	utils.WaitServiceOnline(t, 60, newPort)

	// unset env. and validate the default port comes online
	utils.SnapUnset(t, platformSnap, envServicePort)
	utils.SnapRestart(t, supportSchedulerService)
	utils.WaitServiceOnline(t, 60, supportSchedulerServicePort)
}
