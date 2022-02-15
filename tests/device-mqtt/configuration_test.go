package test

import (
	"edgex-snap-testing/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupSubtestConfiguration(t *testing.T) {
	t.Log("[SUBTEST SETUP]")
	utils.RunCommand(t, "sudo snap start --enable edgex-device-mqtt.device-mqtt")
}

func TestConfiguration(t *testing.T) {
	setupSubtestConfiguration(t)

	t.Cleanup(func() {
		t.Log("[SUBTEST CLEANUP]")
		utils.RunCommand(t, "sudo snap stop --disable edgex-device-mqtt.device-mqtt")
	})

	t.Run("change-the-maximum-startup-duration", func(t *testing.T) {
		startupDurationValue := "120"

		utils.RunCommand(t, "sudo snap set edgex-device-mqtt startup-duration="+startupDurationValue)

		stdout, _ := utils.RunCommand(t, "sudo snap get edgex-device-mqtt startup-duration")
		assert.Equal(t, startupDurationValue+"\n", stdout, "maximum startup-duration does not set successfully")

	})

	t.Run("change-the-maximum-interval-between-retries", func(t *testing.T) {
		t.Log("Test change the interval between retries")

		startupIntervalValue := "3"

		utils.RunCommand(t, "sudo snap set edgex-device-mqtt startup-interval="+startupIntervalValue)

		stdout, _ := utils.RunCommand(t, "sudo snap get edgex-device-mqtt startup-interval")
		assert.Equal(t, startupIntervalValue+"\n", stdout, "maximum startup-interval does not set successfully")

	})
}
