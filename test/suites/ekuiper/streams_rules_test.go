package test

import (
	"edgex-snap-testing/test/utils"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type Reading struct {
	TotalCount int `json:"totalCount"`
}

type RuleStatus struct {
	LogCount int `json:"sink_log_0_0_records_out_total"`
}

func TestStreamsAndRules(t *testing.T) {
	t.Cleanup(func() {
		utils.SnapStop(t,
			ekuiperSnap,
			deviceVirtualSnap)
	})

	utils.SnapStart(t,
		ekuiperSnap,
		deviceVirtualSnap)

	t.Run("create stream", func(t *testing.T) {
		utils.Exec(t, `edgex-ekuiper.kuiper create stream stream1 '()WITH(FORMAT="JSON",TYPE="edgex")'`)
	})

	t.Run("create rule_log", func(t *testing.T) {
		utils.Exec(t,
			`edgex-ekuiper.kuiper create rule rule_log '
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
			`edgex-ekuiper.kuiper create rule rule_edgex_message_bus '
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
	if err := utils.WaitServiceOnline(t, 60, deviceVirtualPort); err != nil {
		t.Fatal(err)
	}
	utils.WaitForReadings(t, true)

	t.Run("check rule_log", func(t *testing.T) {
		var ruleStatus RuleStatus

		// waiting for readings to come from edgex to ekuiper
		for i := 1; ; i++ {
			time.Sleep(1 * time.Second)
			resp, err := http.Get("http://localhost:59720/rules/rule_log/status")
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if err := json.NewDecoder(resp.Body).Decode(&ruleStatus); err != nil {
				t.Fatal(err)
			}

			t.Logf("Waiting for readings to come from edgex to ekuiper, current retry count: %d/60", i)

			if i <= 60 && ruleStatus.LogCount > 0 {
				t.Logf("Readings are coming to ekuiper now")
				break
			}

			if i > 60 && ruleStatus.LogCount <= 0 {
				t.Logf("Waiting for readings to come from edgex to ekuiper, reached maximum retry count of 60")
				break
			}
		}

		require.Greaterf(t, ruleStatus.LogCount, 0, "No readings have been published to log by ekuiper")
	})

	t.Run("check rule_edgex_message_bus", func(t *testing.T) {
		// wait device-virtual to come online and produce readings
		if err := utils.WaitServiceOnline(t, 60, deviceVirtualPort); err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodGet, "http://localhost:59880/api/v2/reading/device/name/device-test", nil)
		require.NoError(t, err)

		idToken := utils.LoginTestUser(t)
		req.Header.Set("Authorization", "Bearer "+idToken)

		var reading Reading
		client := &http.Client{}
		resp, err := client.Do(req)

		require.NoError(t, err)
		defer resp.Body.Close()

		if err = json.NewDecoder(resp.Body).Decode(&reading); err != nil {
			t.Fatal(err)
		}

		require.Greaterf(t, reading.TotalCount, 0, "No readings have been re-published to EdgeX message bus by ekuiper")
	})
}
