package utils

import "testing"

func SetEnvConfig(t *testing.T, snap, service, defaultServicePort string) {
	if !FullConfigTest {
		// make this subtest optional to save testing time,
		t.Skip("Full config test is disabled by default, and similar full config tests have been operated in device-mqtt test suite.")
	} else {
		// start clean
		SnapStop(t, snap)

		t.Run("change service port", func(t *testing.T) {
			t.Cleanup(func() {
				SnapUnset(t, snap, "env")
				SnapStop(t, snap)
			})

			const newPort = "11111"

			// make sure the port is available before using it
			RequirePortAvailable(t, newPort)

			// set env. and validate the new port comes online
			SnapSet(t, snap, "env.service.port", newPort)
			SnapStart(t, snap)

			WaitServiceOnline(t, 60, newPort)

			// unset env. and validate the default port comes online
			SnapUnset(t, snap, "env.service.port")
			SnapRestart(t, service)
			WaitServiceOnline(t, 60, defaultServicePort)
		})
	}
}

func SetAppConfig(t *testing.T, snap, service, appName, defaultServicePort string) {
	// start clean
	SnapStop(t, snap)

	t.Run("set and unset apps.", func(t *testing.T) {
		t.Cleanup(func() {
			SnapUnset(t, snap, "apps")
			SnapUnset(t, snap, "config-enabled")
			SnapStop(t, snap)
		})

		const newPort = "22222"

		// enable new apps option to aviod mixed options issue with old env option
		SnapSet(t, snap, "config-enabled", "true")

		// make sure the port is available before using it
		RequirePortAvailable(t, newPort)

		// set apps. and validate the new port comes online
		SnapSet(t, snap, "apps."+appName+".config.service-port", newPort)
		SnapStart(t, snap)

		WaitServiceOnline(t, 60, newPort)

		// unset apps. and validate the default port comes online
		SnapUnset(t, snap, "apps."+appName+".config.service-port")
		SnapRestart(t, service)

		WaitServiceOnline(t, 60, defaultServicePort)
	})
}

func SetGlobalConfig(t *testing.T, snap, service, defaultServicePort string) {
	// start clean
	SnapStop(t, snap)

	t.Run("set and unset apps.", func(t *testing.T) {
		t.Cleanup(func() {
			SnapUnset(t, snap, "config")
			SnapUnset(t, snap, "config-enabled")
			SnapStop(t, snap)
		})

		const newPort = "33333"

		// enable new config option to aviod mixed options issue with old env option
		SnapSet(t, snap, "config-enabled", "true")

		// make sure the port is available before using it
		RequirePortAvailable(t, newPort)

		// set config. and validate the new port comes online
		SnapSet(t, snap, "config.service-port", newPort)
		SnapStart(t, snap)

		WaitServiceOnline(t, 60, newPort)

		// unset config. and validate the default port comes online
		SnapUnset(t, snap, "config.service-port")
		SnapRestart(t, service)

		WaitServiceOnline(t, 60, defaultServicePort)
	})
}

func SetMixedConfig(t *testing.T, snap, service, appName, defaultServicePort string) {
	if !FullConfigTest {
		// make this subtest optional to save testing time,
		// similar full config tests have been operated in device-mqtt test suite
		t.Skip("Full config test is disabled by default, and similar full config tests have been operated in device-mqtt test suite.")
	} else {
		// start clean
		SnapStop(t, snap)

		t.Run("use apps. and config. for different values", func(t *testing.T) {
			t.Cleanup(func() {
				SnapUnset(t, snap, "apps")
				SnapUnset(t, snap, "config")
				SnapUnset(t, snap, "config-enabled")
				SnapStop(t, service)
			})

			const newAppPort = "44444"
			const newConfigPort = "55555"

			// enable new apps/config options to aviod mixed options issue with old env option
			SnapSet(t, snap, "config-enabled", "true")

			// make sure the ports are available before using it
			RequirePortAvailable(t, newAppPort)
			RequirePortAvailable(t, newConfigPort)

			// set apps. and config. with different values,
			// and validate that app-specific option has been picked up because it has higher precedence
			SnapSet(t, snap, "apps."+appName+".config.service-port", newAppPort)
			SnapSet(t, snap, "config.service-port", newConfigPort)
			SnapStart(t, snap)

			WaitServiceOnline(t, 60, newAppPort)
		})
	}
}
