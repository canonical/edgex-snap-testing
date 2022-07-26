package utils

import (
	"testing"
)

type Refresh struct {
	TestRefreshServices    bool
	TestRefreshConfigPaths bool
}

func TestRefresh(t *testing.T, snapName string, conf Refresh) {
	t.Run("refreshing", func(t *testing.T) {
		if conf.TestRefreshServices {
			testRefreshServices(t, snapName)
		}
		if conf.TestRefreshConfigPaths {
			testRefreshConfigPaths()
		}
	})
}

func testRefreshServices(t *testing.T, snapName string) {
	const (
		upgradedVersion = "latest/beta"
		platformSnap    = "edgexfoundry"
	)

	t.Cleanup(func() {
		if LocalSnap != "" {
			SnapInstallFromFile(t, LocalSnap)
		} else {
			SnapInstallFromStore(t, platformSnap, ServiceChannel)
		}
		WaitPlatformOnline(t)
	})

	if LocalSnap != "" {
		originalVersion := SnapVersion(t, snapName)
		SnapRefresh(t, snapName, upgradedVersion)
		WaitPlatformOnline(t)
		t.Logf("Successfully upgraded:\n\tfrom: %s\n\tto: %s", originalVersion, upgradedVersion)
	} else {
		SnapRefresh(t, snapName, upgradedVersion)
		WaitPlatformOnline(t)
		t.Logf("Successfully upgraded:\n\tfrom: %s\n\tto: %s", PlatformChannel, upgradedVersion)
	}
}

func testRefreshConfigPaths() {}
