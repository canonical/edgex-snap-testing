package test

import (
	"edgex-snap-testing/test/utils"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestConfigOption(t *testing.T) {
	const (
		defaultPort = "6379"
		newPort     = "6000"
		portKey     = "config.edgex.default.port"
	)

	t.Cleanup(func() {
		utils.SnapUnset(t, ekuiperSnap, portKey)
		utils.SnapRestart(t, ekuiperSnap)
	})

	t.Run("Test config option", func(t *testing.T) {
		t.Log("Set and verify new EdgeX Redis Message bus port to:", newPort)
		utils.SnapSet(t, ekuiperSnap, portKey, newPort)
		startTime := time.Now()
		utils.SnapStart(t, ekuiperSnap)

		t.Log("Creating stream and rule to trigger the process of applying config option in edgex-ekuiper")
		utils.Exec(t, `edgex-ekuiper.kuiper create stream stream1 '()WITH(FORMAT="JSON",TYPE="edgex")'`)
		utils.Exec(t,
			`edgex-ekuiper.kuiper create rule rule_log '
				{
					"sql":"SELECT * FROM stream1 WHERE meta(deviceName) != \"device-test\"",
					"actions":[
						{
							"log":{}
						}
					]
				}'`)

		require.True(t, waitForRedisPortInLogs(t, ekuiperSnap, newPort, startTime),
			"new port: %s", newPort)

		t.Log("Unset and check default port")
		utils.SnapUnset(t, ekuiperSnap, portKey)
		startTime = time.Now()
		utils.SnapRestart(t, ekuiperSnap)
		require.True(t, waitForRedisPortInLogs(t, ekuiperSnap, defaultPort, startTime),
			"default port: %s", defaultPort)
	})
}

func waitForRedisPortInLogs(t *testing.T, snap, expectedPort string, since time.Time) bool {
	const maxRetry = 60

	utils.WaitServiceOnline(t, 60,
		utils.ServicePort(ekuiperApp),
		utils.ServicePort(ekuiperRestfulApi),
	)

	for i := 1; i <= maxRetry; i++ {
		time.Sleep(1 * time.Second)
		t.Logf("Retry %d/%d: Waiting for expected port in logs: %s", i, maxRetry, expectedPort)

		logs := utils.SnapLogs(t, since, snap)
		if strings.Contains(logs, fmt.Sprintf("port:%s", expectedPort)) ||
			strings.Contains(logs, fmt.Sprintf(`port:%s`, expectedPort)) {
			t.Logf("Found expected port: %s", expectedPort)
			return true
		}
	}

	t.Logf("Time out: reached max %d retries.", maxRetry)
	return false
}
