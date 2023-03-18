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

func TestChangeStartupMsg_app(t *testing.T) {
	const (
		newStartupMsg = "snap-testing (app)"
		startupMsgKey = "apps.support-scheduler.config.service-startupmsg"
	)

	t.Cleanup(func() {
		utils.SnapUnset(t, platformSnap, startupMsgKey)
		utils.SnapRestart(t, supportSchedulerService)
	})

	t.Log("Set and verify new startup message:", newStartupMsg)
	utils.SnapSet(t, platformSnap, startupMsgKey, newStartupMsg)
	ts := time.Now()
	utils.SnapRestart(t, supportSchedulerService)

	require.True(t, checkStartupMsg(t, supportSchedulerService, newStartupMsg, ts),
		"new startup message = %s", newStartupMsg)

	t.Log("Unset and check default message")
	utils.SnapUnset(t, platformSnap, startupMsgKey)
	ts = time.Now()
	utils.SnapRestart(t, supportSchedulerService)
	require.True(t, checkStartupMsg(t, supportSchedulerService, supportSchedulerStartupMsg, ts),
		"default startup message = %s", supportSchedulerStartupMsg)
}

func TestChangeStartupMsg_global(t *testing.T) {
	const (
		newStartupMsg = "snap-testing (global)"
		startupMsgKey = "config.service-startupmsg"
	)

	t.Cleanup(func() {
		utils.SnapUnset(t, platformSnap, startupMsgKey)
		utils.SnapRestart(t, supportSchedulerService)
	})

	t.Log("Set and verify new startup message:", newStartupMsg)
	utils.SnapSet(t, platformSnap, startupMsgKey, newStartupMsg)
	ts := time.Now()
	utils.SnapRestart(t, supportSchedulerService)

	require.True(t, checkStartupMsg(t, supportSchedulerService, newStartupMsg, ts),
		"new startup message = %s", newStartupMsg)

	t.Log("Unset and check default message")
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

		globalNewStartupMsg = "snap-testing (global override)"
		globalStartupMsgKey = "config.service-startupmsg"
	)

	t.Cleanup(func() {
		utils.SnapUnset(t, platformSnap, globalStartupMsgKey)
		utils.SnapRestart(t, supportSchedulerService)
	})

	t.Log("Set local and global startup messages and verify that local has taken precedence")
	utils.SnapSet(t, platformSnap, appStartupMsgKey, appNewStartupMsg)
	utils.SnapSet(t, platformSnap, globalStartupMsgKey, globalNewStartupMsg)
	ts := time.Now()
	utils.SnapRestart(t, supportSchedulerService)
	require.True(t,
		checkStartupMsg(t, supportSchedulerService, appNewStartupMsg, ts),
		"new startup message = %s", appNewStartupMsg)

	t.Log("Unset and check default message")
	utils.SnapUnset(t, platformSnap, appStartupMsgKey)
	utils.SnapUnset(t, platformSnap, globalStartupMsgKey)
	ts = time.Now()
	utils.SnapRestart(t, supportSchedulerService)
	require.True(t,
		checkStartupMsg(t, supportSchedulerService, supportSchedulerStartupMsg, ts),
		"default startup message = %s", supportSchedulerStartupMsg)
}

func checkStartupMsg(t *testing.T, snap, expectedMsg string, since time.Time) bool {
	t.Skip("Skip while working on a fix: https://github.com/canonical/edgex-snap-testing/issues/172")
	const maxRetry = 10

	utils.WaitPlatformOnline(t)

	for i := 1; i <= maxRetry; i++ {
		time.Sleep(1 * time.Second)
		t.Logf("Retry %d/%d: Waiting for startup message: %s", i, maxRetry, expectedMsg)

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
