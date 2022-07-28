package utils

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"regexp"
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

	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		err := os.Chdir(currentDir)
		if err != nil {
			t.Fatal(err)
		}

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
		files, err := walkMatch(fmt.Sprintf("/var/snap/%s/%s", snapName, refreshRevision),
			fmt.Sprintf("%s/%s", snapName, originalRevision),
			"consul/data/raft/raft.db")
		require.NoError(t, err)
		require.Empty(t, files)
	})
}

func walkMatch(root, pattern, excludedPattern string) ([]string, error) {
	err := os.Chdir(root)
	if err != nil {
		return nil, err
	}

	var matches []string
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		matched, err := regexp.MatchString(pattern, path)
		if err != nil {
			return err
		}

		matchedConsul, err := regexp.MatchString(excludedPattern, path)
		if err != nil {
			return err
		}

		if matched && !matchedConsul {
			matches = append(matches, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return matches, nil
}
