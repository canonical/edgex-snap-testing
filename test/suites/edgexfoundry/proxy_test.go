package test

import (
	"crypto/tls"
	"edgex-snap-testing/test/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestAddProxyUser(t *testing.T) {
	const (
		tmpDir     = "keys"
		publicKey  = tmpDir + "/public.pem"
		privateKey = tmpDir + "/private.pem"
	)

	// start clean
	require.NoError(t, os.RemoveAll(tmpDir))

	t.Cleanup(func() {
		require.NoError(t, os.RemoveAll(tmpDir))
	})

	// Create temp dir for private and public keys
	require.NoError(t, os.Mkdir(tmpDir, 0755))

	// Get Kong admin JWT token
	utils.Exec(t, fmt.Sprintf("sudo install -m 604 /var/snap/edgexfoundry/current/secrets/security-proxy-setup/kong-admin-jwt ./%s", tmpDir))
	kongAdminJWTFile := tmpDir + "/kong-admin-jwt"
	kongAdminJWT, err := os.ReadFile(kongAdminJWTFile)
	require.NoError(t, err)

	// Generate private and public keys
	utils.Exec(t, fmt.Sprintf("openssl ecparam -genkey -name prime256v1 -noout -out %s", privateKey))
	utils.Exec(t, fmt.Sprintf("openssl ec -in %s -pubout -out %s", privateKey, publicKey))

	// Use secrets-config to add a user example with id 1000
	stdout, _ := utils.Exec(t,
		fmt.Sprintf("edgexfoundry.secrets-config proxy adduser --token-type jwt --user example --algorithm ES256 --public_key %s --id 1000 -jwt %s",
			publicKey,
			kongAdminJWT))
	// On success, the above command prints the user id
	require.Equal(t, "1000\n", stdout)

	// Generate a JWT token for the above example user
	jwt, _ := utils.Exec(t,
		fmt.Sprintf("edgexfoundry.secrets-config proxy jwt --algorithm ES256 --private_key %s --id 1000 --expiration=1h", privateKey))
	jwt = strings.TrimSuffix(jwt, "\n")

	// Call an API on behalf of example user
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", "https://localhost:8443/core-data/api/v2/ping", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwt))

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, 200, resp.StatusCode)
}

func TestTLSCert(t *testing.T) {

	const tmpDir = "certs"

	// start clean
	require.NoError(t, os.RemoveAll(tmpDir))

	t.Cleanup(func() {
		require.NoError(t, os.RemoveAll(tmpDir))
	})

	// Create temp dir for certificates and keys
	require.NoError(t, os.Mkdir(tmpDir, 0755))

	// Get Kong admin JWT token
	utils.Exec(t, fmt.Sprintf("sudo install -m 604 /var/snap/edgexfoundry/current/secrets/security-proxy-setup/kong-admin-jwt ./%s", tmpDir))
	kongAdminJWTFile := tmpDir + "/kong-admin-jwt"
	kongAdminJWT, err := os.ReadFile(kongAdminJWTFile)
	require.NoError(t, err)

	// Add the certificate, using Kong Admin JWT to authenticate
	caCertFile, caKeyFile := generateCerts(tmpDir)
	utils.Exec(t, fmt.Sprintf("edgexfoundry.secrets-config proxy tls --incert %s --inkey %s --admin_api_jwt %s", caCertFile, caKeyFile, kongAdminJWT))

	// Wait the certificate to be fully installed
	caKey, err := os.ReadFile(caKeyFile)
	require.NoError(t, err)
	waitCertInstall(t, string(caKey), 10)

	// Check if TLS is setup correctly returning status code 401
	code, _ := utils.Exec(t, fmt.Sprintf(`curl --show-error --silent --include --output /dev/null --write-out "%%{http_code}" --cacert %s -X GET 'https://localhost:8443/core-data/api/v2/ping' -H "Authorization: Bearer %s"`, caCertFile, "testToken"))
	require.Equal(t, "401\n", code)
}

// generateCerts generates CA private key, CA cert,
// server private key, server signing request, server cert
func generateCerts(tmpDir string) (string, string) {
	var (
		caKeyFile      = tmpDir + "/ca.key"
		caCertFile     = tmpDir + "/ca.crt"
		serverCsrFile  = tmpDir + "/server.csr"
		serverKeyFile  = tmpDir + "/server.key"
		serverCertFile = tmpDir + "/server.crt"
	)

	// Generate the Certificate Authority (CA) Private Key
	utils.Exec(nil, fmt.Sprintf("openssl ecparam -name prime256v1 -genkey -noout -out %s", caKeyFile))
	// Generate the Certificate Authority Certificate
	utils.Exec(nil, fmt.Sprintf(`openssl req -new -x509 -sha256 -key %s -out %s -subj "/CN=localhost"`, caKeyFile, caCertFile))

	// Generate the Server Certificate Private Keys
	utils.Exec(nil, fmt.Sprintf("openssl ecparam -name prime256v1 -genkey -noout -out %s", serverKeyFile))
	// Generate the Server Certificate Signing Request
	utils.Exec(nil, fmt.Sprintf(`openssl req -new -sha256 -key %s -out %s -subj "/CN=localhost"`, serverKeyFile, serverCsrFile))
	// Generate the Server Certificate
	utils.Exec(nil, fmt.Sprintf(`openssl x509 -req -in %s -CA %s -CAkey %s -CAcreateserial -out %s -days 1000 -sha256`, serverCsrFile, caCertFile, caKeyFile, serverCertFile))

	return caCertFile, caKeyFile
}

// waitCertInstall checks if certificate is fully installed
// up to a maximum retry number
func waitCertInstall(t *testing.T, caKey string, maxRetry int) {

	type ResponseBody struct {
		Data []struct {
			Key string
		}
	}

	var res ResponseBody

	for i := 1; ; i++ {
		t.Logf("Checking certificate installation. Retry %d/%d", i, maxRetry)

		resp, err := http.Get("http://localhost:8001/certificates")
		require.NoError(t, err)
		defer resp.Body.Close()

		err = json.NewDecoder(resp.Body).Decode(&res)
		require.NoError(t, err)

		totalCerts := len(res.Data)
		if i == maxRetry && totalCerts == 0 {
			t.Fatalf("Time out: reached max %d retries.", maxRetry)
		} else if totalCerts == 0 {
			time.Sleep(1 * time.Second)
		} else {
			break
		}
	}

	require.Equal(t, res.Data[0].Key, caKey)
}
