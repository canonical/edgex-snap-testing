package test

import (
	"edgex-snap-testing/test/utils"
	"testing"
)

// Deprecated
func TestEnvConfig(t *testing.T) {
	utils.FullConfigTest = true
	utils.SetEnvConfig(t, deviceMqttSnap, deviceMqttService, defaultServicePort)
}

func TestAppConfig(t *testing.T) {
	utils.SetAppConfig(t, deviceMqttSnap, deviceMqttService, appName, defaultServicePort)
}

func TestGlobalConfig(t *testing.T) {
	// start clean
	utils.SetGlobalConfig(t, deviceMqttSnap, deviceMqttService, defaultServicePort)
}

func TestMixedConfig(t *testing.T) {
	utils.FullConfigTest = true
	utils.SetMixedConfig(t, deviceMqttSnap, deviceMqttService, appName, defaultServicePort)
}
