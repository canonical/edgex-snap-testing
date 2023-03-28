package test

import (
	"edgex-snap-testing/test/utils"
	"encoding/json"
	"io/ioutil"
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
	deviceVirtualPort = "59900"

	ekuiperSnap           = "edgex-ekuiper"
	ekuiperApp            = "ekuiper"
	ekuiperServerPort     = "20498"
	ekuiperRestfulApiPort = "59720"

	ascSnap             = "edgex-app-service-configurable"
	ascApp              = "app-service-configurable"
	ascServiceRulesPort = "59701"
)

func TestRulesEngine(t *testing.T) {
	teardown, err := subtestSetup(t)
	if err != nil {
		log.Fatalf("Failed to setup tests: %s", err)
	}

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

	teardown()
}

func subtestSetup(t *testing.T) (teardown func(), err error) {
	log.Println("[SUBTEST CLEAN]")
	utils.SnapRemove(t, deviceVirtualSnap)
	utils.SnapRemove(t, ekuiperSnap)
	utils.SnapRemove(t, ascSnap)

	log.Println("[SUBTEST SETUP]")
	start := time.Now()

	teardown = func() {
		log.Println("[SUBTEST TEARDOWN]")
		utils.SnapDumpLogs(t, start, deviceVirtualSnap)
		utils.SnapDumpLogs(t, start, ekuiperSnap)
		utils.SnapDumpLogs(t, start, ascSnap)

		log.Println("Removing installed snap:", !utils.SkipTeardownRemoval)
		if !utils.SkipTeardownRemoval {
			utils.SnapRemove(t, deviceVirtualSnap)
			utils.SnapRemove(t, ekuiperSnap)
			utils.SnapRemove(t, ascSnap)
		}
	}

	if err = utils.SnapInstallFromStore(t, deviceVirtualSnap, utils.ServiceChannel); err != nil {
		teardown()
		return
	}

	if err = utils.SnapInstallFromStore(t, ekuiperSnap, utils.ServiceChannel); err != nil {
		teardown()
		return
	}

	if err = utils.SnapInstallFromStore(t, ascSnap, utils.ServiceChannel); err != nil {
		teardown()
		return
	}

	// turn security off
	utils.SnapSet(t, deviceVirtualSnap, "config.edgex-security-secret-store", "false")
	utils.SnapSet(t, ascSnap, "config.edgex-security-secret-store", "false")
	utils.SnapSet(t, ekuiperSnap, "config.edgex-security-secret-store", "false")

	// use ASC for event filtering
	utils.SnapSet(t, ekuiperSnap, "config.edgex.default.topic", "rules-events")
	utils.SnapSet(t, ekuiperSnap, "config.edgex.default.messagetype", "event")
	utils.SnapSet(t, ascSnap, "profile", "rules-engine")

	// set tests to run without a config provider when testing config options as a temporary solution. 
	// update this once the following PR has been merged: https://github.com/canonical/edgex-snap-testing/pull/175
	disableConfigProviderServiceSnap(t, deviceVirtualSnap, deviceVirtualApp)
	disableConfigProviderServiceSnap(t, ascSnap, ascApp)

	// make sure all services are online before starting the tests
	utils.SnapStart(t,
		ekuiperSnap,
		deviceVirtualSnap,
		ascSnap)

	if err = utils.WaitServiceOnline(t, 60,
		deviceVirtualPort,
		ekuiperServerPort,
		ekuiperRestfulApiPort,
		ascServiceRulesPort,
	); err != nil {
		teardown()
		return
	}

	// wait device-virtual to produce readings
	utils.WaitForReadings(t, false)

	return
}

// disableConfigProviderServiceSnap disables the config provider for the specified app,
// copies the common configuration file from the platform snap to the service snap,
// and sets the common configuration path.
func disableConfigProviderServiceSnap(t *testing.T, snap, app string) {
	utils.SnapSet(t, snap, "apps."+app+".config.edgex-config-provider", "none")

	t.Logf("Copying coommon config file from platform snap to service snap: %s", snap)

	sourceFile := "/snap/edgexfoundry/current/config/core-common-config-bootstrapper/res/configuration.yaml"
	destFile := "/var/snap/" + snap + "/current/config/common-config.yaml"
	// read the source common config file
	source, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		t.Fatal(err)
	}

	// write the source file contents to the destination file
	err = ioutil.WriteFile(destFile, source, 0644)
	if err != nil {
		t.Fatal(err)
	}

	utils.SnapSet(t, snap, "apps."+app+".config.edgex-common-config", destFile)
}
