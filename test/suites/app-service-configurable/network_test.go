package test

import (
	"edgex-snap-testing/test/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNetworkInterface(t *testing.T) {
	t.Cleanup(func() {
		utils.SnapStop(t, ascService)
	})

	utils.SnapStart(t, ascService)

	t.Run("listen default port "+appRulesEngineServicePort, func(t *testing.T) {
		utils.WaitServiceOnline(t, appRulesEngineServicePort)
	})

	t.Run("not listen on all interfaces", func(t *testing.T) {
		isConnected := utils.PortConnectionAllInterface(t, appRulesEngineServicePort)
		require.False(t, isConnected)
	})

	t.Run("listen localhost", func(t *testing.T) {
		isConnected := utils.PortConnectionLocalhost(t, appRulesEngineServicePort)
		require.True(t, isConnected)
	})
}
