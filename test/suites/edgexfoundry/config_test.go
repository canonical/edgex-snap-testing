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

	const (
		newPort        = "11111"
		envServicePort = "env." + supportSchedulerApp + ".service.port"
	)

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
		defaultStartupMsg = "This is the Support Scheduler Microservice"

		newStartupMsg = "snap-testing-startup-message"
		startupMsgKey = "config.service-startupmsg"
	)

	// start clean
	utils.SnapStop(t, supportSchedulerService)

	t.Cleanup(func() {
		utils.SnapUnset(t, platformSnap, startupMsgKey)
		utils.SnapStop(t, supportSchedulerService)
	})

	// enable new config option to avoid mixed options issue with old env option
	utils.SnapSet(t, platformSnap, "app-options", "true")

	// make sure the startupmsg is the default testing message before setting it
	utils.SnapStart(t, supportSchedulerService)
	require.True(t, checkStartupMsg(t, supportSchedulerService, defaultStartupMsg))

	// set config. and validate the startupmsg is the testing message
	utils.SnapSet(t, platformSnap, startupMsgKey, newStartupMsg)
	utils.SnapRestart(t, supportSchedulerService)
	require.True(t, checkStartupMsg(t, supportSchedulerService, newStartupMsg))

	// unset config. and validate the startupmsg is not the testing message anymore
	utils.SnapUnset(t, platformSnap, startupMsgKey)
	utils.SnapRestart(t, supportSchedulerService)
	require.True(t, checkStartupMsg(t, supportSchedulerService, defaultStartupMsg))
}

func TestChangeStartupMsg_mixedGlobalApp(t *testing.T) {
	const (
		defaultStartupMsg = "This is the Support Scheduler Microservice"

		appNewStartupMsg = "snap testing startup message (set by app option)"
		appStartupMsgKey = "apps." + supportSchedulerApp + ".config.service-startupmsg"

		globalNewStartupMsg = "snap testing startup message (set by config option)"
		globalStartupMsgKey = "config.service-startupmsg"
	)

	// start clean
	utils.SnapStop(t, supportSchedulerService)

	t.Cleanup(func() {
		utils.SnapUnset(t, platformSnap, globalStartupMsgKey)
		utils.SnapStop(t, supportSchedulerService)
	})

	// enable new config option to avoid mixed options issue with old env option
	utils.SnapSet(t, platformSnap, "app-options", "true")

	// make sure the startupmsg is the default testing message before setting it
	utils.SnapStart(t, supportSchedulerService)
	require.True(t, checkStartupMsg(t, supportSchedulerService, defaultStartupMsg))

	// set apps. and config. with different testing message,
	// and validate that app-specific option has been picked up because it has higher precedence
	utils.SnapSet(t, platformSnap, appStartupMsgKey, appNewStartupMsg)
	utils.SnapSet(t, platformSnap, globalStartupMsgKey, globalNewStartupMsg)
	utils.SnapRestart(t, supportSchedulerService)
	require.True(t, checkStartupMsg(t, supportSchedulerService, appNewStartupMsg))

	// unset apps. and config. and validate the startupmsg is back to default
	utils.SnapUnset(t, platformSnap, appStartupMsgKey)
	utils.SnapUnset(t, platformSnap, globalStartupMsgKey)
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
