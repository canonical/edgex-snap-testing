#!/bin/bash -ex

username=test
password=$(sudo edgexfoundry.secrets-config proxy adduser --user ${username} --useRootToken | jq -r '.password')

vault_token=$(curl -sS "http://localhost:8200/v1/auth/userpass/login/${username}" -d "{\"password\":\"${password}\"}" | jq -re '.auth.client_token')

id_token=$(curl -sS -H "Authorization: Bearer ${vault_token}" "http://localhost:8200/v1/identity/oidc/token/${username}" | jq -re '.data.token')

echo $id_token
