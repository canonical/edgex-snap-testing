package test

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"edgex-snap-testing/test/utils"
	"github.com/stretchr/testify/require"
)

const coreDataPingEndpoint = "https://localhost:8443/core-data/api/v3/ping"

// TestAddAPIGatewayUser creates an example user, generates a JWT token for this user,
// and then accesses the core-data service via the API gateway using the JWT token.
// https://docs.edgexfoundry.org/3.0/getting-started/Ch-GettingStartedSnapUsers/#adding-api-gateway-users
func TestAddAPIGatewayUser(t *testing.T) {
	t.Log("Create an example user and generate a JWT token")
	idToken := utils.LoginTestUser(t)

	t.Log("Calling on behalf of example user:", coreDataPingEndpoint)

	req, err := http.NewRequest("GET", coreDataPingEndpoint, nil)
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

// TestChangeTLSCert creats new TLS certificate and calls API gateway to verify the use of new certificate
// https://docs.edgexfoundry.org/3.0/getting-started/Ch-GettingStartedSnapUsers/#changing-tls-certificates
func TestChangeTLSCert(t *testing.T) {
	const tmpDir = "./tmp"
	t.Cleanup(func() {
		require.NoError(t, os.RemoveAll(tmpDir))
	})

	t.Log("Create TLS certificate")
	createTLSCert(t)

	t.Log("Calling API gateway using new TLS certificates:", coreDataPingEndpoint)
	const caCertFile = tmpDir + "/ca.cert"

	// Note: %%	is a literal percent sign
	// Note: The ca.cert file is created by create-tls-certificates.sh
	code, _, _ := utils.Exec(t, fmt.Sprintf(
		"curl --verbose --show-error --silent --include --output /dev/null --write-out '%%{http_code}' --cacert %s '%s'",
		caCertFile,
		coreDataPingEndpoint))

	// A success response should return status 401 because the endpoint is protected.
	require.Equal(t, "401", strings.TrimSpace(code))
}

func createTLSCert(t *testing.T) {
	// The script path relative to the testing suites
	const createLtsCert = "../../utils/create-tls-certificates.sh"
	utils.Exec(t, createLtsCert)
}
