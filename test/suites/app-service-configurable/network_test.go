package test

import (
	"edgex-snap-testing/test/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

const defaultServicePort = "59701"

func TestNetworkInterface(t *testing.T) {
	t.Cleanup(func() {
		utils.Exec(t, "sudo snap stop edgex-app-service-configurable.app-service-configurable")
	})

	utils.Exec(t, "sudo snap set edgex-app-service-configurable env.service.port="+defaultServicePort)
	utils.Exec(t, "sudo snap start edgex-app-service-configurable.app-service-configurable")

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
