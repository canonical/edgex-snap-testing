package test

import (
	"edgex-snap-testing/test/utils"
	"testing"
)

const newPort = "56789"

// Deprecated
func TestEnvConfig(t *testing.T) {
	const envServicePort = "env.device-virtual.service.port"

	t.Cleanup(func() {
		utils.SnapStop(t, snapAppName)
		utils.SnapUnset(t, platformSnap, envServicePort)
	})
	t.Run("change device-virtual service port", func(t *testing.T) {
		// make sure the port is available before using it
		utils.RequirePortAvailable(t, newPort)

		// check if service port can be changed
		utils.SnapStop(t, snapAppName)
		utils.SnapSet(t, platformSnap, envServicePort, newPort)
		utils.SnapStart(t, snapAppName)
		utils.WaitServiceOnline(t, newPort)

		// check if service port can be unset and revert to the default
		utils.SnapStop(t, snapAppName)
		utils.SnapUnset(t, platformSnap, envServicePort)
		utils.SnapStart(t, snapAppName)
		utils.WaitServiceOnline(t, deviceVirtualDefaultServicePort)
	})
}

func TestAppConfig(t *testing.T) {
	t.Skip("TODO")
}
