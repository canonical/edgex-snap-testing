package test

import (
	"edgex-snap-testing/test/utils"
	"testing"
)

// Deprecated
func TestEnvConfig(t *testing.T) {
	const newPort = "11111"
	const envServicePort = "env." + deviceVirtualApp + ".service.port"

	// start clean
	utils.SnapStop(t, platformSnap)

	t.Cleanup(func() {
		utils.SnapUnset(t, platformSnap, envServicePort)
		utils.SnapStop(t, platformSnap)
	})
	t.Run("change device-virtual service port", func(t *testing.T) {

		// make sure the port is available before using it
		utils.RequirePortAvailable(t, newPort)

		// set env. and validate the new port comes online
		utils.SnapSet(t, platformSnap, envServicePort, newPort)
		utils.SnapStart(t, platformSnap)

		utils.WaitServiceOnline(t, 60, newPort)

		// unset env. and validate the default port comes online
		utils.SnapUnset(t, platformSnap, envServicePort)
		utils.SnapRestart(t, snapAppName)
		utils.WaitServiceOnline(t, 60, deviceVirtualDefaultServicePort)

	})
}

func TestAppConfig(t *testing.T) {
	t.Cleanup(func() {
		utils.SnapStop(t, platformSnap)
	})

	utils.SnapStart(t, platformSnap)
	utils.SetAppConfig(t, platformSnap, snapAppName, deviceVirtualApp, deviceVirtualDefaultServicePort)
}

func TestGlobalConfig(t *testing.T) {
	t.Cleanup(func() {
		utils.SnapStop(t, platformSnap)
	})
	utils.SnapStart(t, platformSnap)

	utils.SetGlobalConfig(t, platformSnap, snapAppName, deviceVirtualDefaultServicePort)
}

func TestMixedConfig(t *testing.T) {
	utils.FullConfigTest = true

	t.Cleanup(func() {
		utils.SnapStop(t, platformSnap)
	})
	utils.SnapStart(t, platformSnap)

	utils.SetMixedConfig(t, platformSnap, snapAppName, deviceVirtualApp, deviceVirtualDefaultServicePort)
}
