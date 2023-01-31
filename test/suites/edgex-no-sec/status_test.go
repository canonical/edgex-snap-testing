package test

import (
	"edgex-snap-testing/test/utils"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
	"fmt"
	"strings"
)

func TestServiceStatus(t *testing.T) {
	t.Run("security services", func(t *testing.T) {
		var securityServices = []string{"kong-daemon", "postgres", "vault"}

		for _, service := range securityServices {
			require.False(t, utils.SnapServicesEnabled(t, "edgexfoundry."+service))
			require.False(t, utils.SnapServicesActive(t, "edgexfoundry."+service))
		}
	})
}

func TestAccess(t *testing.T) {
	t.Run("consul", func(t *testing.T) {
		t.Log("Getting Consul token")
		consulToken, _, _ := utils.Exec(t, "sudo cat /var/snap/edgexfoundry/current/secrets/consul-acl-token/bootstrap_token.json | jq -r '.SecretID'")

		t.Log("Access Consul locally")
		client := &http.Client{}
		req, err := http.NewRequest("GET", "http://localhost:8500/v1/kv/edgex/v3/core-data/Service/Port", nil)
		require.NoError(t, err)
		req.Header.Set("X-Consul-Token", strings.TrimSpace(consulToken))
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, 200, resp.StatusCode)
	})
}
