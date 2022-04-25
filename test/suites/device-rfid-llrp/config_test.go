package test

import (
	"edgex-snap-testing/test/utils"
	"testing"
)

// Deprecated
func TestEnvConfig(t *testing.T) {
	utils.SetEnvConfig(t, deviceRfidLlrpSnap, deviceRfidLlrpService, defaultServicePort)
}

func TestAppConfig(t *testing.T) {
	utils.SetAppConfig(t, deviceRfidLlrpSnap, deviceRfidLlrpService, deviceRfidApp, defaultServicePort)
}

func TestGlobalConfig(t *testing.T) {
	// start clean
	utils.SetGlobalConfig(t, deviceRfidLlrpSnap, deviceRfidLlrpService, defaultServicePort)
}

func TestMixedConfig(t *testing.T) {
	utils.SetMixedConfig(t, deviceRfidLlrpSnap, deviceRfidLlrpService, deviceRfidApp, defaultServicePort)
}
