package test

import (
	"edgex-snap-testing/test/utils"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// Deprecated
func TestEnvConfig(t *testing.T) {
	// start clean
	utils.SnapStop(t, ascService)

	t.Run("change service port", func(t *testing.T) {
		t.Cleanup(func() {
			utils.SnapStop(t, ascService)
			utils.SnapUnset(t, ascSnap, "env.service.port")
		})

		const newPort = "56789"

		// make sure the port is available before using it
		utils.CheckPortAvailable(t, newPort)

		utils.SnapSet(t, ascSnap, "env.service.port", newPort)
		utils.SnapStart(t, ascSnap)
		utils.WaitServiceOnline(t, newPort)
	})

	t.Run("set profile", func(t *testing.T) {
		t.Cleanup(func() {
			utils.SnapStop(t, ascService)
		})

		time := time.Now()
		profile := "http-export"

		// set profile
		utils.SnapSet(t, ascSnap, "profile", profile)
		utils.SnapStart(t, ascSnap)

		//check logs for the record of expected profile
		utils.Exec(t, "sleep 1")
		logs := utils.SnapLogsJournal(t, time, ascSnap)
		expectLog := "app=app-" + profile
		require.True(t, strings.Contains(logs, expectLog))
	})
}

func TestAppConfig(t *testing.T) {
	t.Skip("TODO")
}
