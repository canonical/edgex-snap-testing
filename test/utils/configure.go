package utils

import "testing"

func TestEnvChangeServicePort(t *testing.T, snapName, defaultServicePort string) {
	const envServicePort = "env.service.port"
	const newPort = "56789"

	t.Cleanup(func() {
		SnapStop(t, snapName)
		SnapUnset(t, snapName, envServicePort)
	})

	// make sure the port is available before using it
	RequirePortAvailable(t, newPort)

	// check if service port can be changed
	SnapStop(t, snapName)
	SnapSet(t, snapName, envServicePort, newPort)
	SnapStart(t, snapName)
	WaitServiceOnline(t, 60, newPort)

	// check if service port can be unset and revert to the default
	SnapStop(t, snapName)
	SnapUnset(t, snapName, envServicePort)
	SnapStart(t, snapName)
	WaitServiceOnline(t, 60, defaultServicePort)
}
