package test

import (
	"edgex-snap-testing/test/utils"
	"testing"
	"time"
)

func TestConfigOption(t *testing.T) {
	const (
		// eKuiper will be configured to use this port to connect to Redis
		newPort = "11111"
		portKey = "config.edgex.default.port"
	)

	var defaultRedisPort = utils.ServicePort("redis")

	t.Cleanup(func() {
		utils.SnapUnset(t, ekuiperSnap, portKey)
		utils.SnapStop(t, ekuiperSnap)
	})

	t.Run("Test config option", func(t *testing.T) {
		t.Logf("Set EdgeX Redis Message bus port to %s and verify", newPort)
		utils.SnapSet(t, ekuiperSnap, portKey, newPort)
		startTime := time.Now()
		utils.SnapStart(t, ekuiperSnap)

		utils.WaitServiceOnline(t, 60,
			utils.ServicePort(ekuiperApp),
			utils.ServicePort(ekuiperRestfulApi),
		)

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

		utils.WaitForLogMessage(t, ekuiperSnap, "port:"+newPort, startTime)

		t.Log("Unset and check default port")
		utils.SnapUnset(t, ekuiperSnap, portKey)
		startTime = time.Now()
		utils.SnapRestart(t, ekuiperSnap)
		utils.WaitForLogMessage(t, ekuiperSnap, "port:"+defaultRedisPort, startTime)
	})
}
