package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// RequireSnapSemver checks that a snap has semantic versioning
func RequireSnapSemver(t *testing.T, snapName string) {
	require.Regexp(t,
		"^([0-9]+).([0-9]+).([0-9]+).*$",
		SnapVersion(t, snapName),
	)
}
