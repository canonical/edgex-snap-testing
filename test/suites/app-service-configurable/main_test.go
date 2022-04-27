package test

import (
	"edgex-snap-testing/test/utils"
	"log"
	"os"
	"testing"
	"time"
)

const (
	ascSnap                    = "edgex-app-service-configurable"
	ascApp                     = "app-service-configurable"
	ascService                 = ascSnap + "." + ascApp
	defaultProfile             = "rules-engine"
	appServiceRulesServicePort = "59701"
)

var start = time.Now()

func TestMain(m *testing.M) {

	log.Println("[SETUP]")

	// start clean
	utils.SnapRemove(nil,
		ascSnap,
		"edgexfoundry",
	)

	// install the app-service-configurable snap before edgexfoundry
	// to catch build error sooner and stop
	if utils.LocalSnap != "" {
		utils.SnapInstallFromFile(nil, utils.LocalSnap)
	} else {
		utils.SnapInstallFromStore(nil, ascSnap, utils.ServiceChannel)
	}
	utils.SnapInstallFromStore(nil, "edgexfoundry", utils.PlatformChannel)

	// make sure all services are online before starting the tests
	utils.WaitPlatformOnline(nil)

	// for local build, the interface isn't auto-connected.
	// connect manually regardless
	utils.SnapConnect(nil,
		"edgexfoundry:edgex-secretstore-token",
		ascSnap+":edgex-secretstore-token",
	)

	// set profile to rules engine
	utils.SnapSet(nil, ascSnap, "profile", defaultProfile)

	// Start the service so that the default config gets uploaded to consul.
	// Otherwise, settings that get passed using environment variables on first start get uploaded
	// and become the default.
	utils.SnapStart(nil, ascService)
	utils.WaitServiceOnline(nil, 60, appServiceRulesServicePort)

	exitCode := m.Run()

	log.Println("[TEARDOWN]")

	utils.SnapDumpLogs(nil, start, ascSnap)

	utils.SnapRemove(nil,
		ascSnap,
		"edgexfoundry",
	)

	os.Exit(exitCode)
}
