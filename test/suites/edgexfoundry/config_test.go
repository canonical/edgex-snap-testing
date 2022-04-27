package test

import (
	"edgex-snap-testing/test/utils"
	"testing"
)

var FullConfigTest = true

// Deprecated
func TestEnvConfig(t *testing.T) {
	const newPort = "11111"
	const envServicePort = "env." + deviceVirtualApp + ".service.port"

	// start clean
	utils.SnapStop(t, deviceVirtualService)

	t.Cleanup(func() {
		utils.SnapUnset(t, platformSnap, envServicePort)
		utils.SnapStop(t, deviceVirtualService)
	})
	t.Run("change device-virtual service port", func(t *testing.T) {

		// make sure the port is available before using it
		utils.RequirePortAvailable(t, newPort)

		// set env. and validate the new port comes online
		utils.SnapSet(t, platformSnap, envServicePort, newPort)
		utils.SnapStart(t, deviceVirtualService)
		utils.WaitServiceOnline(t, 60, newPort)

		// unset env. and validate the default port comes online
		utils.SnapUnset(t, platformSnap, envServicePort)
		utils.SnapRestart(t, deviceVirtualService)
		utils.WaitServiceOnline(t, 60, deviceVirtualDefaultServicePort)

	})
}

func TestAppConfig(t *testing.T) {
	t.Cleanup(func() {
		utils.SnapStop(t, deviceVirtualService)
	})

	utils.SnapStart(t, deviceVirtualService)
	utils.SetAppConfig(t, platformSnap, deviceVirtualApp, deviceVirtualDefaultServicePort)
}

func TestGlobalConfig(t *testing.T) {
	t.Cleanup(func() {
		utils.SnapStop(t, deviceVirtualService)
	})

	utils.SnapStart(t, deviceVirtualService)
	utils.SetGlobalConfig(t, platformSnap, deviceVirtualApp, deviceVirtualDefaultServicePort)
}

func TestMixedConfig(t *testing.T) {
	t.Cleanup(func() {
		utils.SnapStop(t, deviceVirtualService)
	})

	utils.SnapStart(t, deviceVirtualService)
	utils.SetMixedConfig(t, platformSnap, deviceVirtualApp, deviceVirtualDefaultServicePort)
}
