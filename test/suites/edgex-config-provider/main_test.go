package test

import (
	"edgex-snap-testing/test/utils"
	"log"
	"os"
	"testing"
	"time"
)

const (
	platformSnap = "edgexfoundry"
	provider     = "edgex-config-provider-example"
)

var services = []string{
	"app-service-configurable",
	"app-rfid-llrp-inventory",
	"device-gpio",
	"device-modbus",
	"device-mqtt",
	"device-rest",
	"device-rfid-llrp",
	"device-snmp",
	"device-usb-camera",
	"device-virtual",
	"device-onvif-camera",
}

const startupMsg = "CONFIG BY EXAMPLE PROVIDER"

func TestMain(m *testing.M) {
	teardown, err := setup()
	if err != nil {
		log.Fatalf("Failed to setup tests: %s", err)
	}

	code := m.Run()
	teardown()

	os.Exit(code)
}

func setup() (teardown func(), err error) {
	log.Println("[CLEAN]")
	utils.SnapRemove(nil, platformSnap)
	utils.SnapRemove(nil, provider)

	log.Println("[SETUP]")
	start := time.Now()

	teardown = func() {
		log.Println("[TEARDOWN]")
		utils.SnapDumpLogs(nil, start, platformSnap)
		utils.SnapRemove(nil, platformSnap)
		utils.SnapRemove(nil, provider)
		// remove cloned directory
		os.RemoveAll(provider)
	}

	// install the provider
	if utils.LocalSnap() {
		if err = utils.SnapInstallFromFile(nil, utils.LocalSnapPath); err != nil {
			teardown()
			return
		}
	} else {
		const workDir = provider + "/"
		// clone the example provider
		if _, _, err = utils.Exec(nil, "git clone https://github.com/canonical/edgex-config-provider.git --branch=snap-testing --depth=1 "+workDir); err != nil {
			teardown()
			return
		}

		// build the example provider snap
		if err = utils.SnapBuild(nil, workDir); err != nil {
			teardown()
			return
		}

		const configProviderSnapFile = workDir + provider + "_*_*.snap"
		if err = utils.SnapInstallFromFile(nil, configProviderSnapFile); err != nil {
			teardown()
			return
		}
	}

	if err = utils.SnapInstallFromStore(nil, platformSnap, utils.PlatformChannel); err != nil {
		teardown()
		return
	}

	utils.SnapStart(nil, platformSnap)

	// make sure all services are online before starting the tests
	if err = utils.WaitPlatformOnline(nil); err != nil {
		teardown()
		return
	}

	return
}

func TestConfigProvider(t *testing.T) {

	for _, name := range services {
		t.Run(name, func(t *testing.T) {
			snapName := "edgex-" + name

			// clean start
			utils.SnapRemove(t, snapName)

			start := time.Now()

			t.Cleanup(func() {
				utils.SnapDumpLogs(t, start, snapName)
				utils.SnapRemove(t, snapName)
			})

			// install the consumer
			utils.SnapInstallFromStore(t, snapName, utils.ServiceChannel)

			// add basic config to make the services run
			switch name {
			case "device-mqtt":
				const mqttBroker = "mosquitto"
				if !utils.SnapInstalled(t, mqttBroker) {
					utils.SnapInstallFromStore(t, mqttBroker, "latest/stable")
					t.Cleanup(func() {
						utils.SnapRemove(t, mqttBroker)
					})
				}
				if !utils.SnapServicesActive(t, mqttBroker) {
					utils.SnapStart(t, mqttBroker)
					t.Cleanup(func() {
						utils.SnapStop(t, mqttBroker)
					})
				}
			case "app-service-configurable":
				utils.SnapSet(nil, snapName, "profile", "http-export")
			}

			// NOTE: Connect after setting options to work around the error resulted
			// from creation/removal of an env file in read-only file system.
			// See https://warthogs.atlassian.net/browse/EDGEX-586
			//
			// connect to provider's slot
			interfaceName := name + "-config"
			utils.SnapConnect(t,
				snapName+":"+interfaceName,
				provider+":"+interfaceName)

			utils.SnapStart(t, snapName)

			utils.WaitForLogMessage(t, snapName, `msg="`+startupMsg+`"`, start)
		})
	}
}
