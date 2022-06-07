package test

import (
	"edgex-snap-testing/test/utils"
	"testing"
)

func TestNetworkInterface(t *testing.T) {
	t.Cleanup(func() {
		utils.SnapStop(t, coreMetadataService)
	})

	utils.SnapStart(t, coreMetadataService)
	// check network interface status for all platform ports except for:
	// Kong’s port: 8000
	// Kong-db's port: 5432
	// Redis's port: 6379
	for _, port := range utils.PlatformPorts {
		if port != "8000" && port != "5432" && port != "6379" {
			t.Run("platform port "+port+" not listen on all interfaces", func(t *testing.T) {
				utils.RequireListenAllInterfaces(t, false, port)
			})

			t.Run("platform port "+port+" listen localhost", func(t *testing.T) {
				utils.RequireListenLoopback(t, port)
				utils.RequirePortOpen(t, port)
			})
		}
	}
}
