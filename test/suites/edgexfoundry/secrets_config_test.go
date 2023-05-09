package test

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"edgex-snap-testing/test/utils"
)
// TestAddAPIGatewayUser creates an example user, generates a JWT token for this user,
// and then accesses the core-data service via the API gateway using the JWT token.
// https://docs.edgexfoundry.org/3.0/getting-started/Ch-GettingStartedSnapUsers/#adding-api-gateway-users
func TestAddAPIGatewayUser(t *testing.T) {
	t.Log("Create an example user and generate a JWT token")
	var idToken string
	idToken = utils.LoginTestUser(t)

	t.Log("Call an API on behalf of example user")
	coreDataEndpoint := "https://localhost:8443/core-data/api/v3/ping"

	req, err := http.NewRequest("GET", coreDataEndpoint, nil)
	require.NoError(t, err)

	req.Header.Set("Authorization", "Bearer "+idToken)

	// InsecureSkipVerify because the API Gateway uses the built-in self-signed certificate
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, 200, resp.StatusCode, "Unexpected HTTP response")

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	t.Logf("Output: %s", body)
}
