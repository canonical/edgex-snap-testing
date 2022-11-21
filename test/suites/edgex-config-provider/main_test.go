package test

import (
	"edgex-snap-testing/test/utils"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	platformSnap = "edgexfoundry"

	deviceVirtualSnap = "edgex-device-virtual"
)

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
	utils.SnapRemove(nil, platformSnap, deviceVirtualSnap)

	log.Println("[SETUP]")
	start := time.Now()

	teardown = func() {
		log.Println("[TEARDOWN]")

		utils.SnapDumpLogs(nil, start, platformSnap)
		utils.SnapDumpLogs(nil, start, deviceVirtualSnap)

		utils.SnapRemove(nil, platformSnap)
		utils.SnapRemove(nil, deviceVirtualSnap)

	}

	// clone the example provider
	const workDir = "edgex-config-provider"
	utils.Exec(nil, "git clone https://github.com/canonical/edgex-config-provider.git --depth=1 "+workDir)

	// TODO: add other config sources

	// change startup message, for the sake of testing

	utils.Exec(nil, fmt.Sprintf(`find %s -type f -name 'configuration.toml' | xargs \
    		sed --in-place --regexp-extended 's/StartupMsg.*/StartupMsg="%s"/'`,
		workDir, startupMsg))

	// build the example provider snap
	utils.SnapBuild(nil, workDir)

	const configProviderSnapFile = workDir + "/edgex-config-provider-example_2.3_amd64.snap"
	if err = utils.SnapInstallFromFile(nil, configProviderSnapFile); err != nil {
		teardown()
		return
	}

	if err = utils.SnapInstallFromStore(nil, deviceVirtualSnap, utils.ServiceChannel); err != nil {
		teardown()
		return
	}

	// connect
	const interfaceName = "device-virtual-config"
	if err = utils.SnapConnect(nil, deviceVirtualSnap+":"+interfaceName, "edgex-config-provider-example:"+interfaceName); err != nil {
		teardown()
		return
	}

	err = utils.SnapInstallFromStore(nil, platformSnap, utils.PlatformChannel)
	if err != nil {
		teardown()
		return
	}

	return
}

func TestConfigProvider(t *testing.T) {
	start := time.Now()
	utils.SnapStart(t, deviceVirtualSnap)

	require.True(t, checkStartupMsg(t, deviceVirtualSnap, startupMsg, start))
}

func checkStartupMsg(t *testing.T, snap, expectedMsg string, since time.Time) bool {
	const maxRetry = 10

	utils.WaitPlatformOnline(t)

	for i := 1; i <= maxRetry; i++ {
		time.Sleep(1 * time.Second)
		t.Logf("Waiting for startup message. Retry %d/%d", i, maxRetry)

		logs := utils.SnapLogs(t, since, snap)
		if strings.Contains(logs, fmt.Sprintf("msg=%s", expectedMsg)) ||
			strings.Contains(logs, fmt.Sprintf(`msg="%s"`, expectedMsg)) {
			t.Logf("Found startup message: %s", expectedMsg)
			return true
		}
	}
	t.Logf("Time out: reached max %d retries.", maxRetry)
	return false
}
