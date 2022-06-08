package utils

import (
	"testing"
)

func TestCommon(t *testing.T, params TestParams) {

	if params.TestEnvConfig == true {
		SetEnvConfig(t, params.Snap, params.App, params.DefaultServicePort)
	}
	if params.TestAppConfig == true {
		SetAppConfig(t, params.Snap, params.App, params.DefaultServicePort)
	}
	if params.TestGlobalConfig == true {
		SetGlobalConfig(t, params.Snap, params.App, params.DefaultServicePort)
	}
	if params.TestMixedConfig == true {
		SetMixedConfig(t, params.Snap, params.App, params.DefaultServicePort)
	}

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

	if params.TestSemanticSnapVersion == true {
		RequireSnapSemver(t, params.Snap)
	}
}
