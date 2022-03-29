package test

import (
	"edgex-snap-testing/test/utils"
	"testing"
)

// Deprecated
func TestEnvConfig(t *testing.T) {
	// start clean
	utils.Exec(t, "sudo snap stop edgex-app-service-configurable.app-service-configurable")

	t.Run("change service port", func(t *testing.T) {
		t.Cleanup(func() {
			utils.Exec(t, "sudo snap stop edgex-app-service-configurable.app-service-configurable")
			utils.Exec(t, "sudo snap unset edgex-app-service-configurable env.service.port")
		})

		const newPort = "56789"

		// make sure the port is available before using it
		utils.CheckPortAvailable(t, newPort)

		utils.Exec(t, "sudo snap set edgex-app-service-configurable env.service.port="+newPort)
		utils.Exec(t, "sudo snap start edgex-app-service-configurable")
		utils.WaitServiceOnline(t, newPort)
	})
}

func TestAppConfig(t *testing.T) {
	t.Skip("TODO")
}
