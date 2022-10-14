package test

import (
	"edgex-snap-testing/test/utils"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestServiceStatus(t *testing.T) {
	var securityServices = []string{"kong-daemon", "postgres", "vault"}
	
	t.Run("security services", func(t *testing.T) {
		for _, service := range securityServices {
			require.Equal(t, "disabled", utils.SnapServicesStartup(t, "edgexfoundry."+service))
			require.Equal(t, "inactive", utils.SnapServicesCurrent(t, "edgexfoundry."+service))
		}
	})
}

func TestAccess(t *testing.T) {
	t.Run("consul", func(t *testing.T) {
		t.Log("Access Consul locally")
		resp, err := http.Get("http://localhost:8500/v1/kv/edgex/core/2.0/core-data/Service/Port")
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, 200, resp.StatusCode)
	})
}
