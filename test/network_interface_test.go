package test

import (
	"MQTT-test-project/utils"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func setupSubtest(t *testing.T) {
	t.Logf("[SETUP]")
	utils.Command("sudo snap install edgexfoundry --channel=latest/beta")
	utils.Command("sudo snap install edgex-device-mqtt --channel=latest/beta")
	utils.Command("sudo snap start --enable edgex-device-mqtt.device-mqtt")
}

func teardownSubtest(t *testing.T) {
	t.Logf("[TEARDOWN]")
	utils.Command("sudo snap remove --purge edgex-device-mqtt")
	utils.Command("sudo snap remove --purge edgexfoundry")
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestNetworkInterface(t *testing.T) {
	cases := []struct {
		name           string
		command        string
		expectedOutput bool
		errorMesssage  string
	}{
		{
			name:           "Test if the service is listening on all the configured network interfaces",
			command:        "sudo lsof -nPi :59982 | grep \\*",
			expectedOutput: false,
			errorMesssage:  "[Error] This service is listening on all the configured network interface which is not allowed.",
		},
		{
			name:           "Test if the service is only bound to the local machine",
			command:        "sudo lsof -nPi :59982 | grep 127.0.0.1",
			expectedOutput: true,
			errorMesssage:  "[Error] This service is not bound to the local machine.",
		},
	}

	setupSubtest(t)
	defer teardownSubtest(t)

	for _, c := range cases {
		t.Logf("[SUBTEST] %s \n", c.name)

		t.Run(c.command, func(t *testing.T) {
			var output bool
			if utils.Command(c.command) != "" {
				output = true
			} else {
				output = false
			}
			assert.Equal(t, c.expectedOutput, output, c.errorMesssage)
		})
	}
}
