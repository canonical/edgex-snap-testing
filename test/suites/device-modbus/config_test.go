package test

import (
	"edgex-snap-testing/test/utils"
	"testing"
)

// Deprecated
func TestEnvConfig(t *testing.T) {
	utils.SetEnvConfig(t, deviceModbusSnap, deviceModbusService, defaultServicePort)
}

func TestAppConfig(t *testing.T) {
	utils.SetAppConfig(t, deviceModbusSnap, deviceModbusService, deviceModbusApp, defaultServicePort)
}

func TestGlobalConfig(t *testing.T) {
	// start clean
	utils.SetGlobalConfig(t, deviceModbusSnap, deviceModbusService, defaultServicePort)
}

func TestMixedConfig(t *testing.T) {
	utils.SetMixedConfig(t, deviceModbusSnap, deviceModbusService, deviceModbusApp, defaultServicePort)
}
