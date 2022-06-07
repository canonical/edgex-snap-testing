package utils

import (
	"testing"
)

func TestCommon(t *testing.T, params *TestParams) {
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
	if params.TestOpenPorts != nil {
		WaitServiceOnline(t, 60, params.DefaultServicePort)
	}
	if params.TestBindAddrLoopback {
		RequireListenAllInterfaces(t, false, params.DefaultServicePort)
		RequireListenLoopback(t, params.DefaultServicePort)
		RequirePortOpen(t, params.DefaultServicePort)
	}

	if params.TestSemanticSnapVersion == true {
		RequireSnapSemver(nil, params.Snap)
	}
}
