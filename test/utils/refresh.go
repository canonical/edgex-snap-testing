package utils

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

type Refresh struct {
	TestRefreshServicesAndConfigPaths bool
}

func TestRefresh(t *testing.T, snapName string, conf Refresh) {
	t.Run("refresh", func(t *testing.T) {
		if conf.TestRefreshServicesAndConfigPaths {
			testRefresh(t, snapName)
		}
	})
}

func testRefresh(t *testing.T, snapName string) {
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
	originalRevision := SnapRevision(t, snapName)

	t.Run("refresh services", func(t *testing.T) {
		SnapRefresh(t, snapName, refreshChannel)
		refreshVersion := SnapVersion(t, snapName)
		refreshRevision = SnapRevision(t, snapName)
		WaitPlatformOnline(t)

		t.Logf("Successfully upgraded from %s(%s) to %s(%s)",
			originalVersion, originalRevision, refreshVersion, refreshRevision)
	})

	t.Run("refresh config paths", func(t *testing.T) {
		if originalRevision == refreshRevision {
			t.Skip("Upgraded to the same revision. Skipping test")
		}

		t.Logf("Checking for files with original snap revision %s", originalRevision)
		// exclude the file consul/data/raft/raft.db which has an old revision number in the path
		stdout, stderr := exec(t, fmt.Sprintf(`cd /var/snap/%s/current && grep --dereference-recursive --line-number %s/%s | grep --invert-match %s`,
			snapName, snapName, originalRevision, "raft.db"),
			true)
		require.Empty(t, stdout)
		require.Empty(t, stderr)
	})
}
