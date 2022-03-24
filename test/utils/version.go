package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func CheckSemver(t *testing.T, snapName string) {
	require.Regexp(t,
		"^([0-9]+).([0-9]+).([0-9]+).*$",
		SnapVersion(t, snapName),
	)
}
