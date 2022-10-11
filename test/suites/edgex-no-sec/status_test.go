package test

import (
	"crypto/tls"
	"edgex-snap-testing/test/utils"
	"github.com/stretchr/testify/require"
	"net/http"
	"strings"
	"testing"
)

func TestServiceStatus(t *testing.T) {
	var securityServices = []string{"kong-daemon", "postgres", "vault"}
	
	t.Run("security services", func(t *testing.T) {
		for _, service := range securityServices {
			require.Equal(t, "inactive", utils.SnapServices(t, "edgexfoundry."+service))
		}
	})
}

func TestAccess(t *testing.T) {
	t.Run("consul", func(t *testing.T) {
		t.Log("Generate a Consul token")
		consulToken, _, _ := utils.Exec(t, "sudo cat /var/snap/edgexfoundry/current/secrets/consul-acl-token/bootstrap_token.json | jq -r '.SecretID'")

		t.Log("Access Consul locally")
		req, err := http.NewRequest(http.MethodGet, "http://localhost:8500/v1/kv/edgex/core/2.0/core-data/Service/Port", nil)
		require.NoError(t, err)
		req.Header.Set("X-Consul-Token", strings.TrimSpace(consulToken))

		// InsecureSkipVerify
		client := http.Client{Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, 200, resp.StatusCode)
	})
}
