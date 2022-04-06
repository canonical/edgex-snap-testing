package test

import (
	"edgex-snap-testing/test/utils"
	"testing"
)

func TestNetworkInterface(t *testing.T) {
	t.Cleanup(func() {
		utils.SnapStop(t, snapAppName)
	})

	utils.SnapStart(t, snapAppName)

	t.Run("listen default port "+deviceVirtualDefaultServicePort, func(t *testing.T) {
		utils.WaitServiceOnline(t, deviceVirtualDefaultServicePort)
	})

	t.Run("not listen on all interfaces", func(t *testing.T) {
		utils.RequireListenAllInterfaces(t, false, deviceVirtualDefaultServicePort)
	})

	t.Run("listen localhost", func(t *testing.T) {
		utils.RequireListenLoopback(t, deviceVirtualDefaultServicePort)
		utils.RequirePortOpen(t, deviceVirtualDefaultServicePort)
	})
}
