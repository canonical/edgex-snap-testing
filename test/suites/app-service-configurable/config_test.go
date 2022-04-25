package test

import (
	"edgex-snap-testing/test/utils"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const profile = "http-export"

// Deprecated
func TestEnvConfig(t *testing.T) {
	utils.SetEnvConfig(t, ascSnap, ascService, appSrviceRulesServicePort)
}

func TestAppConfig(t *testing.T) {
	utils.SetAppConfig(t, ascSnap, ascService, appName, appSrviceRulesServicePort)
}

func TestGlobalConfig(t *testing.T) {
	// start clean
	utils.SetGlobalConfig(t, ascSnap, ascService, appSrviceRulesServicePort)
}

func TestMixedConfig(t *testing.T) {
	utils.SetMixedConfig(t, ascSnap, ascService, appName, appSrviceRulesServicePort)
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
