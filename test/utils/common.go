package utils

import (
	"testing"
)

func TestCommon(t *testing.T, params Params) {
	// config tests
	t.Run("config", func(t *testing.T) {
		t.Run("change service port", func(t *testing.T) {
			changePort := params.Config.TestChangePort

			// start once so that default configs get uploaded to the registry
			service := params.Snap + "." + changePort.App
			SnapStart(nil, service)
			WaitServiceOnline(nil, 60, changePort.DefaultPort)
			SnapStop(nil, service)

			if changePort.TestLegacyEnvConfig {
				SetEnvConfig(t, params.Snap, changePort.App, changePort.DefaultPort)
			}
			if changePort.TestAppConfig {
				SetAppConfig(t, params.Snap, changePort.App, changePort.DefaultPort)
			}
			if changePort.TestGlobalConfig {
				SetGlobalConfig(t, params.Snap, changePort.App, changePort.DefaultPort)
			}
			if changePort.TestMixedGlobalAppConfig {
				SetMixedConfig(t, params.Snap, changePort.App, changePort.DefaultPort)
			}
		})
	})

	// networking tests
	t.Run("net", func(t *testing.T) {
		if params.Net.StartSnap {
			t.Cleanup(func() {
				SnapStop(t, params.Snap)
			})
			SnapStart(t, params.Snap)
		}

		if len(params.Net.TestOpenPorts) > 0 {
			t.Run("ports open", func(t *testing.T) {
				WaitServiceOnline(t, 60, params.Net.TestOpenPorts...)
			})
		}
		if len(params.Net.TestBindLoopback) > 0 {
			WaitServiceOnline(t, 60, params.Net.TestBindLoopback...)

			t.Run("ports not listening on all interfaces", func(t *testing.T) {
				RequireListenAllInterfaces(t, false, params.Net.TestBindLoopback...)
			})

			t.Run("ports listening on localhost", func(t *testing.T) {
				RequireListenLoopback(t, params.Net.TestBindLoopback...)
				// RequirePortOpen(t, params.TestBindAddrLoopback...)
			})
		}

	})

	// packaging tests
	t.Run("packaging", func(t *testing.T) {
		if params.Packaging.TestSemanticSnapVersion == true {
			RequireSnapSemver(t, params.Snap)
		}
	})
}
