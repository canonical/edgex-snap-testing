package test

import (
	"edgex-snap-testing/test/utils"
	"testing"
)

const defaultServicePort = "59900"

func TestNetworkInterface(t *testing.T) {
	t.Cleanup(func() {
		utils.SnapStop(t, deviceRestService)
	})

	utils.SnapStart(t, deviceRestService)

	t.Run("not listen on all interfaces", func(t *testing.T) {
		utils.RequireListenAllInterfaces(t, false, defaultServicePort)
	})

	t.Run("listen localhost", func(t *testing.T) {
		utils.RequireListenLoopback(t, defaultServicePort)
		utils.RequirePortOpen(t, defaultServicePort)
	})
}
