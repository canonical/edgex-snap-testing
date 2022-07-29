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

		// The command should not return error even if nothing is grepped, hence the "|| true"
		stdout, stderr := exec(t, fmt.Sprintf(`sudo grep -RnI "%s/%s" /var/snap/%s/current || true`,
			snapName, originalRevision, snapName),
			true)
		require.Empty(t, stdout,
			fmt.Sprintf(`files not upgraded to use "current" symlink in config files:%s`, stdout))
		require.Empty(t, stderr)
	})
}
