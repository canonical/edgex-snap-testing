package test

import (
	"edgex-snap-testing/test/utils"
	"log"
	"os"
	"testing"
)

const deviceSnmpSnap = "edgex-device-snmp"
const deviceSnmpService = "edgex-device-snmp.device-snmp"

func TestMain(m *testing.M) {

	log.Println("[SETUP]")

	// start clean
	utils.SnapRemove(nil,
		deviceSnmpSnap,
		"edgexfoundry",
	)

	// install the device snap before edgexfoundry
	// to catch build error sooner and stop
	if utils.LocalSnap != "" {
		utils.SnapInstallFromFile(nil, utils.LocalSnap)
	} else {
		utils.SnapInstallFromStore(nil, deviceSnmpSnap, utils.ServiceChannel)
	}
	utils.SnapInstallFromStore(nil, "edgexfoundry", utils.PlatformChannel)

	// make sure all services are online before starting the tests
	utils.WaitServiceOnline(nil, 60, utils.PlatformPorts...)

	// for local build, the interface isn't auto-connected.
	// connect manually regardless
	utils.SnapConnect(nil,
		"edgexfoundry:edgex-secretstore-token",
		deviceSnmpSnap+":edgex-secretstore-token",
	)

	exitCode := m.Run()

	log.Println("[TEARDOWN]")

	utils.SnapDumpLogs(nil, deviceSnmpSnap)

	utils.SnapRemove(nil,
		deviceSnmpSnap,
		"edgexfoundry",
	)

	os.Exit(exitCode)
}
