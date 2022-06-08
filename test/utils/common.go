package utils

import (
	"testing"
)

func TestCommon(t *testing.T, params *TestParams) {
	if len(params.DefaultServicePort) > 0 {
		for _, port := range params.DefaultServicePort {
			WaitServiceOnline(t, 60, port)
		}
	}

	if params.TestEnvConfig == true {
		SetEnvConfig(t, params.Snap, params.App, params.DefaultServicePort[0])
	}
	if params.TestAppConfig == true {
		SetAppConfig(t, params.Snap, params.App, params.DefaultServicePort[0])
	}
	if params.TestGlobalConfig == true {
		SetGlobalConfig(t, params.Snap, params.App, params.DefaultServicePort[0])
	}
	if params.TestMixedConfig == true {
		SetMixedConfig(t, params.Snap, params.App, params.DefaultServicePort[0])
	}

	if len(params.TestOpenPorts) > 0 {
		for _, port := range params.TestOpenPorts {
			WaitServiceOnline(t, 60, port)
		}
	}
	if params.TestBindAddrLoopback {
		RequireListenAllInterfaces(t, false, params.DefaultServicePort[0])
		RequireListenLoopback(t, params.DefaultServicePort[0])
		RequirePortOpen(t, params.DefaultServicePort[0])
	}

	if params.TestSemanticSnapVersion == true {
		RequireSnapSemver(nil, params.Snap)
	}
}
