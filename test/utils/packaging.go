package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type Packaging struct {
	TestSemanticSnapVersion bool
}

func TestPackaging(t *testing.T, snapName string, conf Packaging) {
	t.Run("packaging", func(t *testing.T) {
		if conf.TestSemanticSnapVersion {
			testSemanticSnapVersion(t, snapName)
		}
	})
}

func testSemanticSnapVersion(t *testing.T, snapName string) {
	t.Run("semantic snap version", func(t *testing.T) {
		require.Regexp(t,
			"^([0-9]+).([0-9]+).([0-9]+).*$",
			SnapVersion(t, snapName),
		)
	})
}
