package test

import (
	"edgex-snap-testing/test/utils"
	"testing"
)

const defaultServicePort = "59982"

func TestNetworkInterface(t *testing.T) {
	t.Cleanup(func() {
		utils.SnapStop(t, deviceMqttService)
	})

	utils.SnapStart(t, deviceMqttService)

	t.Run("listen default port "+defaultServicePort, func(t *testing.T) {
		utils.RequireServiceOnline(t, defaultServicePort)
	})

	t.Run("not listen on all interfaces", func(t *testing.T) {
		utils.RequireNotListenAllInterfaces(t, defaultServicePort)
	})

	t.Run("listen localhost", func(t *testing.T) {
		utils.RequireListenLoopback(t, defaultServicePort)
		utils.RequirePortOpenLocalhost(t, defaultServicePort)
	})
}
