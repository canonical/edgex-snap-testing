package test

import (
	"edgex-snap-testing/test/utils"
	"testing"
)

const defaultServicePort = "4000"

func TestNetworkInterface(t *testing.T) {
	t.Cleanup(func() {
		utils.SnapStop(t, uiService)
	})

	utils.SnapStart(t, uiService)

	t.Run("listen default port "+defaultServicePort, func(t *testing.T) {
		utils.WaitServiceOnline(t, 60, defaultServicePort)
	})

	t.Run("not listen on all interfaces", func(t *testing.T) {
		utils.RequireListenAllInterfaces(t, true, defaultServicePort)
	})
}
