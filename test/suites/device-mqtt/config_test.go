package test

import (
	"edgex-snap-testing/test/utils"
	"testing"
)

// Deprecated
func TestEnvConfig(t *testing.T) {
	// start clean
	utils.Exec(t, "sudo snap stop edgex-device-mqtt.device-mqtt")

	t.Run("change service port", func(t *testing.T) {
		t.Cleanup(func() {
			utils.Exec(t, "sudo snap stop edgex-device-mqtt.device-mqtt")
			utils.Exec(t, "sudo snap unset edgex-device-mqtt env.service.port")
		})

		const newPort = "56789"

		// make sure the port is available before using it
		utils.CheckPortAvailable(t, newPort)

		utils.SnapSet(t, snap, "env.service.port", newPort)
		utils.Exec(t, "sudo snap start edgex-device-mqtt")
		utils.WaitServiceOnline(t, newPort)
	})
}

func TestAppConfig(t *testing.T) {
	t.Skip("TODO")
}
