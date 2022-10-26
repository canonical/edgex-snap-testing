package utils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

type Autostart struct {
	TestAutoStart bool
}

func TestAutoStart(t *testing.T, snapName string, conf Autostart) {
	t.Run("autostart", func(t *testing.T) {
		if conf.TestAutoStart {
			testAutoStart(t, snapName)
		}
	})
}

func testAutoStart(t *testing.T, snapName string) {
	t.Run("set and unset autostart", func(t *testing.T) {
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
