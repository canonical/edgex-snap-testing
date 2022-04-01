package test

import (
	"edgex-snap-testing/test/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

const defaultServicePort = "59910"

func TestNetworkInterface(t *testing.T) {
	t.Cleanup(func() {
		utils.SnapStop(t, deviceGpioService)
	})

	utils.SnapStart(t, deviceGpioService)

	t.Run("listen default port "+defaultServicePort, func(t *testing.T) {
		utils.WaitServiceOnline(t, defaultServicePort)
	})

	t.Run("not listen on all interfaces", func(t *testing.T) {
		isConnected := utils.PortConnectionAllInterface(t, defaultServicePort)
		require.False(t, isConnected)
	})

	t.Run("listen localhost", func(t *testing.T) {
		isConnected := utils.PortConnectionLocalhost(t, defaultServicePort)
		require.True(t, isConnected)
	})
}
