package test

import (
	"edgex-snap-testing/test/utils"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	log.Println("[SETUP]")

	// start clean
	utils.SnapRemove(nil,
		"edgex-ekuiper",
		"edgexfoundry",
	)

	// install the ekuiper snap before edgexfoundry
	// to catch build error sooner and stop
	if utils.LocalSnap != "" {
		utils.SnapInstallFromFile(nil, utils.LocalSnap)
	} else {
		utils.SnapInstallFromStore(nil, "edgex-ekuiper", utils.ServiceChannel)
	}
	utils.SnapInstallFromStore(nil, "edgexfoundry", utils.PlatformChannel)

	// for local build, the interface isn't auto-connected.
	// connect manually regardless
	utils.SnapConnect(nil,
		"edgexfoundry:edgex-secretstore-token",
		"edgex-ekuiper:edgex-secretstore-token",
	)
	utils.Exec(nil,
		"sudo snap restart edgex-ekuiper.kuiper",
	)

	// security on (default)
	exitCode := m.Run()
	if exitCode != 0 {
		goto TEARDOWN
	}

	// security off
	utils.Exec(nil, "sudo snap set edgexfoundry security-secret-store=off")
	utils.SnapRemove(nil,
		"edgex-ekuiper")

	if utils.LocalSnap != "" {
		utils.SnapInstallFromFile(nil, utils.LocalSnap)
	} else {
		utils.SnapInstallFromStore(nil, "edgex-ekuiper", utils.ServiceChannel)
	}

	exitCode = m.Run()

TEARDOWN:
	log.Println("[TEARDOWN]")

	if exitCode != 0 {
		log.Printf("Snap logs:\n%s\n",
			utils.SnapLogs(nil, "edgex-ekuiper"))
	}

	utils.SnapRemove(nil,
		"edgex-ekuiper",
		"edgexfoundry",
	)

	os.Exit(exitCode)
}
