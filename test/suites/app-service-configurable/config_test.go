package test

import (
	"edgex-snap-testing/test/utils"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const profile = "http-export"
const appRulesEngineServicePort = "59720"

// Deprecated
func TestEnvConfig(t *testing.T) {
	// start clean
	utils.SnapStop(t, ascService)

	t.Run("change service port", func(t *testing.T) {
		t.Cleanup(func() {
			utils.SnapStop(t, ascService)
			utils.SnapUnset(t, ascSnap, "env.service.port")
			utils.SnapSet(t, ascSnap, "env.service.port", appRulesEngineServicePort)
		})

		const newPort = "56789"

		// make sure the port is available before using it
		utils.RequirePortAvailable(t, newPort)

		utils.SnapSet(t, ascSnap, "env.service.port", newPort)
		utils.SnapStart(t, ascSnap)
		utils.WaitServiceOnline(t, 60, newPort)
	})
}

func TestAppConfig(t *testing.T) {
	t.Skip("TODO")
}

func TestProfileConfig(t *testing.T) {
	// start clean
	utils.SnapStop(t, ascService)

	t.Run("set profile", func(t *testing.T) {
		t.Cleanup(func() {
			utils.SnapStop(t, ascService)
			utils.SnapUnset(t, ascSnap, "profile")
			// set profile back to default for upcoming tests
			utils.SnapSet(t, ascSnap, "profile", defaultProfile)
		})

		var start = time.Now()

		// set profile
		utils.SnapSet(t, ascSnap, "profile", profile)
		utils.SnapStart(t, ascSnap)

		// check logs for the record of expected profile

		//check logs for the record of expected profile
		time.Sleep(1 * time.Second)
		logs := utils.SnapLogs(t, start, ascSnap)
		expectLog := "app=app-" + profile

		require.True(t, strings.Contains(logs, expectLog))
	})
}
