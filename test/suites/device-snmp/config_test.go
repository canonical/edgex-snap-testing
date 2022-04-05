package test

import (
	"edgex-snap-testing/test/utils"
	"testing"
)

// Deprecated
func TestEnvConfig(t *testing.T) {
	// start clean
	utils.SnapStop(t, deviceSnmpService)

	t.Run("change service port", func(t *testing.T) {
		t.Cleanup(func() {
			utils.SnapStop(t, deviceSnmpService)
			utils.SnapUnset(t, deviceSnmpSnap, "env.service.port")
		})

		const newPort = "56789"

		// make sure the port is available before using it
		utils.RequirePortAvailable(t, newPort)

		utils.SnapSet(t, deviceSnmpSnap, "env.service.port", newPort)
		utils.SnapStart(t, deviceSnmpSnap)
		utils.WaitServiceOnline(t, newPort)
	})
}

func TestAppConfig(t *testing.T) {
	t.Skip("TODO")
}
