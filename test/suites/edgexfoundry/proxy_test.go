package test

import (
	"crypto/tls"
	"edgex-snap-testing/test/utils"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var kongAdminJwtFile = "/var/snap/edgexfoundry/current/secrets/security-proxy-setup/kong-admin-jwt"

func TestAddProxyUser(t *testing.T) {

	var tmpDir = "keys"

	// start clean
	err := os.RemoveAll(tmpDir)
	require.NoError(t, err)

	t.Cleanup(func() {
		err = os.RemoveAll(tmpDir)
		require.NoError(t, err)
	})

	// Get Kong admin JWT token
	out, err := os.ReadFile(kongAdminJwtFile)
	require.NoError(t, err)
	kongAdminJWT := string(out)
	require.NotEmpty(t, kongAdminJWT)

	// Create temp dir for private and public keys
	err = os.Mkdir(tmpDir, 0755)
	require.NoError(t, err)

	publicKey := tmpDir + "/public.pem"
	privateKey := tmpDir + "/private.pem"

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
}

func TestTLSCert(t *testing.T) {

	var tmpDir = "certs"

	// start clean
	err := os.RemoveAll(tmpDir)
	require.NoError(t, err)

	t.Cleanup(func() {
		err = os.RemoveAll(tmpDir)
		require.NoError(t, err)
	})

	// Get Kong admin JWT token
	out, err := os.ReadFile(kongAdminJwtFile)
	require.NoError(t, err)
	kongAdminJWT := string(out)
	require.NotEmpty(t, kongAdminJWT)

	// Create temp dir for certificates and keys
	err = os.Mkdir(tmpDir, 0755)
	require.NoError(t, err)

	// Add the certificate, using Kong Admin JWT to authenticate
	caCertFile, caKeyFile := certGenerator(tmpDir)
	utils.Exec(t, fmt.Sprintf("edgexfoundry.secrets-config proxy tls --incert %s --inkey %s --admin_api_jwt %s", caCertFile, caKeyFile, kongAdminJWT))

	code, _ := utils.Exec(t, fmt.Sprintf(`curl --show-error --silent --include --output /dev/null --write-out "%%{http_code}" --cacert %s -X GET 'https://localhost:8443/core-data/api/v2/ping?' -H "Authorization: Bearer %s"`, caCertFile, "BadToken"))
	require.Equal(t, "401\n", code)
}

func certGenerator(tmpDir string) (string, string) {
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
