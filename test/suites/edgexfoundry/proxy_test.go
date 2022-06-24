package test

import (
	"edgex-snap-testing/test/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	snapData        = "/var/snap/edgexfoundry/current/temp"
	kongAminJwtFile = "/var/snap/edgexfoundry/current/secrets/security-proxy-setup/kong-admin-jwt"
)

func TestAddProxyUser(t *testing.T) {
	const (
		publicKey  = snapData + "/public.pem"
		privateKey = snapData + "/private.pem"
	)

	// start clean
	utils.Exec(t, `sudo rm -rf `+snapData)

	t.Cleanup(func() {
		utils.Exec(t, `sudo rm -rf `+snapData)
	})

	// Due to confinement issues when running this test, we write the files to SNAP_DATA
	utils.Exec(t, `sudo mkdir -p `+snapData)

	// Read the API Gateway token
	kongAdminJWT, _ := utils.Exec(t, "sudo cat "+kongAminJwtFile)
	require.NotEmpty(t, kongAdminJWT)

	// Create private and public keys
	utils.Exec(t, `sudo openssl ecparam -genkey -name prime256v1 -noout -out `+privateKey)
	utils.Exec(t, `sudo openssl ec -in `+privateKey+` -pubout -out `+publicKey)

	// Use secrets-config to add a user example with id 1000
	// On success, the above command prints the user id
	stdout, _ := utils.Exec(t, `
			sudo edgexfoundry.secrets-config proxy adduser \
			--token-type jwt --user example --algorithm ES256 \
			--public_key `+publicKey+` --id 1000 -jwt `+kongAdminJWT)
	require.Equal(t, "1000\n", stdout)

}

func TestTLSCert(t *testing.T) {
	// start clean
	utils.Exec(t, `sudo rm -rf `+snapData)

	t.Cleanup(func() {
		utils.Exec(t, `sudo rm -rf `+snapData)
	})

	// Due to confinement issues when running this test, we write the files to SNAP_DATA
	utils.Exec(t, `sudo mkdir -p `+snapData)

	// Read the API Gateway token
	kongAminJwt, _ := utils.Exec(t, "sudo cat "+kongAminJwtFile)
	require.NotEmpty(t, kongAminJwt)

	// Add the certificate, using Kong Admin JWT to authenticate
	serverKeyFile, serverCertFile := certGenerator()
	utils.Exec(t, `sudo edgexfoundry.secrets-config proxy tls --incert `+serverCertFile+` --inkey `+serverKeyFile+` --admin_api_jwt `+kongAminJwt)

	code, _ := utils.Exec(t, `curl --show-error --silent --include \
		--output /dev/null --write-out "%{http_code}" \
		--cacert `+serverCertFile+` \
		-X GET 'https://localhost:8443/core-data/api/v2/ping?' \
		-H "Authorization: Bearer $TOKEN"`)
	require.Equal(t, "200\n", code)
}

func certGenerator() (string, string) {
	const (
		caKeyFile  = snapData + "/ca.key"
		caCertFile = snapData + "/ca.crt"

		serverCsrFile  = snapData + "/server.csr"
		serverKeyFile  = snapData + "/server.key"
		serverCertFile = snapData + "/server.crt"
	)

	// Generate the Certificate Authority (CA) Private Key
	utils.Exec(nil, `sudo openssl ecparam -name prime256v1 -genkey -noout -out `+caKeyFile)
	// Generate the Certificate Authority Certificate
	utils.Exec(nil, `sudo openssl req -new -x509 -sha256 -key `+caKeyFile+` -out `+caCertFile+` -subj "/CN=localhost"`)

	// Generate the Server Certificate Private Key
	utils.Exec(nil, `sudo openssl ecparam -name prime256v1 -genkey -noout -out `+serverKeyFile)
	// Generate the Server Certificate Signing Request
	utils.Exec(nil, `sudo openssl req -new -sha256 -key `+serverKeyFile+` -out `+serverCsrFile+` -subj "/CN=localhost"`)
	// Generate the Server Certificate
	utils.Exec(nil, `sudo openssl x509 -req -in `+serverCsrFile+` -CA `+caCertFile+` -CAkey `+caKeyFile+` -CAcreateserial -out `+serverCertFile+` -days 1000 -sha256`)

	return serverKeyFile, serverCertFile
}
