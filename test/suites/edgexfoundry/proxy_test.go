package test

import (
	"crypto/tls"
	"crypto/x509"
	"edgex-snap-testing/test/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// Test seeding an admin user using snap options
// https://docs.edgexfoundry.org/2.2/getting-started/Ch-GettingStartedSnapUsers/#adding-api-gateway-users
func TestAddProxyUser(t *testing.T) {
	const (
		tmpDir         = "./tmp"
		publicKeyFile  = tmpDir + "/public.pem"
		privateKeyFile = tmpDir + "/private.pem"
	)

	// start clean
	require.NoError(t, os.RemoveAll(tmpDir))

	t.Cleanup(func() {
		require.NoError(t, os.RemoveAll(tmpDir))
		utils.SnapUnset(t, platformSnap, "apps")
		utils.SnapUnset(t, platformSnap, "app-options")
	})

	t.Log("Generate private and public keys")
	require.NoError(t, os.Mkdir(tmpDir, 0755))
	utils.Exec(t, fmt.Sprintf("openssl ecparam -genkey -name prime256v1 -noout -out %s", privateKeyFile))
	utils.Exec(t, fmt.Sprintf("openssl ec -in %s -pubout -out %s", privateKeyFile, publicKeyFile))

	t.Log("Add the public key for admin user")
	publicKey, err := os.ReadFile(publicKeyFile)
	require.NoError(t, err)

	utils.SnapSet(t, platformSnap, "app-options", "true")
	utils.SnapSet(t, platformSnap, "apps.secrets-config.proxy.admin.public-key", string(publicKey))

	t.Log("Generate a JWT token for the admin user")
	// The seedable "admin" has id 1
	jwt, _ := utils.Exec(t,
		fmt.Sprintf("edgexfoundry.secrets-config proxy jwt --algorithm ES256 --private_key %s --id 1 --expiration=1h", privateKeyFile))

	t.Log("Call an API on behalf of admin user")
	req, err := http.NewRequest(http.MethodGet, "https://localhost:8443/core-data/api/v2/ping", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", strings.TrimSpace(jwt)))

	// InsecureSkipVerify because the proxy uses the built-in self-signed certificate
	client := http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, 200, resp.StatusCode)
}

// Test seeding a custom TLS certificate using snap options
// https://docs.edgexfoundry.org/2.2/getting-started/Ch-GettingStartedSnapUsers/#changing-tls-certificates
func TestTLSCert(t *testing.T) {
	t.Cleanup(func() {
		utils.SnapUnset(t, platformSnap, "apps")
		utils.SnapUnset(t, platformSnap, "app-options")
	})

	t.Logf("Generate a self-signed certificate")
	_, caCertFile, _, serverKeyFile, serverCertFile := generateCerts(t)

	serverKey, err := os.ReadFile(serverKeyFile)
	require.NoError(t, err)
	serverCert, err := os.ReadFile(serverCertFile)
	require.NoError(t, err)

	t.Logf("Add the self-signed certificate")
	utils.SnapSet(t, platformSnap, "app-options", "true")
	// The options must be set together
	utils.Exec(t, fmt.Sprintf(
		"sudo snap set %s apps.secrets-config.proxy.tls.key='%s' apps.secrets-config.proxy.tls.cert='%s'",
		platformSnap,
		serverKey,
		serverCert,
	))

	t.Logf("Verify certificate installation by querying Kong's Admin API")
	resp, err := http.Get("http://localhost:8001/certificates")
	require.NoError(t, err)
	defer resp.Body.Close()

	var res struct{ Data []struct{ Cert string } }
	err = json.NewDecoder(resp.Body).Decode(&res)
	require.NoError(t, err)
	require.Len(t, res.Data, 1)
	require.Equal(t, res.Data[0].Cert, string(serverCert))

	t.Logf("Query a service via the proxy to verify the use of new certificate")
	caCert, err := ioutil.ReadFile(caCertFile)
	require.NoError(t, err)
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	client := http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs: caCertPool,
		},
	}}

	const maxRetry = 30
	for i := 1; i <= maxRetry; i++ {
		resp, err = client.Get("https://localhost:8443/core-data/api/v2/ping")
		if err != nil {
			t.Logf("Retry %d/%d: %s", i, maxRetry, err)
			time.Sleep(time.Second)
			continue
		}
		// A success response should return status 401 because the endpoint is protected.
		require.Equal(t, 401, resp.StatusCode)
	}
}

// generateCerts generates CA private key, CA cert,
// server private key, server signing request, server cert
func generateCerts(t *testing.T) (caKeyFile, caCertFile, serverCsrFile, serverKeyFile, serverCertFile string) {
	const tmpDir = "./tmp"

	// start clean
	require.NoError(t, os.RemoveAll(tmpDir))

	t.Cleanup(func() {
		require.NoError(t, os.RemoveAll(tmpDir))
	})

	// Create temp dir for certificates and keys
	require.NoError(t, os.Mkdir(tmpDir, 0755))

	caKeyFile = tmpDir + "/ca.key"
	caCertFile = tmpDir + "/ca.crt"
	serverCsrFile = tmpDir + "/server.csr"
	serverKeyFile = tmpDir + "/server.key"
	serverCertFile = tmpDir + "/server.crt"

	// Generate the Certificate Authority (CA) Private Key
	utils.Exec(nil, fmt.Sprintf("openssl ecparam -name prime256v1 -genkey -noout -out %s",
		caKeyFile))
	// Generate the Certificate Authority Certificate
	utils.Exec(nil, fmt.Sprintf("openssl req -new -x509 -sha256 -key %s -out %s -subj '/CN=snap-testing-ca' -addext 'subjectAltName = DNS:localhost'",
		caKeyFile, caCertFile))

	// Generate the Server Certificate Private Keys
	utils.Exec(nil, fmt.Sprintf("openssl ecparam -name prime256v1 -genkey -noout -out %s",
		serverKeyFile))
	// Generate the Server Certificate Signing Request
	utils.Exec(nil, fmt.Sprintf("openssl req -new -sha256 -key %s -out %s -subj '/CN=snap-tests' -addext 'subjectAltName = DNS:localhost'",
		serverKeyFile, serverCsrFile))
	// Generate the Server Certificate
	utils.Exec(nil, fmt.Sprintf("openssl x509 -req -in %s -CA %s -CAkey %s -CAcreateserial -out %s -days 1 -sha256",
		serverCsrFile, caCertFile, caKeyFile, serverCertFile))

	return
}
