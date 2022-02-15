package test

import (
	"log"
	"os"
	"testing"

	"MQTT-test-project/utils"
)

func TestMain(m *testing.M) {

	log.Println("[GLOBAL SETUP]")

	// TODO
	// utils.InstallSnaps("edgex-device-mqtt", "edgexfoundry")
	// utils.RemoveSnaps("edgex-device-mqtt", "edgexfoundry")
	// the function can be variadic and take zero or more inputs
	// e.g. https://github.com/canonical/edgex-snap-hooks/blob/50df6237c8eb5b49d497d3a6f978f83391905308/utils.go#L254

	utils.RunCommand(nil,
		"sudo snap remove --purge edgex-device-mqtt",
		"sudo snap remove --purge edgexfoundry",
		"sudo snap install edgexfoundry --channel=latest/stable",
		"sudo snap install edgex-device-mqtt --channel=latest/stable")

	exitCode := m.Run()

	log.Println("[GLOBAL TEARDOWN]")

	utils.RunCommand(nil,
		"sudo snap remove --purge edgex-device-mqtt",
		"sudo snap remove --purge edgexfoundry")

	os.Exit(exitCode)
}
