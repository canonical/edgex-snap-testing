package utils

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

type Refresh struct {
	TestRefreshServicesAndConfigPaths bool
	Regexes                           []string
}

func TestRefresh(t *testing.T, snapName string, conf Refresh) {
	t.Run("refresh", func(t *testing.T) {
		if conf.TestRefreshServicesAndConfigPaths {
			testRefresh(t, snapName, conf.Regexes...)
		}
	})
}

func testRefresh(t *testing.T, snapName string, regexes ...string) {
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

		// exclude files that have an old revision number in the path using regex
		command := fmt.Sprintf(`cd /var/snap/%s/current && grep --dereference-recursive --line-number %s/%s | grep --invert-match `,
			snapName, snapName, originalRevision)
		for _, reg := range regexes {
			command += fmt.Sprintf("--word-regexp %s ", reg)
		}

		stdout, stderr := exec(t, command, true)
		require.Empty(t, stdout, fmt.Sprintf(`files not upgraded to use "current" symlink in config files:%s`, stdout))
		require.Empty(t, stderr)
	})
}
