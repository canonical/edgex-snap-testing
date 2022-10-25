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
	start := time.Now()
	utils.TestDeviceVirtualReading(t)

	t.Run("check rule_log", func(t *testing.T) {
		//check logs for the record of expected log
		time.Sleep(15 * time.Second)
		logs := utils.SnapLogs(t, start, ekuiperSnap)
		expectLog := "sink result for rule rule_log"

		require.True(t, strings.Contains(logs, expectLog))
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
