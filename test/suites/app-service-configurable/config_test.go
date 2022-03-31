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
	utils.Exec(t, "sudo snap stop edgex-app-service-configurable.app-service-configurable")

	t.Run("change service port", func(t *testing.T) {
		t.Cleanup(func() {
			utils.Exec(t, "sudo snap stop edgex-app-service-configurable.app-service-configurable")
			utils.Exec(t, "sudo snap unset edgex-app-service-configurable env.service.port")
		})

		const newPort = "56789"

		// make sure the port is available before using it
		utils.CheckPortAvailable(t, newPort)

		utils.Exec(t, "sudo snap set edgex-app-service-configurable env.service.port="+newPort)
		utils.Exec(t, "sudo snap start edgex-app-service-configurable")
		utils.WaitServiceOnline(t, newPort)
	})

	t.Run("set profile", func(t *testing.T) {
		t.Cleanup(func() {
			utils.Exec(t, "sudo snap stop edgex-app-service-configurable.app-service-configurable")
		})

		time := time.Now()
		profile := "http-export"

		// set profile
		utils.Exec(t,
			"sudo snap set edgex-app-service-configurable profile="+profile,
			"sudo snap start edgex-app-service-configurable",
		)

		//check logs for the record of expected profile
		utils.Exec(t, "sleep 1")
		logs := utils.SnapLogsJournal(t, time, "edgex-app-service-configurable")
		expectLog := "app=app-" + profile
		require.True(t, strings.Contains(logs, expectLog))
	})
}

func TestAppConfig(t *testing.T) {
	t.Skip("TODO")
}
