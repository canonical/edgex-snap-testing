package test

import (
	"edgex-snap-testing/test/utils"
	"testing"
)

// Deprecated
func TestEnvConfig(t *testing.T) {
	t.Run("change service port", func(t *testing.T) {
		utils.TestEnvChangeServicePort(t, platformSnap, deviceVirtualApp, deviceVirtualDefaultServicePort)
	})
}

func TestAppConfig(t *testing.T) {
	t.Skip("TODO")
}
