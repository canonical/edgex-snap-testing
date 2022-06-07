package utils

import "testing"

const serviceWaitTimeout = 60 // seconds

func SetEnvConfig(t *testing.T, snap, app, servicePort string) {
	service := snap + "." + app
	if !FullConfigTest {
		t.Skip("Full config test is disabled.")
	}
	// start clean
	SnapStop(t, service)

	t.Run("change service port", func(t *testing.T) {
		t.Cleanup(func() {
			SnapUnset(t, snap, "env")
			SnapStop(t, service)
		})

		const newPort = "11111"

		// make sure the port is available before using it
		RequirePortAvailable(t, newPort)

		// set env. and validate the new port comes online
		SnapSet(t, snap, "env.service.port", newPort)
		SnapStart(t, service)

		WaitServiceOnline(t, serviceWaitTimeout, newPort)

		// unset env. and validate the default port comes online
		SnapUnset(t, snap, "env.service.port")
		SnapRestart(t, service)
		WaitServiceOnline(t, serviceWaitTimeout, servicePort)
	})

}

func SetAppConfig(t *testing.T, snap, app, servicePort string) {
	service := snap + "." + app

	// start clean
	SnapStop(t, service)

	t.Run("set and unset apps.", func(t *testing.T) {
		t.Cleanup(func() {
			SnapUnset(t, snap, "apps")
			SnapUnset(t, snap, "app-options")
			SnapStop(t, service)
		})

		const newPort = "22222"

		// enable new apps option to aviod mixed options issue with old env option
		SnapSet(t, snap, "app-options", "true")

		// make sure the port is available before using it
		RequirePortAvailable(t, newPort)

		// set apps. and validate the new port comes online
		SnapSet(t, snap, "apps."+app+".config.service-port", newPort)
		SnapStart(t, service)

		WaitServiceOnline(t, serviceWaitTimeout, newPort)

		// unset apps. and validate the default port comes online
		SnapUnset(t, snap, "apps."+app+".config.service-port")
		SnapRestart(t, service)

		WaitServiceOnline(t, serviceWaitTimeout, servicePort)
	})
}

func SetGlobalConfig(t *testing.T, snap, app, servicePort string) {
	service := snap + "." + app

	// start clean
	SnapStop(t, service)

	t.Run("set and unset apps.", func(t *testing.T) {
		t.Cleanup(func() {
			SnapUnset(t, snap, "config")
			SnapUnset(t, snap, "app-options")
			SnapStop(t, service)
		})

		const newPort = "33333"

		// enable new config option to aviod mixed options issue with old env option
		SnapSet(t, snap, "app-options", "true")

		// make sure the port is available before using it
		RequirePortAvailable(t, newPort)

		// set config. and validate the new port comes online
		SnapSet(t, snap, "config.service-port", newPort)
		SnapStart(t, service)

		WaitServiceOnline(t, serviceWaitTimeout, newPort)

		// unset config. and validate the default port comes online
		SnapUnset(t, snap, "config.service-port")
		SnapRestart(t, service)

		WaitServiceOnline(t, serviceWaitTimeout, servicePort)
	})
}

func SetMixedConfig(t *testing.T, snap, app, servicePort string) {
	service := snap + "." + app

	if !FullConfigTest {
		t.Skip("Full config test is disabled.")
	}
	// start clean
	SnapStop(t, service)

	t.Run("use apps. and config. for different values", func(t *testing.T) {
		t.Cleanup(func() {
			SnapUnset(t, snap, "apps")
			SnapUnset(t, snap, "config")
			SnapUnset(t, snap, "app-options")
			SnapStop(t, service)
		})

		const newAppPort = "44444"
		const newConfigPort = "55555"

		// enable new apps/config options to aviod mixed options issue with old env option
		SnapSet(t, snap, "app-options", "true")

		// make sure the ports are available before using it
		RequirePortAvailable(t, newAppPort)
		RequirePortAvailable(t, newConfigPort)

		// set apps. and config. with different values,
		// and validate that app-specific option has been picked up because it has higher precedence
		SnapSet(t, snap, "apps."+app+".config.service-port", newAppPort)
		SnapSet(t, snap, "config.service-port", newConfigPort)
		SnapStart(t, service)

		WaitServiceOnline(t, serviceWaitTimeout, newAppPort)
	})

}
