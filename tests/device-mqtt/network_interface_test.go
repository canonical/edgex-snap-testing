package test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"MQTT-test-project/utils"
)

var port = []string{"59982"}

func setupSubtestNetworkInterface(t *testing.T) {
	t.Log("[SUBTEST SETUP]")
	stdout, stderr, err := utils.Command("sudo snap start --enable edgex-device-mqtt.device-mqtt")
	utils.CommandLog(t, stdout, stderr, err)

	err = utils.WaitServiceOnline(t, port)
	utils.CommandLog(t, "", "", err)
}

func TestNetworkInterface(t *testing.T) {
	setupSubtestNetworkInterface(t)

	t.Cleanup(func() {
		t.Log("[SUBTEST CLEANUP]")
		stdout, stderr, err := utils.Command("sudo snap stop --disable edgex-device-mqtt.device-mqtt")
		utils.CommandLog(t, stdout, stderr, err)
	})

	t.Run("listen-all-interfaces", func(t *testing.T) {
		t.Log("Test if the service is listening on all the configured network interfaces which is not allowed")

		//stdout, stderr, err := utils.Command("sudo lsof -nPi :59982 | { grep \\* || true; }")
		//utils.CommandLog(t, stdout, stderr, err)
		//require.Empty(t, stdout, "This service is listening on all the configured network interface which is not allowed.")
		isConnected := utils.PortConnectionAllInterface(t, port)
		require.False(t, isConnected, "This service is listening on all the configured network interface which is not allowed.")
	})

	t.Run("listen-localhost-interfaces", func(t *testing.T) {
		t.Log("Test if the service is only bound to the local machine")

		//stdout, stderr, err := utils.Command("sudo lsof -nPi :59982 | { grep 127.0.0.1 || true; }")
		//utils.CommandLog(t, stdout, stderr, err)
		//require.NotEmpty(t, stdout, "This service is not bound to the local machine.")
		isConnected := utils.PortConnectionLocalhost(t, port)
		require.True(t, isConnected, "This service is listening on all the configured network interface which is not allowed.")
	})
}
