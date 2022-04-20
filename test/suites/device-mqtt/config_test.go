package test

import (
	"edgex-snap-testing/test/utils"
	"testing"
)

// Deprecated
func TestEnvConfig(t *testing.T) {
	// start clean
	utils.SnapStop(t, deviceMqttSnap)

	t.Run("change service port", func(t *testing.T) {
		t.Cleanup(func() {
			utils.SnapUnset(t, deviceMqttSnap, "env.service.port")
			utils.SnapUnset(t, deviceMqttSnap, "env")
			utils.SnapRestart(t, deviceMqttSnap)
		})

		const newPort = "56789"

		// make sure the port is available before using it
		utils.RequirePortAvailable(t, newPort)

		utils.SnapStart(t, deviceMqttSnap)
		utils.SnapSet(t, deviceMqttSnap, "env.service.port", newPort)
		utils.SnapRestart(t, deviceMqttSnap)

		utils.WaitServiceOnline(t, 60, newPort)
	})
}

func TestAppConfig(t *testing.T) {
	// start clean
	utils.SnapStop(t, deviceMqttSnap)

	t.Run("use apps. and config. for the same option", func(t *testing.T) {
		t.Cleanup(func() {
			utils.SnapUnset(t, deviceMqttSnap, "apps.device-mqtt.config.service.port")
			utils.SnapUnset(t, deviceMqttSnap, "config.service.port")
			utils.SnapRestart(t, deviceMqttSnap)
		})

		const newPort = "8888"

		// make sure the port is available before using it
		utils.RequirePortAvailable(t, newPort)

		utils.SnapStart(t, deviceMqttSnap)
		utils.SnapSet(t, deviceMqttSnap, "apps.device-mqtt.config.service.port", newPort)
		utils.SnapSet(t, deviceMqttSnap, "config.service.port", newPort)
		utils.SnapRestart(t, deviceMqttSnap)

		utils.WaitServiceOnline(t, 60, newPort)
	})

	t.Run("use apps. and config. for different options", func(t *testing.T) {
		t.Cleanup(func() {
			utils.SnapUnset(t, deviceMqttSnap, "apps.device-mqtt.config.service.port")
			utils.SnapUnset(t, deviceMqttSnap, "config.service.port")
			utils.SnapRestart(t, deviceMqttSnap)
		})

		const newAppPort = "11111"
		const newConfigPort = "22222"

		// make sure the ports are available before using it
		utils.RequirePortAvailable(t, newAppPort)
		utils.RequirePortAvailable(t, newConfigPort)

		utils.SnapStart(t, deviceMqttSnap)
		utils.SnapSet(t, deviceMqttSnap, "apps.device-mqtt.config.service.port", newAppPort)
		utils.SnapSet(t, deviceMqttSnap, "config.service.port", newConfigPort)
		utils.SnapRestart(t, deviceMqttSnap)

		utils.WaitServiceOnline(t, 60, newAppPort)
	})
}
