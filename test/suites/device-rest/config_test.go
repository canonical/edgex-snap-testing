package test

import (
	"edgex-snap-testing/test/utils"
	"testing"
)

// Deprecated
func TestEnvConfig(t *testing.T) {
	utils.SetEnvConfig(t, deviceRestSnap, deviceRestApp, defaultServicePort)
}

func TestAppConfig(t *testing.T) {
	utils.SetAppConfig(t, deviceRestSnap, deviceRestApp, defaultServicePort)
}

func TestGlobalConfig(t *testing.T) {
	utils.SetGlobalConfig(t, deviceRestSnap, deviceRestApp, defaultServicePort)
}

func TestMixedConfig(t *testing.T) {
	utils.SetMixedConfig(t, deviceRestSnap, deviceRestApp, defaultServicePort)
}
