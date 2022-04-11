package test

import (
	"edgex-snap-testing/test/utils"
	"testing"
)

const defaultServicePort = "59993"

func TestNetworkInterface(t *testing.T) {
	t.Cleanup(func() {
		utils.SnapStop(t, deviceSnmpService)
	})

	utils.SnapStart(t, deviceSnmpService)

	t.Run("listen default port "+defaultServicePort, func(t *testing.T) {
		utils.WaitServiceOnline(t, 60, defaultServicePort)
	})

	t.Run("not listen on all interfaces", func(t *testing.T) {
		utils.RequireListenAllInterfaces(t, false, defaultServicePort)
	})

	t.Run("listen localhost", func(t *testing.T) {
		utils.RequireListenLoopback(t, defaultServicePort)
		utils.RequirePortOpen(t, defaultServicePort)
	})
}
