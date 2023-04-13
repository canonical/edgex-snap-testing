package test

import (
	"edgex-snap-testing/test/utils"
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

	utils.DoNotUseConfigProviderPlatformSnap(t, platformSnap, supportSchedulerApp)

	t.Log("Set and verify new startup message:", newStartupMsg)
	utils.SnapSet(t, platformSnap, startupMsgKey, newStartupMsg)
	ts := time.Now()
	utils.SnapRestart(t, supportSchedulerService)

	require.True(t, utils.WaitForLogMessage(t, supportSchedulerService, `msg="`+newStartupMsg+`"`, ts),
		"new startup message = %s", newStartupMsg)

	t.Log("Unset and check default message")
	utils.SnapUnset(t, platformSnap, startupMsgKey)
	ts = time.Now()
	utils.SnapRestart(t, supportSchedulerService)
	require.True(t, utils.WaitForLogMessage(t, supportSchedulerService, `msg="`+supportSchedulerStartupMsg+`"`, ts),
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

	utils.DoNotUseConfigProviderPlatformSnap(t, platformSnap, supportSchedulerApp)

	t.Log("Set and verify new startup message:", newStartupMsg)
	utils.SnapSet(t, platformSnap, startupMsgKey, newStartupMsg)
	ts := time.Now()
	utils.SnapRestart(t, supportSchedulerService)

	require.True(t, utils.WaitForLogMessage(t, supportSchedulerService, `msg="`+newStartupMsg+`"`, ts),
		"new startup message = %s", newStartupMsg)

	t.Log("Unset and check default message")
	utils.SnapUnset(t, platformSnap, startupMsgKey)
	ts = time.Now()
	utils.SnapRestart(t, supportSchedulerService)
	require.True(t, utils.WaitForLogMessage(t, supportSchedulerService, `msg="`+supportSchedulerStartupMsg+`"`, ts),
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

	utils.DoNotUseConfigProviderPlatformSnap(t, platformSnap, supportSchedulerApp)

	t.Log("Set local and global startup messages and verify that local has taken precedence")
	utils.SnapSet(t, platformSnap, appStartupMsgKey, appNewStartupMsg)
	utils.SnapSet(t, platformSnap, globalStartupMsgKey, globalNewStartupMsg)
	ts := time.Now()
	utils.SnapRestart(t, supportSchedulerService)
	require.True(t,
		utils.WaitForLogMessage(t, supportSchedulerService, `msg="`+appNewStartupMsg+`"`, ts),
		"new startup message = %s", appNewStartupMsg)

	t.Log("Unset and check default message")
	utils.SnapUnset(t, platformSnap, appStartupMsgKey)
	utils.SnapUnset(t, platformSnap, globalStartupMsgKey)
	ts = time.Now()
	utils.SnapRestart(t, supportSchedulerService)
	require.True(t,
		utils.WaitForLogMessage(t, supportSchedulerService, `msg="`+supportSchedulerStartupMsg+`"`, ts),
		"default startup message = %s", supportSchedulerStartupMsg)
}
