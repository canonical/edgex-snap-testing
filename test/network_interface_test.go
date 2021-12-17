package test

import (
	"MQTT-test-project/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func setupSubtest(t *testing.T) {
	t.Logf("[SETUP]")
	utils.Command(t, "sudo snap install edgexfoundry --channel=latest/beta")
	utils.Command(t, "sudo snap install edgex-device-mqtt --channel=latest/beta")
	utils.Command(t, "sudo snap start --enable edgex-device-mqtt.device-mqtt")
}

func TestNetworkInterface(t *testing.T) {
	setupSubtest(t)
	t.Cleanup(func() {
		t.Logf("[CLEANUP]")
		utils.Command(t, "sudo snap remove --purge edgex-device-mqtt")
		utils.Command(t, "sudo snap remove --purge edgexfoundry")
		// TODO improvement
		// utils.RemoveSnaps("edgex-device-mqtt", "edgexfoundry")
		// the function can be variadic and take zero or more inputs
		// e.g. https://github.com/canonical/edgex-snap-hooks/blob/50df6237c8eb5b49d497d3a6f978f83391905308/utils.go#L254
	})

	t.Run("listen-all-interfaces", func(t *testing.T) {
		t.Logf("[SUBTEST] Test if the service is listening on all the configured network interfaces")
		t.Cleanup(func() {
			// subtest cleanup
		})
		output := utils.Command(t, "sudo lsof -nPi :59982 | { grep \\* || true; }")
		assert.Equal(t, "", output, "This service is listening on all the configured network interface which is not allowed.")
	})

	t.Run("listen-localhost-interfaces", func(t *testing.T) {
		t.Logf("[SUBTEST] Test if the service is only bound to the local machine")
		t.Cleanup(func() {
			// subtest cleanup
		})
		output := utils.Command(t, "sudo lsof -nPi :59982 | { grep 127.0.0.1 || true; }")
		assert.NotEmpty(t, output, "This service is not bound to the local machine.")
	})
}
