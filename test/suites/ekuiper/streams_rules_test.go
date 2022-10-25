package test

import (
	"edgex-snap-testing/test/utils"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"
)

type Reading struct {
	TotalCount int `json:"totalCount"`
}

type RuleStatus struct {
	SouceCount int `json:"source_stream1_0_records_out_total"`
	LogCount   int `json:"sink_log_0_0_records_out_total"`
}

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
				"sql":"SELECT * FROM stream1 WHERE meta(deviceName) != \"device-test\"",
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
			   "sql":"SELECT * FROM stream1 WHERE meta(deviceName) != \"device-test\"",
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
		var ruleStatus RuleStatus
		var body string

		// waiting for readings to come from edgex to ekuiper
		for i := 1; ; i++ {
			time.Sleep(1 * time.Second)
			stdout, _, err := utils.Exec(t, "edgex-ekuiper.kuiper-cli getstatus rule rule_log")
			if err != nil {
				t.Fatal(err)
			}

			startIndex := strings.Index(stdout, "{")

			if startIndex < 0 {
				t.Fatal("Error getting status of rule rule_log in JSON format")
			} else {
				body = stdout[startIndex:]
			}

			err = json.Unmarshal([]byte(body), &ruleStatus)
			if err != nil {
				t.Fatal(err)
			}

			t.Logf("Waiting for readings to come from edgex to ekuiper, current retry count: %d/60", i)

			if i <= 60 && ruleStatus.SouceCount > 0 {
				t.Logf("Readings are coming to ekuiper now")
				break
			}

			if i > 60 && ruleStatus.SouceCount <= 0 {
				t.Logf("Waiting for readings to come from edgex to ekuiper, reached maximum retry count of 60")
				break
			}
		}

		stdout, _, err := utils.Exec(t, "edgex-ekuiper.kuiper-cli getstatus rule rule_log")
		if err != nil {
			t.Fatal(err)
		}

		startIndex := strings.Index(stdout, "{")

		if startIndex < 0 {
			t.Fatal("Error getting status of rule rule_log in JSON format")
		} else {
			body = stdout[startIndex:]
		}

		err = json.Unmarshal([]byte(body), &ruleStatus)
		if err != nil {
			t.Fatal(err)
		}

		require.Greaterf(t, ruleStatus.LogCount, 0, "No readings have been published to log by ekuiper")

	})

	t.Run("check rule_edgex_message_bus", func(t *testing.T) {
		// wait device-virtual to come online and produce readings
		utils.WaitServiceOnline(t, 60, deviceVirtualPort)

		var reading Reading
		resp, err := http.Get("http://localhost:59880/api/v2/reading/device/name/device-test")
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		err = json.Unmarshal(body, &reading)
		if err != nil {
			t.Fatal(err)
		}

		require.Greaterf(t, reading.TotalCount, 0, "No readings have been re-published to EdgeX message bus by ekuiper")
	})
}
