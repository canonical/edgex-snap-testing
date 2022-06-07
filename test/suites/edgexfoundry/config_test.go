package test

import (
	"edgex-snap-testing/test/utils"
	"testing"
)

var FullConfigTest = true

// Deprecated
func TestEnvConfig(t *testing.T) {
	const newPort = "11111"
	const envServicePort = "env." + coreMetadataApp + ".service.port"

	// start clean
	utils.SnapStop(t, coreMetadataService)

	t.Cleanup(func() {
		utils.SnapUnset(t, platformSnap, envServicePort)
		utils.SnapStop(t, coreMetadataService)
	})
	t.Run("change core-metadata service port", func(t *testing.T) {

		// make sure the port is available before using it
		utils.RequirePortAvailable(t, newPort)

		// set env. and validate the new port comes online
		utils.SnapSet(t, platformSnap, envServicePort, newPort)
		utils.SnapStart(t, coreMetadataService)
		utils.WaitServiceOnline(t, 60, newPort)

		// unset env. and validate the default port comes online
		utils.SnapUnset(t, platformSnap, envServicePort)
		utils.SnapRestart(t, coreMetadataService)
		utils.WaitServiceOnline(t, 60, coreMetadataDefaultServicePort)

	})
}
