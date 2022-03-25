package test

import (
	"edgex-snap-testing/test/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

const defaultServicePort = "59982"

func setupSubtestNetworkInterface(t *testing.T) {
	t.Log("[SUBTEST SETUP]")
	utils.Exec(t, "sudo snap start --enable edgex-device-mqtt.device-mqtt")

	err := utils.WaitServiceOnline(t, defaultServicePort)
	require.NoError(t, err, "Error waiting for services to come online.")
}

func TestNetworkInterface(t *testing.T) {
	setupSubtestNetworkInterface(t)

	t.Cleanup(func() {
		t.Log("[SUBTEST CLEANUP]")
		utils.Exec(t, "sudo snap stop --disable edgex-device-mqtt.device-mqtt")
	})

	t.Run("listen-all-interfaces", func(t *testing.T) {
		t.Log("Test if the service is listening on all the configured network interfaces which is not allowed")

		//stdout, stderr, err := utils.Command("sudo lsof -nPi :59982 | { grep \\* || true; }")
		//utils.CommandLog(t, stdout, stderr, err)
		//require.Empty(t, stdout, "This service is listening on all the configured network interface which is not allowed.")
		isConnected := utils.PortConnectionAllInterface(t, defaultServicePort)
		require.False(t, isConnected, "This service is listening on all the configured network interface which is not allowed.")
	})

	t.Run("listen-localhost-interfaces", func(t *testing.T) {
		t.Log("Test if the service is only bound to the local machine")

		//stdout, stderr, err := utils.Command("sudo lsof -nPi :59982 | { grep 127.0.0.1 || true; }")
		//utils.CommandLog(t, stdout, stderr, err)
		//require.NotEmpty(t, stdout, "This service is not bound to the local machine.")
		isConnected := utils.PortConnectionLocalhost(t, defaultServicePort)
		require.True(t, isConnected, "This service is listening on all the configured network interface which is not allowed.")
	})
}
