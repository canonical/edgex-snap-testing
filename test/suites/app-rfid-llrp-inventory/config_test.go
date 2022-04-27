package test

import (
	"edgex-snap-testing/test/utils"
	"testing"
)

// Deprecated
func TestEnvConfig(t *testing.T) {
	utils.SetEnvConfig(t, appRfidLlrpSnap, appRfidLlrpApp, defaultServicePort)
}

func TestAppConfig(t *testing.T) {
	utils.SetAppConfig(t, appRfidLlrpSnap, appRfidLlrpApp, defaultServicePort)
}

func TestGlobalConfig(t *testing.T) {
	// start clean
	utils.SetGlobalConfig(t, appRfidLlrpSnap, appRfidLlrpApp, defaultServicePort)
}

func TestMixedConfig(t *testing.T) {
	utils.SetMixedConfig(t, appRfidLlrpSnap, appRfidLlrpApp, defaultServicePort)
}
