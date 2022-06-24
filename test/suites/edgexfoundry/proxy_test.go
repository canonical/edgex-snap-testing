package test

import (
	"edgex-snap-testing/test/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	EDGEXFOUNDRY_SNAP_DATA = "/var/snap/edgexfoundry/current/temp"
	PUBLIC_KEY             = EDGEXFOUNDRY_SNAP_DATA + "/public.pem"
	PRIVATE_KEY            = EDGEXFOUNDRY_SNAP_DATA + "/private.pem"
	CA_KEY_FILE            = EDGEXFOUNDRY_SNAP_DATA + "/ca.key"
	CA_CERT_FILE           = EDGEXFOUNDRY_SNAP_DATA + "/ca.crt"
	SERVER_KEY_FILE        = EDGEXFOUNDRY_SNAP_DATA + "/server.key"
	SERVER_CSR_FILE        = EDGEXFOUNDRY_SNAP_DATA + "/server.csr"
	SERVER_CERT_FILE       = EDGEXFOUNDRY_SNAP_DATA + "/server.crt"
	KONG_ADMIN_JWT_FILE    = EDGEXFOUNDRY_SNAP_DATA + "/secrets/security-proxy-setup/kong-admin-jwt"
)


	// start clean
	utils.Exec(t, `rm -rf `+CA_KEY_FILE+` `+CA_CERT_FILE+` `+SERVER_KEY_FILE+` `+SERVER_CSR_FILE+` `+SERVER_CERT_FILE+` `+EDGEXFOUNDRY_SNAP_DATA)

	t.Cleanup(func() {
		utils.Exec(t, `rm -rf `+CA_KEY_FILE+` `+CA_CERT_FILE+` `+SERVER_KEY_FILE+` `+SERVER_CSR_FILE+` `+SERVER_CERT_FILE+` `+EDGEXFOUNDRY_SNAP_DATA)
	})

	// Due to confinement issues when running this test, we write the private key to SNAP_DATA
	// Create private and public keys
	utils.Exec(t, `sudo mkdir -p `+EDGEXFOUNDRY_SNAP_DATA)
	utils.Exec(t, `sudo openssl ecparam -genkey -name prime256v1 -noout -out `+PRIVATE_KEY)
	utils.Exec(t, `sudo openssl ec -in `+PRIVATE_KEY+` -pubout -out `+PUBLIC_KEY)

	// Read the API Gateway token
	KONG_ADMIN_JWT, _ := utils.Exec(t, "sudo cat /var/snap/edgexfoundry/current/secrets/security-proxy-setup/kong-admin-jwt")
	require.NotEmpty(t, KONG_ADMIN_JWT)

	t.Run("add proxy admin user", func(t *testing.T) {
		// Use secrets-config to add a user example with id 1000
		// On success, the above command prints the user id
		stdout, _ := utils.Exec(t, `
			sudo edgexfoundry.secrets-config proxy adduser \
			--token-type jwt --user example --algorithm ES256 \
			--public_key `+PUBLIC_KEY+` --id 1000 -jwt `+KONG_ADMIN_JWT)
		require.Equal(t, "1000\n", stdout)
	})

	t.Run("set custom TLS cert", func(t *testing.T) {
		certGenerator(CA_KEY_FILE, CA_CERT_FILE, SERVER_KEY_FILE, SERVER_CSR_FILE, SERVER_CERT_FILE)

		// Setting security-proxy certificate
		utils.SnapSet(t, "edgexfoundry", "env.security-proxy.tls-certificate", SERVER_CERT_FILE)

		// Setting security-proxy certificate private key
		utils.SnapSet(t, "edgexfoundry", "env.security-proxy.tls-private-key", SERVER_KEY_FILE)

		utils.Exec(t, `cp `+EDGEXFOUNDRY_SNAP_DATA+`/ca.crt `+EDGEXFOUNDRY_SNAP_DATA)

		// Add the certificate, using Kong Admin JWT to authenticate
		utils.Exec(t, `sudo edgexfoundry.secrets-config proxy tls --incert cert.pem --inkey `+PRIVATE_KEY+` --admin_api_jwt `+KONG_ADMIN_JWT)
		code, _ := utils.Exec(t, `curl --show-error --silent --include \
		--output /dev/null --write-out "%{http_code}" \
		--cacert `+CA_KEY_FILE+` \
		-X GET 'https://localhost:8443/core-data/api/v2/ping?' \
		-H "Authorization: Bearer $TOKEN"`)
		require.Equal(t, 200, code)

	})
}

func certGenerator(CA_KEY_FILE, CA_CERT_FILE, SERVER_KEY_FILE, SERVER_CSR_FILE, SERVER_CERT_FILE string) {
	// Generate the Certificate Authority (CA) Private Key
	utils.Exec(nil, `sudo openssl ecparam -name prime256v1 -genkey -noout -out `+CA_KEY_FILE)
	// Generate the Certificate Authority Certificate
	utils.Exec(nil, `sudo openssl req -new -x509 -sha256 -key `+CA_KEY_FILE+` -out `+CA_CERT_FILE+` -subj "/CN=test-ca"`)
	// Generate the Server Certificate Private Key
	utils.Exec(nil, `sudo openssl ecparam -name prime256v1 -genkey -noout -out `+SERVER_KEY_FILE)
	// Generate the Server Certificate Signing Request
	utils.Exec(nil, `sudo openssl req -new -sha256 -key `+SERVER_KEY_FILE+` -out `+SERVER_CSR_FILE+` -subj "/CN=localhost"`)
	// Generate the Server Certificate
	utils.Exec(nil, `sudo openssl x509 -req -in `+SERVER_CSR_FILE+` -CA `+CA_CERT_FILE+` -CAkey `+CA_KEY_FILE+` -CAcreateserial -out `+SERVER_CERT_FILE+` -days 1000 -sha256`)

}
