package test

import (
	"edgex-snap-testing/test/utils"
	"testing"
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
		utils.RequireListenAllInterfaces(t, false, appRulesEngineServicePort)
	})

	t.Run("listen localhost", func(t *testing.T) {
		utils.RequireListenLoopback(t, appRulesEngineServicePort)
		utils.RequirePortOpen(t, appRulesEngineServicePort)
	})
}
