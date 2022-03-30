package test

import (
	"edgex-snap-testing/test/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
	"time"
)

const deviceVirtualPort = "59900"

func TestStreamsAndRules(t *testing.T) {

	t.Cleanup(func() {
		utils.SnapStop(t,
			snapService,
			"edgexfoundry.app-service-configurable",
			"edgexfoundry.device-virtual")
	})

	utils.SnapStart(t,
		snapService,
		"edgexfoundry.app-service-configurable",
		"edgexfoundry.device-virtual")

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

	utils.WaitServiceOnline(t, deviceVirtualPort)

	// wait device-virtual producing readings with maximum 60 seconds
	for i := 1; ; i++ {
		time.Sleep(1 * time.Second)
		req, err := http.NewRequest("GET", "http://localhost:59880/api/v2/event/count", nil)
		if err != nil {
			fmt.Print(err)
			return
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Print(err)
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Print(err)
			return
		}

		mapContainer := make(map[string]json.RawMessage)
		err = json.Unmarshal(body, &mapContainer)
		if err != nil {
			fmt.Print(err)
			return
		}

		count := mapContainer["Count"]
		countToInt, _ := strconv.Atoi(string(count))

		fmt.Printf("waiting for device-virtual produce readings, current retry count: %d/60\n", i)

		if i <= 60 && countToInt > 0 {
			fmt.Println("device-virtual is producing readings now")
			break
		}

		if i > 60 && countToInt <= 0 {
			fmt.Println("waiting for device-virtual produce readings, reached maximum retry count of 60")
			break
		}
	}

	t.Run("check rule_log", func(t *testing.T) {
		utils.Exec(t, `edgex-ekuiper.kuiper-cli getstatus rule rule_log`)
	})

	t.Run("check rule_edgex_message_bus", func(t *testing.T) {
		utils.Exec(t, `edgex-ekuiper.kuiper-cli getstatus rule rule_edgex_message_bus`)
	})
}
