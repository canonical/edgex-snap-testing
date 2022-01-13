package test

import (
	"log"
	"os"
	"testing"

	"MQTT-test-project/utils"
)

func TestMain(m *testing.M) {

	log.Println("[GLOBAL SETUP]")

	stdout, stderr, err := utils.Command(
		"sudo snap remove --purge edgex-device-mqtt",
		"sudo snap remove --purge edgexfoundry",
		"sudo snap install edgexfoundry --channel=latest/stable",
		"sudo snap install edgex-device-mqtt --channel=latest/stable")
	utils.CommandLog(nil, stdout, stderr, err)

	exitCode := m.Run()

	log.Println("[GLOBAL TEARDOWN]")

	stdout, stderr, err = utils.Command(
		"sudo snap remove --purge edgex-device-mqtt",
		"sudo snap remove --purge edgexfoundry")
	utils.CommandLog(nil, stdout, stderr, err)

	// TODO improvement
	// utils.RemoveSnaps("edgex-device-mqtt", "edgexfoundry")
	// the function can be variadic and take zero or more inputs
	// e.g. https://github.com/canonical/edgex-snap-hooks/blob/50df6237c8eb5b49d497d3a6f978f83391905308/utils.go#L254

	os.Exit(exitCode)
}
