package test

import (
	"edgex-snap-testing/utils"
	"testing"
)

var deviceVirtualPort = []string{"59900"}

func setupSubtestStreamsAndRules(t *testing.T) {
	t.Log("[SUBTEST SETUP]")
	utils.Exec(t,
		"sudo snap start --enable edgex-ekuiper.kuiper",
		"sudo snap set edgexfoundry app-service-configurable=on",
		"sudo snap set edgexfoundry device-virtual=on")
}

func TestStreamsAndRuels(t *testing.T) {
	setupSubtestStreamsAndRules(t)

	t.Cleanup(func() {
		t.Log("[SUBTEST CLEANUP]")
		utils.Exec(t,
			"sudo snap stop --disable edgex-ekuiper.kuiper",
			"sudo snap set edgexfoundry app-service-configurable=off",
			"sudo snap set edgexfoundry device-virtual=off")
	})

	t.Run("create-stream", func(t *testing.T) {
		t.Log("Test if create stream works")

		utils.Exec(t, `edgex-ekuiper.kuiper-cli create stream stream1 '()WITH(FORMAT="JSON",TYPE="edgex")'`)
	})

	t.Run("create-rule-log", func(t *testing.T) {
		t.Log("Test if create rule_log works")

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

	t.Run("create-rule-edgex-message-bus", func(t *testing.T) {
		t.Log("Test if create rule_mqtt works")

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

	utils.WaitServiceOnline(t, deviceVirtualPort)

	t.Run("check-rule-log", func(t *testing.T) {
		t.Log("Test if rule_log is running without errors")
		// utils.Exec(t, `edgex-ekuiper.kuiper-cli getstatus rule rule_log | jq '.sink_log_0_0_records_out_total'`)
		utils.Exec(t, `edgex-ekuiper.kuiper-cli getstatus rule rule_log`)
	})

	t.Run("check-rule-edgex-message-bus", func(t *testing.T) {
		t.Log("Test if rule_edgex_message_bus is running without errors")
		// utils.Exec(t, `edgex-ekuiper.kuiper-cli getstatus rule rule_edgex_message_bus | jq '.sink_edgex_message_bus_0_0_records_out_total'`)
		utils.Exec(t, `edgex-ekuiper.kuiper-cli getstatus rule rule_edgex_message_bus`)
	})
}
