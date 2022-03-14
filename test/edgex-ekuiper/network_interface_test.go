package test

import (
	"edgex-snap-testing/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

var serverPort = []string{"20498"}
var restfulApiPort = []string{"59720"}

func setupSubtestNetworkInterface(t *testing.T) {
	t.Log("[SUBTEST SETUP]")
	utils.Exec(t, "sudo snap start --enable edgex-ekuiper.kuiper")
}

func TestNetworkInterface(t *testing.T) {
	setupSubtestNetworkInterface(t)

	t.Cleanup(func() {
		t.Log("[SUBTEST CLEANUP]")
		utils.Exec(t, "sudo snap stop --disable edgex-ekuiper.kuiper")
	})

	t.Run("kuiper-server", func(t *testing.T) {
		t.Log("Test if kuiper server is listening on port 20498")

		err := utils.WaitServiceOnline(t, serverPort)
		require.NoError(t, err, "kuiper server is not listening on port 20498.")
	})

	t.Run("restful-api", func(t *testing.T) {
		t.Log("Test if restful api is listening on port 59720")

		err := utils.WaitServiceOnline(t, restfulApiPort)
		require.NoError(t, err, "restful api is not listening on port 59720.")
	})
}
