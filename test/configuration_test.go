package test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"MQTT-test-project/utils"
)

func setupSubtestConfiguration(t *testing.T) {
	t.Log("[SUBTEST SETUP]")
	stdout, stderr, err := utils.Command("sudo snap start --enable edgex-device-mqtt.device-mqtt")
	utils.CommandLog(t, stdout, stderr, err)
}

func TestConfiguration(t *testing.T) {
	setupSubtestConfiguration(t)

	t.Cleanup(func() {
		t.Log("[SUBTEST CLEANUP]")
		stdout, stderr, err := utils.Command("sudo snap stop --disable edgex-device-mqtt.device-mqtt")
		utils.CommandLog(t, stdout, stderr, err)
	})

	t.Run("change-the-maximum-startup-duration", func(t *testing.T) {
		startupDurationValue := "120"

		stdout, stderr, err := utils.Command("sudo snap set edgex-device-mqtt startup-duration=" + startupDurationValue)
		utils.CommandLog(t, stdout, stderr, err)

		stdout, stderr, err = utils.Command("sudo snap get edgex-device-mqtt startup-duration")
		utils.CommandLog(t, stdout, stderr, err)
		assert.Equal(t, startupDurationValue+"\n", stdout, "maximum startup-duration does not set successfully")

	})

	t.Run("change-the-maximum-interval-between-retries", func(t *testing.T) {
		t.Log("Test change the interval between retries")

		startupIntervalValue := "3"

		stdout, stderr, err := utils.Command("sudo snap set edgex-device-mqtt startup-interval=" + startupIntervalValue)
		utils.CommandLog(t, stdout, stderr, err)

		stdout, stderr, err = utils.Command("sudo snap get edgex-device-mqtt startup-interval")
		utils.CommandLog(t, stdout, stderr, err)
		assert.Equal(t, startupIntervalValue+"\n", stdout, "maximum startup-interval does not set successfully")

	})
}
