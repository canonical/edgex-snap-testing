package utils

import (
	"testing"
)

func TestCommon(t *testing.T, params Params) {
	// config tests
	if params.TestEnvConfig {
		SetEnvConfig(t, params.Snap, params.App, params.DefaultServicePort)
	}
	if params.TestAppConfig {
		SetAppConfig(t, params.Snap, params.App, params.DefaultServicePort)
	}
	if params.TestGlobalConfig {
		SetGlobalConfig(t, params.Snap, params.App, params.DefaultServicePort)
	}
	if params.TestMixedConfig {
		SetMixedConfig(t, params.Snap, params.App, params.DefaultServicePort)
	}

	// networking tests
	if len(params.TestOpenPorts) > 0 {
		for _, port := range params.TestOpenPorts {
			WaitServiceOnline(t, 60, port)
		}
	}
	if params.TestBindAddrLoopback {
		for _, port := range params.TestOpenPorts {
			t.Run("platform port "+port+" not listen on all interfaces", func(t *testing.T) {
				RequireListenAllInterfaces(t, false, port)
			})

			t.Run("platform port "+port+" listen localhost", func(t *testing.T) {
				RequireListenLoopback(t, port)
				RequirePortOpen(t, port)
			})
		}
	}

	// packaging tests
	if params.TestSemanticSnapVersion == true {
		RequireSnapSemver(t, params.Snap)
	}
}
