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
		utils.SnapStop(t, ekuiperService)
	})

	utils.SnapStart(t, ekuiperService)

	t.Run("listen default port "+serverPort, func(t *testing.T) {
		utils.WaitServiceOnline(t, 60, serverPort)
	})

	t.Run("listen default restful api port "+restfulApiPort, func(t *testing.T) {
		utils.WaitServiceOnline(t, 60, restfulApiPort)
	})
}
