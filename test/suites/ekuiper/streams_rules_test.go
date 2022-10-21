package test

import (
	"edgex-snap-testing/test/utils"
	"testing"
)

func TestStreamsAndRules(t *testing.T) {
	t.Cleanup(func() {
		utils.SnapStop(t,
			ekuiperService,
			deviceVirtualSnap)
	})

	utils.SnapStart(t,
		ekuiperService,
		deviceVirtualSnap)

	t.Run("create stream", func(t *testing.T) {
		utils.Exec(t, `edgex-ekuiper.kuiper-cli create stream stream1 '()WITH(FORMAT="JSON",TYPE="edgex")'`)
	})

	t.Run("create rule_log", func(t *testing.T) {
		utils.Exec(t,
			`edgex-ekuiper.kuiper-cli create rule rule_log '
			{
				"sql":"SELECT * from stream1",
				"actions":[
					{
						"log":{}
					}
				]
			}'`)
	})

	t.Run("create rule_edgex_message_bus", func(t *testing.T) {
		utils.Exec(t,
			`edgex-ekuiper.kuiper-cli create rule rule_edgex_message_bus '
			{
			   "sql":"SELECT * from stream1",
			   "actions": [
				  {
					 "edgex": {
						"connectionSelector": "edgex.redisMsgBus",
						"topicPrefix": "edgex/events/device", 
						"messageType": "request",
						"deviceName": "device-test"
					 }
				  }
			   ]
			}'`)
	})

	// wait device-virtual to come online and produce readings
	utils.WaitServiceOnline(t, 60, deviceVirtualPort)
	utils.TestDeviceVirtualReading(t)

	t.Run("check rule_log", func(t *testing.T) {
		utils.Exec(t, `edgex-ekuiper.kuiper-cli getstatus rule rule_log`)
	})

	t.Run("check rule_edgex_message_bus", func(t *testing.T) {
		utils.Exec(t, `edgex-ekuiper.kuiper-cli getstatus rule rule_edgex_message_bus`)
	})
}
