package test

import (
	"edgex-snap-testing/test/utils"
	"testing"
)

var FullConfigTest = true

// Deprecated
func TestEnvConfig(t *testing.T) {
	utils.SetEnvConfig(t, deviceMqttSnap, deviceMqttApp, defaultServicePort)
}

func TestAppConfig(t *testing.T) {
	utils.SetAppConfig(t, deviceMqttSnap, deviceMqttApp, defaultServicePort)
}

func TestGlobalConfig(t *testing.T) {
	// start clean
	utils.SetGlobalConfig(t, deviceMqttSnap, deviceMqttApp, defaultServicePort)
}

func TestMixedConfig(t *testing.T) {
	utils.SetMixedConfig(t, deviceMqttSnap, deviceMqttApp, defaultServicePort)
}
