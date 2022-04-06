package utils

import "testing"

func TestEnvChangeServicePort(t *testing.T, snapName, appName, defaultServicePort string) {
	const newPort = "56789"

	var envServicePort = "env.service.port"

	if appName != "" {
		envServicePort = "env." + appName + ".service.port"
	}

	t.Cleanup(func() {
		if appName != "" {
			SnapStop(t, snapName+"."+appName)
		} else {
			SnapStop(t, snapName)
		}

		SnapUnset(t, snapName, envServicePort)
	})

	// make sure the port is available before using it
	RequirePortAvailable(t, newPort)

	// check if service port can be changed
	if appName != "" {
		SnapStop(t, snapName+"."+appName)
	} else {
		SnapStop(t, snapName)
	}

	SnapSet(t, snapName, envServicePort, newPort)

	if appName != "" {
		SnapStart(t, snapName+"."+appName)
	} else {
		SnapStart(t, snapName)
	}

	WaitServiceOnline(t, newPort)

	// check if service port can be unset and revert to the default
	if appName != "" {
		SnapStop(t, snapName+"."+appName)
	} else {
		SnapStop(t, snapName)
	}

	SnapUnset(t, snapName, envServicePort)

	if appName != "" {
		SnapStart(t, snapName+"."+appName)
	} else {
		SnapStart(t, snapName)
	}

	WaitServiceOnline(t, defaultServicePort)
}
