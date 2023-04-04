package utils

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"text/template"
	"time"

	"github.com/stretchr/testify/require"
)

type Config struct {
	TestChangePort ConfigChangePort
	TestAutoStart  bool
}

type ConfigChangePort struct {
	App                      string
	DefaultPort              string
	TestAppConfig            bool
	TestGlobalConfig         bool
	TestMixedGlobalAppConfig bool
}

const serviceWaitTimeout = 60 // seconds

func TestConfig(t *testing.T, snapName string, conf Config) {
	t.Run("config", func(t *testing.T) {
		TestChangePort(t, snapName, conf.TestChangePort)
		TestAutoStart(t, snapName, conf.TestAutoStart)
	})
}

func TestChangePort(t *testing.T, snapName string, conf ConfigChangePort) {
	if conf.TestAppConfig || conf.TestGlobalConfig || conf.TestMixedGlobalAppConfig {
		t.Run("change service port", func(t *testing.T) {

			// start once so that default configs get uploaded to the registry
			service := snapName + "." + conf.App
			SnapStart(t, service)
			WaitServiceOnline(t, 60, conf.DefaultPort)
			SnapStop(t, service)

			if conf.TestAppConfig {
				testChangePort_app(t, snapName, conf.App, conf.DefaultPort)
			}
			if conf.TestGlobalConfig {
				testChangePort_global(t, snapName, conf.App, conf.DefaultPort)
			}
			if conf.TestMixedGlobalAppConfig {
				testChangePort_mixedGlobalApp(t, snapName, conf.App, conf.DefaultPort)
			}
		})
	}
}

func testChangePort_app(t *testing.T, snap, app, servicePort string) {
	t.Run("app config", func(t *testing.T) {
		service := snap + "." + app

		// start clean
		SnapStop(t, service)

		t.Cleanup(func() {
			SnapUnset(t, snap, "apps")
			SnapStop(t, service)
		})

		const newPort = "22222"

		// make sure the port is available before using it
		RequirePortAvailable(t, newPort)

		DisableConfigProviderServiceSnap(t, snap, app)

		// set apps. and validate the new port comes online
		SnapSet(t, snap, "apps."+app+".config.service-port", newPort)
		SnapStart(t, service)

		WaitServiceOnline(t, serviceWaitTimeout, newPort)

		// unset apps. and validate the default port comes online
		SnapUnset(t, snap, "apps."+app+".config.service-port")
		SnapRestart(t, service)

		WaitServiceOnline(t, serviceWaitTimeout, servicePort)
	})
}

func testChangePort_global(t *testing.T, snap, app, servicePort string) {
	t.Run("global config", func(t *testing.T) {
		service := snap + "." + app

		// start clean
		SnapStop(t, service)

		t.Cleanup(func() {
			SnapUnset(t, snap, "config")
			SnapStop(t, service)
		})

		const newPort = "33333"

		// make sure the port is available before using it
		RequirePortAvailable(t, newPort)

		DisableConfigProviderServiceSnap(t, snap, app)

		// set config. and validate the new port comes online
		SnapSet(t, snap, "config.service-port", newPort)
		SnapStart(t, service)

		WaitServiceOnline(t, serviceWaitTimeout, newPort)

		// unset config. and validate the default port comes online
		SnapUnset(t, snap, "config.service-port")
		SnapRestart(t, service)

		WaitServiceOnline(t, serviceWaitTimeout, servicePort)
	})
}

func testChangePort_mixedGlobalApp(t *testing.T, snap, app, servicePort string) {
	t.Run("app+global config for different values", func(t *testing.T) {
		service := snap + "." + app

		if !FullConfigTest {
			t.Skip("Full config test is disabled.")
		}
		// start clean
		SnapStop(t, service)

		t.Cleanup(func() {
			SnapUnset(t, snap, "apps")
			SnapUnset(t, snap, "config")
			SnapStop(t, service)
		})

		const newAppPort = "44444"
		const newConfigPort = "55555"

		// make sure the ports are available before using it
		RequirePortAvailable(t, newAppPort)
		RequirePortAvailable(t, newConfigPort)

		DisableConfigProviderServiceSnap(t, snap, app)

		// set apps. and config. with different values,
		// and validate that app-specific option has been picked up because it has higher precedence
		SnapSet(t, snap, "apps."+app+".config.service-port", newAppPort)
		SnapSet(t, snap, "config.service-port", newConfigPort)
		SnapStart(t, service)

		WaitServiceOnline(t, serviceWaitTimeout, newAppPort)
	})

}

func TestAutoStart(t *testing.T, snapName string, testAutoStart bool) {
	t.Run("autostart", func(t *testing.T) {
		if testAutoStart {
			TestAutostartGlobal(t, snapName)
		}
	})
}

func TestAutostartGlobal(t *testing.T, snapName string) {
	t.Run("set and unset global autostart", func(t *testing.T) {
		t.Cleanup(func() {
			SnapUnset(t, snapName, "autostart")
			SnapStop(t, snapName)
		})

		SnapStop(t, snapName)
		require.False(t, SnapServicesEnabled(t, snapName))
		require.False(t, SnapServicesActive(t, snapName))

		SnapSet(t, snapName, "autostart", "true")
		require.True(t, SnapServicesEnabled(t, snapName))
		require.True(t, SnapServicesActive(t, snapName))

		SnapUnset(t, snapName, "autostart")
		require.True(t, SnapServicesEnabled(t, snapName))
		require.True(t, SnapServicesActive(t, snapName))

		SnapSet(t, snapName, "autostart", "false")
		require.False(t, SnapServicesEnabled(t, snapName))
		require.False(t, SnapServicesActive(t, snapName))
	})
}

func WaitStartupMsg(t *testing.T, snap, expectedMsg string, since time.Time, retries int) {
	for i := 1; i <= retries; i++ {
		time.Sleep(1 * time.Second)
		t.Logf("Retry %d/%d: Waiting for startup message: %s", i, retries, expectedMsg)

		logs := SnapLogs(t, since, snap)
		if strings.Contains(logs, fmt.Sprintf("msg=%s", expectedMsg)) ||
			strings.Contains(logs, fmt.Sprintf(`msg="%s"`, expectedMsg)) {
			t.Logf(`Found startup message: "%s"`, expectedMsg)
			return
		}
	}
	t.Fatalf(`Time out: reached max %d retries looking for "%s"`, retries, expectedMsg)
}

// DisableConfigProviderPlatformSnap disables the config provider for the specified app
// and sets the common configuration path
func DisableConfigProviderPlatformSnap(t *testing.T, snap, app string) {
	SnapSet(t, snap, "apps."+app+".config.edgex-config-provider", "none")
	SnapSet(t, snap, "apps."+app+".config.edgex-common-config", "./config/core-common-config-bootstrapper/res/configuration.yaml")
}

// DisableConfigProviderServiceSnap disables the config provider for the specified app,
// copies the common configuration file from the platform snap to the service snap,
// and sets the common configuration path.
func DisableConfigProviderServiceSnap(t *testing.T, snap, app string) {
	SnapSet(t, snap, "apps."+app+".config.edgex-config-provider", "none")

	t.Logf("Copying common config file from platform snap to service snap: %s", snap)

	sourceFile := "/snap/edgexfoundry/current/config/core-common-config-bootstrapper/res/configuration.yaml"
	destFile := "/var/snap/" + snap + "/current/config/common-config.yaml"
	Exec(t, "sudo cp "+sourceFile+" "+destFile)

	SnapSet(t, snap, "apps."+app+".config.edgex-common-config", destFile)
}

// TODO: remove once issue (TBA) is resolved
func injectDevicesAndProfilesDirConfig(app string) error {
	const confPathTempl = "/var/snap/edgex-{{.}}/current/config/{{.}}/res/configuration.yaml"
	const confTempl = `
# Local config to override the common setting
Device:
  ProfilesDir: /var/snap/edgex-{{.}}/current/config/{{.}}/res/profiles
  DevicesDir: /var/snap/edgex-{{.}}/current/config/{{.}}/res/devices
`
	var tpl bytes.Buffer
	template.Must(template.New("temp").Parse(confTempl)).Execute(&tpl, app)
	conf := tpl.String()

	tpl.Reset()
	template.Must(template.New("temp").Parse(confPathTempl)).Execute(&tpl, app)
	confPath := tpl.String()

	stdout, stderr, err := exec(nil, "echo '"+conf+"' | sudo tee -a "+confPath, true)
	if err != nil {
		return fmt.Errorf("%s %s %s", stdout, stderr, err)
	}
	return nil
}
