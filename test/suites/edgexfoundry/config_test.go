package test

import (
	"edgex-snap-testing/test/utils"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const supportSchedulerStartupMsg = "This is the Support Scheduler Microservice"

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
		newStartupMsg = "snap-testing (global)"
		startupMsgKey = "config.service-startupmsg"
	)

	// start clean
	utils.SnapStop(t, supportSchedulerService)

	t.Cleanup(func() {
		utils.SnapUnset(t, platformSnap, startupMsgKey)
		utils.SnapStop(t, supportSchedulerService)
	})

	utils.SnapSet(t, platformSnap, "app-options", "true")

	ts := time.Now()
	utils.SnapStart(t, supportSchedulerService)
	require.True(t, checkStartupMsg(t, supportSchedulerService, supportSchedulerStartupMsg, ts),
		"default startup message = %s", supportSchedulerStartupMsg)

	utils.SnapSet(t, platformSnap, startupMsgKey, newStartupMsg)
	ts = time.Now()
	utils.SnapRestart(t, supportSchedulerService)
	require.True(t, checkStartupMsg(t, supportSchedulerService, newStartupMsg, ts),
		"new startup message = %s", newStartupMsg)

	// unset config and re-check
	utils.SnapUnset(t, platformSnap, startupMsgKey)
	ts = time.Now()
	utils.SnapRestart(t, supportSchedulerService)
	require.True(t, checkStartupMsg(t, supportSchedulerService, supportSchedulerStartupMsg, ts),
		"default startup message = %s", supportSchedulerStartupMsg)
}

func TestChangeStartupMsg_mixedGlobalApp(t *testing.T) {
	const (
		appNewStartupMsg = "snap-testing (app specific)"
		appStartupMsgKey = "apps." + supportSchedulerApp + ".config.service-startupmsg"

		globalNewStartupMsg = "snap-testing (global)"
		globalStartupMsgKey = "config.service-startupmsg"
	)

	// start clean
	utils.SnapStop(t, supportSchedulerService)

	t.Cleanup(func() {
		utils.SnapUnset(t, platformSnap, globalStartupMsgKey)
		utils.SnapStop(t, supportSchedulerService)
	})

	utils.SnapSet(t, platformSnap, "app-options", "true")

	utils.SnapStart(t, supportSchedulerService)
	ts := time.Now()
	require.True(t,
		checkStartupMsg(t, supportSchedulerService, supportSchedulerStartupMsg, ts),
		"default startup message = %s", supportSchedulerStartupMsg)

	// set apps. and config. with different testing message,
	// and validate that app-specific option has been picked up because it has higher precedence
	utils.SnapSet(t, platformSnap, appStartupMsgKey, appNewStartupMsg)
	utils.SnapSet(t, platformSnap, globalStartupMsgKey, globalNewStartupMsg)
	ts = time.Now()
	utils.SnapRestart(t, supportSchedulerService)
	require.True(t,
		checkStartupMsg(t, supportSchedulerService, appNewStartupMsg, ts),
		"new startup message = %s", appNewStartupMsg)

	// unset config and re-check
	utils.SnapUnset(t, platformSnap, appStartupMsgKey)
	utils.SnapUnset(t, platformSnap, globalStartupMsgKey)
	ts = time.Now()
	utils.SnapRestart(t, supportSchedulerService)
	require.True(t,
		checkStartupMsg(t, supportSchedulerService, supportSchedulerStartupMsg, ts),
		"default startup message = %s", supportSchedulerStartupMsg)
}

func checkStartupMsg(t *testing.T, snap, expectedMsg string, since time.Time) bool {
	const maxRetry = 10

	utils.WaitPlatformOnline(t)

	for i := 1; i <= maxRetry; i++ {
		time.Sleep(1 * time.Second)
		t.Logf("Waiting for startup message. Retry %d/%d", i, maxRetry)

		logs := utils.SnapLogs(t, since, snap)
		if strings.Contains(logs, fmt.Sprintf("msg=%s", expectedMsg)) ||
			strings.Contains(logs, fmt.Sprintf(`msg="%s"`, expectedMsg)) {
			t.Logf("Found startup message: %s", expectedMsg)
			return true
		}
	}
	t.Logf("Time out: reached max %d retries.", maxRetry)
	return false
}
