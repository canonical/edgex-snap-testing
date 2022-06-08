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
		params.TestOpenPorts = getStrictPorts(params)
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

func getStrictPorts(params *TestParams) []string {
	// check network interface status for all platform ports except for:
	// Kong’s port: 8000
	// Kong-db's port: 5432
	// Redis's port: 6379
	strictPort := func(port string) bool {
		return (port != "8000" && port != "5432" && port != "6379")
	}

	temp := params.TestOpenPorts[:0]

	for _, port := range params.TestOpenPorts {
		if strictPort(port) {
			// copy and increment index
			temp = append(temp, port)
		}
	}
	params.TestOpenPorts = temp
	return params.TestOpenPorts
}
