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

	t.Run("listen default port "+appServiceRulesServicePort, func(t *testing.T) {
		utils.WaitServiceOnline(t, 60, appServiceRulesServicePort)
	})

	t.Run("not listen on all interfaces", func(t *testing.T) {
		utils.RequireListenAllInterfaces(t, false, appServiceRulesServicePort)
	})

	t.Run("listen localhost", func(t *testing.T) {
		utils.RequireListenLoopback(t, appServiceRulesServicePort)
		utils.RequirePortOpen(t, appServiceRulesServicePort)
	})
}
