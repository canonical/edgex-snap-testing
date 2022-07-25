package test

import (
	"edgex-snap-testing/test/utils"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
	"time"
)

func TestChangePort_legacyEnv(t *testing.T) {
	if !utils.FullConfigTest {
		t.Skip("Full config test is disabled.")
	}

	const newPort = "11111"
	const envServicePort = "env." + supportSchedulerApp + ".service.port"

	// start clean
	utils.SnapStop(t, supportSchedulerService)

	t.Cleanup(func() {
		utils.SnapUnset(t, platformSnap, envServicePort)
		utils.SnapStop(t, supportSchedulerService)
	})

	// make sure the port is available before using it
	utils.RequirePortAvailable(t, newPort)

	// set env. and validate the new port comes online
	utils.SnapSet(t, platformSnap, envServicePort, newPort)
	utils.SnapStart(t, supportSchedulerService)
	utils.WaitServiceOnline(t, 60, newPort)

	// unset env. and validate the default port comes online
	utils.SnapUnset(t, platformSnap, envServicePort)
	utils.SnapRestart(t, supportSchedulerService)
	utils.WaitServiceOnline(t, 60, supportSchedulerServicePort)
}

func TestChangeStartupMsg_global(t *testing.T) {
	const (
		newStartupMsg    = "snap-testing-startup-message"
		configStartupMsg = "config.service-startupmsg"
	)

	// start clean
	utils.SnapStop(t, supportSchedulerService)

	t.Cleanup(func() {
		utils.SnapUnset(t, platformSnap, configStartupMsg)
		utils.SnapStop(t, supportSchedulerService)
	})

	// enable new config option to avoid mixed options issue with old env option
	utils.SnapSet(t, platformSnap, "app-options", "true")

	// make sure the startupmsg is not testing message before setting it
	utils.SnapStart(t, supportSchedulerService)
	require.False(t, checkStartupMsg(t, supportSchedulerService, newStartupMsg))

	// set config. and validate the startupmsg is testing message
	utils.SnapSet(t, platformSnap, configStartupMsg, newStartupMsg)
	utils.SnapRestart(t, supportSchedulerService)
	require.True(t, checkStartupMsg(t, supportSchedulerService, newStartupMsg))

	// unset config. and validate the startupmsg is not testing message anymore
	utils.SnapUnset(t, platformSnap, configStartupMsg)
	utils.SnapRestart(t, supportSchedulerService)
	require.False(t, checkStartupMsg(t, supportSchedulerService, newStartupMsg))
}

func TestChangeStartupMsg_mixedGlobalApp(t *testing.T) {
	const (
		defaultStartupMsg = "This is the Support Scheduler Microservice"
		appNewStartupMsg  = "snap testing startup message (set by app option)"
		appStartupMsg     = "apps." + supportSchedulerApp + ".config.service-startupmsg"

		globalNewStartupMsg = "snap testing startup message (set by config option)"
		globalStartupMsg    = "config.service-startupmsg"
	)

	// start clean
	utils.SnapStop(t, supportSchedulerService)

	t.Cleanup(func() {
		utils.SnapUnset(t, platformSnap, globalStartupMsg)
		utils.SnapStop(t, supportSchedulerService)
	})

	// enable new config option to avoid mixed options issue with old env option
	utils.SnapSet(t, platformSnap, "app-options", "true")

	// make sure the startupmsg is the default testing message before setting it
	utils.SnapStart(t, supportSchedulerService)
	require.True(t, checkStartupMsg(t, supportSchedulerService, defaultStartupMsg))

	// set apps. and config. with different testing message,
	// and validate that app-specific option has been picked up because it has higher precedence
	utils.SnapSet(t, platformSnap, appStartupMsg, appNewStartupMsg)
	utils.SnapSet(t, platformSnap, globalStartupMsg, globalNewStartupMsg)
	utils.SnapRestart(t, supportSchedulerService)
	require.True(t, checkStartupMsg(t, supportSchedulerService, appNewStartupMsg))

	// unset apps. and config. and validate the startupmsg is back to default
	utils.SnapUnset(t, platformSnap, appStartupMsg)
	utils.SnapUnset(t, platformSnap, globalStartupMsg)
	utils.SnapRestart(t, supportSchedulerService)
	require.True(t, checkStartupMsg(t, supportSchedulerService, defaultStartupMsg))
}

func checkStartupMsg(t *testing.T, snap, expectedMsg string) bool {
	const maxRetry = 10
	var start = time.Now()

	utils.WaitPlatformOnline(t)

	for i := 1; i <= maxRetry; i++ {
		time.Sleep(1 * time.Second)
		t.Logf("Waiting for startup message. Retry %d/%d", i, maxRetry)

		logs := utils.SnapLogs(t, start, snap)
		if strings.Contains(logs, expectedMsg) {
			return true
		}
	}
	t.Logf("Time out: reached max %d retries.", maxRetry)
	return false
}
