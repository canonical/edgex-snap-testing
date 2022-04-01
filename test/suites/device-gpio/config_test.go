package test

import (
	"edgex-snap-testing/test/utils"
	"testing"
)

// Deprecated
func TestEnvConfig(t *testing.T) {
	// start clean
	utils.SnapStop(t, deviceGpioService)

	t.Run("change service port", func(t *testing.T) {
		t.Cleanup(func() {
			utils.SnapStop(t, deviceGpioService)
			utils.SnapUnset(t, deviceGpioSnap, "env.service.port")
		})

		const newPort = "56789"

		// make sure the port is available before using it
		utils.CheckPortAvailable(t, newPort)

		utils.SnapSet(t, deviceGpioSnap, "env.service.port", newPort)
		utils.SnapStart(t, deviceGpioSnap)
		utils.WaitServiceOnline(t, newPort)
	})
}

func TestAppConfig(t *testing.T) {
	t.Skip("TODO")
}
