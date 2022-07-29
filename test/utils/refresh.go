package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestRefresh tests an EdgeX upgrade using snap refresh
func TestRefresh(t *testing.T, snapName string) {
	t.Run("refresh", func(t *testing.T) {
		if LocalSnap != "" {
			t.Skip("Skip refresh for local snap build")
		}

		const stableChannel = "latest/stable"

		if ServiceChannel == stableChannel {
			t.Skipf("Skip refresh on same channel: %s", ServiceChannel)
		}

		// remove and install the older stable revision
		SnapRemove(t, snapName)
		SnapInstallFromStore(t, snapName, stableChannel)

		var refreshRevision string

		t.Cleanup(func() {
			SnapRemove(t, snapName)
			SnapInstallFromStore(t, snapName, ServiceChannel)

			WaitPlatformOnline(t)
		})

		originalVersion := SnapVersion(t, snapName)
		originalRevision := SnapRevision(t, snapName)

		t.Run("check services", func(t *testing.T) {
			SnapRefresh(t, snapName, ServiceChannel)
			refreshVersion := SnapVersion(t, snapName)
			refreshRevision = SnapRevision(t, snapName)
			WaitPlatformOnline(t)

			t.Logf("Successfully upgraded from %s (%s) to %s (%s)",
				originalVersion, originalRevision, refreshVersion, refreshRevision)
		})

		t.Run("check config paths", func(t *testing.T) {
			if originalRevision == refreshRevision {
				t.Skip("Upgraded to the same revision. Skipping test")
			}

			t.Logf("Looking for files containing previous snap revision %s", originalRevision)

			// The command should not return error even if nothing is grepped, hence the "|| true"
			stdout, stderr := exec(t,
				fmt.Sprintf("sudo grep -RnI '%s/%s' /var/snap/%s/current || true",
					snapName, originalRevision, snapName),
				true)
			require.Empty(t, stdout,
				"The following files contain revision %s instead of %s or 'current' symlink: %s",
				originalRevision, refreshRevision, stdout)
			require.Empty(t, stderr)
		})
	})
}
