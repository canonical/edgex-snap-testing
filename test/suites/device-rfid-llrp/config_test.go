package test

import (
	"edgex-snap-testing/test/utils"
	"testing"
)

// Deprecated
func TestEnvConfig(t *testing.T) {

	t.Run("change service port", func(t *testing.T) {
		t.Cleanup(func() {
			utils.SnapStop(t, deviceRfidLlrpService)
			utils.SnapUnset(t, deviceRfidLlrpSnap, "env.service.port")
		})

		const newPort = "56789"

		// make sure the port is available before using it
		utils.RequirePortAvailable(t, newPort)

		utils.SnapStop(t, deviceRfidLlrpSnap)
		utils.SnapSet(t, deviceRfidLlrpSnap, "env.service.port", newPort)
		utils.SnapStart(t, deviceRfidLlrpSnap)
		utils.WaitServiceOnline(t, 60, newPort)
	})
}

func TestAppConfig(t *testing.T) {
	t.Skip("TODO")
}
