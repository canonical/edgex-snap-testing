package test

import (
	"edgex-snap-testing/test/utils"
	"testing"
)

// Deprecated
func TestEnvConfig(t *testing.T) {
	utils.SetEnvConfig(t, deviceRfidLlrpSnap, deviceRfidApp, defaultServicePort)
}

func TestAppConfig(t *testing.T) {
	utils.SetAppConfig(t, deviceRfidLlrpSnap, deviceRfidApp, defaultServicePort)
}

func TestGlobalConfig(t *testing.T) {
	utils.SetGlobalConfig(t, deviceRfidLlrpSnap, deviceRfidApp, defaultServicePort)
}

func TestMixedConfig(t *testing.T) {
	utils.SetMixedConfig(t, deviceRfidLlrpSnap, deviceRfidApp, defaultServicePort)
}
