package common

import (
	"edgex-snap-testing/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVersion(t *testing.T, snapName string) {
	require.Regexp(t,
		"^([0-9]+).([0-9]+).([0-9]+).*$",
		utils.SnapVersion(t, snapName),
	)
}
