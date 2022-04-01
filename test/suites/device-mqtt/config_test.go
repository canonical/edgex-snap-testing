package test

import (
	"edgex-snap-testing/test/utils"
	"testing"
)

// Deprecated
func TestEnvConfig(t *testing.T) {
	// start clean
	utils.SnapStop(t, deviceMqttService)

	t.Run("change service port", func(t *testing.T) {
		t.Cleanup(func() {
			utils.SnapStop(t, deviceMqttService)
			utils.SnapUnset(t, deviceMqttSnap, "env.service.port")
		})

		const newPort = "56789"

		// make sure the port is available before using it
		utils.RequirePortAvailable(t, newPort)

		utils.SnapSet(t, deviceMqttSnap, "env.service.port", newPort)
		utils.SnapStart(t, deviceMqttSnap)
		utils.RequireServiceOnline(t, newPort)
	})
}

func TestAppConfig(t *testing.T) {
	t.Skip("TODO")
}
