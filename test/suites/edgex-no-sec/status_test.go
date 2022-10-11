package test

import (
	"edgex-snap-testing/test/utils"
	"github.com/stretchr/testify/require"
	"testing"
)

var securityServices = []string{"kong-daemon", "postgres", "vault"}

func TestStatus(t *testing.T) {
	t.Run("security services", func(t *testing.T) {
		for _, service := range securityServices {
			require.Equal(t, "inactive", utils.SnapServices(t, "edgexfoundry."+service))
		}
	})
}
