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

	ascSnap = "edgex-app-service-configurable"
	ascApp  = "app-service-configurable"
)

func TestRulesEngine(t *testing.T) {
	start := time.Now()

	t.Cleanup(func() {
		log.Println("[TEARDOWN SUBTEST]")
		utils.SnapDumpLogs(t, start, deviceVirtualSnap)
		utils.SnapDumpLogs(t, start, ekuiperSnap)
		utils.SnapDumpLogs(t, start, ascSnap)

		log.Println("Removing installed snap:", !utils.SkipTeardownRemoval)
		if !utils.SkipTeardownRemoval {
			utils.SnapRemove(t, deviceVirtualSnap)
			utils.SnapRemove(t, ekuiperSnap)
			utils.SnapRemove(t, ascSnap)
		}
	})

	log.Println("[CLEAN SUBTEST]")
	utils.SnapRemove(t, deviceVirtualSnap)
	utils.SnapRemove(t, ekuiperSnap)
	utils.SnapRemove(t, ascSnap)

	utils.SnapInstallFromStore(t, deviceVirtualSnap, utils.ServiceChannel)
	utils.SnapInstallFromStore(t, ekuiperSnap, utils.ServiceChannel)
	utils.SnapInstallFromStore(t, ascSnap, utils.ServiceChannel)

	// turn security off
	utils.SnapSet(t, deviceVirtualSnap, "config.edgex-security-secret-store", "false")
	utils.SnapSet(t, ascSnap, "config.edgex-security-secret-store", "false")
	utils.SnapSet(t, ekuiperSnap, "config.edgex-security-secret-store", "false")

	// use ASC for event filtering
	utils.SnapSet(t, ekuiperSnap, "config.edgex.default.topic", "rules-events")
	utils.SnapSet(t, ekuiperSnap, "config.edgex.default.messagetype", "event")
	utils.SnapSet(t, ascSnap, "profile", "rules-engine")

	// set tests to run without a config provider when testing config options as a temporary solution.
	utils.DisableConfigProviderServiceSnap(t, deviceVirtualSnap, deviceVirtualApp)
	utils.DisableConfigProviderServiceSnap(t, ascSnap, ascApp)

	// make sure all services are online before starting the tests
	utils.SnapStart(t,
		ekuiperSnap,
		deviceVirtualSnap,
		ascSnap)

	utils.WaitServiceOnline(t, 60,
		utils.ServicePort(deviceVirtualApp),
		utils.ServicePort(ekuiperApp),
		utils.ServicePort(ekuiperRestfulApi),
		utils.ServicePort(ascApp),
	)

	// TODO: temporary fix
	err := utils.InjectDevicesAndProfilesDirConfig("device-virtual")
	if err != nil {
		log.Fatalf("Failed to inject devices/profiles dir into config: %s", err)
	}

	// wait device-virtual to produce readings
	utils.WaitForReadings(t, false)

	t.Run("create stream and rule", func(t *testing.T) {
		utils.Exec(t, `edgex-ekuiper.kuiper create stream stream1 '()WITH(FORMAT="JSON",TYPE="edgex")'`)

		utils.Exec(t,
			`edgex-ekuiper.kuiper create rule rule_edgex_message_bus '
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
