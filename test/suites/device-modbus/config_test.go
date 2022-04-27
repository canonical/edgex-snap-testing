package test

import (
	"edgex-snap-testing/test/utils"
	"testing"
)

// Deprecated
func TestEnvConfig(t *testing.T) {
	utils.SetEnvConfig(t, deviceModbusSnap, deviceModbusApp, defaultServicePort)
}

func TestAppConfig(t *testing.T) {
	utils.SetAppConfig(t, deviceModbusSnap, deviceModbusApp, defaultServicePort)
}

func TestGlobalConfig(t *testing.T) {
	utils.SetGlobalConfig(t, deviceModbusSnap, deviceModbusApp, defaultServicePort)
}

func TestMixedConfig(t *testing.T) {
	utils.SetMixedConfig(t, deviceModbusSnap, deviceModbusApp, defaultServicePort)
}
