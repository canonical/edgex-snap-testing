package test

import (
	"edgex-snap-testing/test/utils"
	"testing"
)

// Deprecated
func TestEnvConfig(t *testing.T) {

	t.Cleanup(func() {
		utils.SnapStop(t, snapAppName)
	})
	t.Run("change device-virtual service port", func(t *testing.T) {
		const newPort = "56789"
		const envServicePort = "env.device-virtual.service.port"

		// make sure the port is available before using it
		utils.RequirePortAvailable(t, newPort)

		// check if service port can be changed
		utils.SnapStop(t, snapAppName)
		utils.SnapSet(t, platformSnap, envServicePort, newPort)
		utils.SnapStart(t, snapAppName)
		utils.WaitServiceOnline(t, 60, newPort)

		// check if service port can be unset and revert to the default
		utils.SnapStop(t, snapAppName)
		utils.SnapUnset(t, platformSnap, envServicePort)
		utils.SnapStart(t, snapAppName)
		utils.WaitServiceOnline(t, 60, deviceVirtualDefaultServicePort)

		utils.SnapStop(t, snapAppName)
		utils.SnapUnset(t, platformSnap, envServicePort)
	})
}

func TestAppConfig(t *testing.T) {
	t.Skip("TODO")
}
