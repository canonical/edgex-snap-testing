package test

import (
	"edgex-snap-testing/test/utils"
	"testing"
)

const (
	serverPort     = "20498"
	restfulApiPort = "59720"
)

func TestNetworkInterface(t *testing.T) {
	t.Cleanup(func() {
		utils.SnapStop(t, snapService)
	})

	utils.SnapStart(t, snapService)

	t.Run("listen default port "+serverPort, func(t *testing.T) {
		utils.WaitServiceOnline(t, serverPort)
	})

	t.Run("listen default restful api port "+restfulApiPort, func(t *testing.T) {
		utils.WaitServiceOnline(t, restfulApiPort)
	})
}
