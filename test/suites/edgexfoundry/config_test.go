package test

import (
	"edgex-snap-testing/test/utils"
	"testing"
	"time"
)

const supportSchedulerStartupMsg = "This is the Support Scheduler Microservice"

func TestChangeStartupMsg(t *testing.T) {
	revertCP := utils.DoNotUseConfigProviderPlatformSnap(t, platformSnap, supportSchedulerApp)

	t.Cleanup(func() {
		revertCP()
		utils.SnapRestart(t, supportSchedulerService)
	})

	testChangeStartupMsg_app(t)
	testChangeStartupMsg_global(t)
	testChangeStartupMsg_mixedGlobalApp(t)
}

func testChangeStartupMsg_app(t *testing.T) {
	t.Run("app", func(t *testing.T) {
		const (
			newStartupMsg = "snap-testing (app)"
			startupMsgKey = "apps.support-scheduler.config.service-startupmsg"
		)

		t.Log("Set and verify new startup message:", newStartupMsg)
		utils.SnapSet(t, platformSnap, startupMsgKey, newStartupMsg)
		ts := time.Now()
		utils.SnapRestart(t, supportSchedulerService)

		utils.WaitForLogMessage(t, supportSchedulerService, `msg="`+newStartupMsg+`"`, ts)

		t.Log("Unset and check default message")
		utils.SnapUnset(t, platformSnap, startupMsgKey)
		ts = time.Now()
		utils.SnapRestart(t, supportSchedulerService)
		utils.WaitForLogMessage(t, supportSchedulerService, `msg="`+supportSchedulerStartupMsg+`"`, ts)
	})
}

func testChangeStartupMsg_global(t *testing.T) {
	t.Run("global", func(t *testing.T) {
		const (
			newStartupMsg = "snap-testing (global)"
			startupMsgKey = "config.service-startupmsg"
		)

		t.Log("Set and verify new startup message:", newStartupMsg)
		utils.SnapSet(t, platformSnap, startupMsgKey, newStartupMsg)
		ts := time.Now()
		utils.SnapRestart(t, supportSchedulerService)

		utils.WaitForLogMessage(t, supportSchedulerService, `msg="`+newStartupMsg+`"`, ts)

		t.Log("Unset and check default message")
		utils.SnapUnset(t, platformSnap, startupMsgKey)
		ts = time.Now()
		utils.SnapRestart(t, supportSchedulerService)
		utils.WaitForLogMessage(t, supportSchedulerService, `msg="`+supportSchedulerStartupMsg+`"`, ts)
	})
}

func testChangeStartupMsg_mixedGlobalApp(t *testing.T) {
	t.Run("mixedGlobalApp", func(t *testing.T) {
		const (
			appNewStartupMsg = "snap-testing (app specific)"
			appStartupMsgKey = "apps." + supportSchedulerApp + ".config.service-startupmsg"

			globalNewStartupMsg = "snap-testing (global override)"
			globalStartupMsgKey = "config.service-startupmsg"
		)

		t.Log("Set local and global startup messages and verify that local has taken precedence")
		utils.SnapSet(t, platformSnap, appStartupMsgKey, appNewStartupMsg)
		utils.SnapSet(t, platformSnap, globalStartupMsgKey, globalNewStartupMsg)
		ts := time.Now()
		utils.SnapRestart(t, supportSchedulerService)
		utils.WaitForLogMessage(t, supportSchedulerService, `msg="`+appNewStartupMsg+`"`, ts)

		t.Log("Unset and check default message")
		utils.SnapUnset(t, platformSnap, appStartupMsgKey)
		utils.SnapUnset(t, platformSnap, globalStartupMsgKey)
		ts = time.Now()
		utils.SnapRestart(t, supportSchedulerService)
		utils.WaitForLogMessage(t, supportSchedulerService, `msg="`+supportSchedulerStartupMsg+`"`, ts)
	})
}
