package utils

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
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

		t.Logf("Checking for files with original snap revision %s", originalRevision)
		files, err := WalkMatch(fmt.Sprintf("/var/snap/%s/current", snapName), fmt.Sprintf("*%s/%s*", snapName, originalRevision))
		require.NoError(t, err)
		require.Empty(t, files)
	})
}

func WalkMatch(root, pattern string) ([]string, error) {
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
			// exclude the file consul/data/raft/raft.db which has an old revision number in the path
		} else if matched && path != "consul/data/raft/raft.db" {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return matches, nil
}
