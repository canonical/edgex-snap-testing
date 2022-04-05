package test

import (
	"edgex-snap-testing/test/utils"
	"testing"
)

// Deprecated
func TestEnvConfig(t *testing.T) {
	// start clean
	utils.SnapStop(t, deviceSnmpService)

	t.Run("change service port", func(t *testing.T) {
		t.Cleanup(func() {
			utils.SnapStop(t, deviceSnmpService)
		})

		utils.SnapSetUnset(t, deviceSnmpSnap, defaultServicePort)
	})
}

func TestAppConfig(t *testing.T) {
	t.Skip("TODO")
}
