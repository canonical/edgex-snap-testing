package utils

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

type Refresh struct {
	TestRefreshServicesAndConfigPaths bool
}

func TestRefresh(t *testing.T, snapName, serviceName string, conf Refresh) {
	t.Run("refresh", func(t *testing.T) {
		if conf.TestRefreshServicesAndConfigPaths {
			testRefresh(t, snapName, serviceName)
		}
	})
}

func testRefresh(t *testing.T, snapName, serviceName string) {
	const refreshChannel = "latest/beta"
	var refreshRevision string

	t.Cleanup(func() {
		if LocalSnap != "" {
			SnapRemove(t, snapName)
			SnapInstallFromFile(t, LocalSnap)
		} else {
			SnapRemove(t, snapName)
			SnapInstallFromStore(t, snapName, ServiceChannel)
		}
		WaitPlatformOnline(t)
	})

	originalVersion := SnapVersion(t, snapName)
	originalRevision := SnapRevision(t, serviceName)

	t.Run("refresh services", func(t *testing.T) {
		SnapRefresh(t, snapName, refreshChannel)
		refreshVersion := SnapVersion(t, snapName)
		WaitPlatformOnline(t)

		t.Logf(`Successfully upgraded:
		from: %s
		to:   %s`,
			originalVersion, refreshVersion)

		refreshRevision = SnapRevision(t, serviceName)
		t.Logf(`Successfully upgraded:
		from: %s
		to:   %s`,
			originalRevision, refreshRevision)
	})

	t.Run("refresh config paths", func(t *testing.T) {
		if originalRevision == refreshRevision {
			t.Skip("Upgraded to the same revision. Skipping test")
		}

		stdout, _ := Exec(t, fmt.Sprintf(`cd /var/snap/%s/current && grep -R %s/%s | grep -v "raft.db"`,
			snapName, snapName, originalRevision))
		require.Empty(t, stdout)
	})
}
