package test

import (
	"edgex-snap-testing/test/utils"
	"testing"
)

// Deprecated
func TestEnvConfig(t *testing.T) {
	utils.SetEnvConfig(t, deviceSnmpSnap, deviceSnmpApp, defaultServicePort)
}

func TestAppConfig(t *testing.T) {
	utils.SetAppConfig(t, deviceSnmpSnap, deviceSnmpApp, defaultServicePort)
}

func TestGlobalConfig(t *testing.T) {
	// start clean
	utils.SetGlobalConfig(t, deviceSnmpSnap, deviceSnmpApp, defaultServicePort)
}

func TestMixedConfig(t *testing.T) {
	utils.SetMixedConfig(t, deviceSnmpSnap, deviceSnmpApp, defaultServicePort)
}
