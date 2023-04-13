package test

import (
	"edgex-snap-testing/test/utils"
	"encoding/json"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type Reading struct {
	TotalCount int `json:"totalCount"`
}

const (
	deviceVirtualSnap = "edgex-device-virtual"
	deviceVirtualApp  = "device-virtual"

	ekuiperSnap       = "edgex-ekuiper"
	ekuiperApp        = "ekuiper"
	ekuiperRestfulApi = "ekuiper/rest-api"
)

func TestRulesEngine(t *testing.T) {
	start := time.Now()

	t.Cleanup(func() {
		log.Println("[TEARDOWN SUBTEST]")
		utils.SnapDumpLogs(t, start, deviceVirtualSnap)
		utils.SnapDumpLogs(t, start, ekuiperSnap)

		log.Println("Removing installed snap:", !utils.SkipTeardownRemoval)
		if !utils.SkipTeardownRemoval {
			utils.SnapRemove(t, deviceVirtualSnap)
			utils.SnapRemove(t, ekuiperSnap)
		}
	})

	log.Println("[CLEAN SUBTEST]")
	utils.SnapRemove(t, deviceVirtualSnap)
	utils.SnapRemove(t, ekuiperSnap)

	utils.SnapInstallFromStore(t, deviceVirtualSnap, utils.ServiceChannel)
	utils.SnapInstallFromStore(t, ekuiperSnap, utils.ServiceChannel)

	// TODO: temporary fix
	err := utils.InjectDevicesAndProfilesDirConfig("device-virtual")
	if err != nil {
		log.Fatalf("Failed to inject devices/profiles dir into config: %s", err)
	}

	// turn security off
	utils.SnapSet(t, deviceVirtualSnap, "config.edgex-security-secret-store", "false")
	utils.SnapSet(t, ekuiperSnap, "config.edgex-security-secret-store", "false")

	utils.SnapSet(t, ekuiperSnap, "config.kuiper.basic.debug", "true")

	// set tests to run without a config provider when testing config options as a temporary solution.
	utils.DoNotUseConfigProviderServiceSnap(t, deviceVirtualSnap, deviceVirtualApp)

	// make sure all services are online before starting the tests
	utils.SnapStart(t,
		ekuiperSnap,
		deviceVirtualSnap)

	utils.WaitServiceOnline(t, 60,
		utils.ServicePort(deviceVirtualApp),
		utils.ServicePort(ekuiperApp),
		utils.ServicePort(ekuiperRestfulApi),
	)

	// wait device-virtual to produce readings
	utils.WaitForReadings(t, "Random-Integer-Device", false)

	t.Run("create stream and rule", func(t *testing.T) {
		utils.Exec(t, `edgex-ekuiper.kuiper create stream stream1 '()WITH(FORMAT="JSON",TYPE="edgex")'`)

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

		// wait readings come from ekuiper to edgex message bus
		utils.WaitForReadings(t, "device-test", false)

		req, err := http.NewRequest(http.MethodGet, "http://localhost:59880/api/v2/reading/device/name/device-test", nil)
		require.NoError(t, err)

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		var reading Reading
		if err = json.NewDecoder(resp.Body).Decode(&reading); err != nil {
			t.Fatal(err)
		}

		require.Greaterf(t, reading.TotalCount, 0, "No readings have been re-published to EdgeX message bus by ekuiper")
	})

}
