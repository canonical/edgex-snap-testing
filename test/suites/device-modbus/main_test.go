package test

import (
	"edgex-snap-testing/test/utils"
	"log"
	"os"
	"testing"
)

const deviceModbusSnap = "edgex-device-modbus"
const deviceModbusService = "edgex-device-modbus.device-modbus"

func TestMain(m *testing.M) {

	log.Println("[SETUP]")

	// start clean
	utils.SnapRemove(nil,
		deviceModbusSnap,
		"edgexfoundry",
	)

	// install the device snap before edgexfoundry
	// to catch build error sooner and stop
	if utils.LocalSnap != "" {
		utils.SnapInstallFromFile(nil, utils.LocalSnap)
	} else {
		utils.SnapInstallFromStore(nil, deviceModbusSnap, utils.ServiceChannel)
	}
	utils.SnapInstallFromStore(nil, "edgexfoundry", utils.PlatformChannel)

	// make sure all services are online before starting the tests
	utils.WaitServiceOnline(nil, 60, utils.PlatformPorts...)

	// for local build, the interface isn't auto-connected.
	// connect manually regardless
	utils.SnapConnect(nil,
		"edgexfoundry:edgex-secretstore-token",
		deviceModbusSnap+":edgex-secretstore-token",
	)

	exitCode := m.Run()

	log.Println("[TEARDOWN]")

	utils.SnapDumpLogs(nil, deviceModbusSnap)

	utils.SnapRemove(nil,
		deviceModbusSnap,
		"edgexfoundry",
	)

	os.Exit(exitCode)
}
