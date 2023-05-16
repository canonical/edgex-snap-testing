#!/bin/bash -ex

username=test
password=$(sudo edgexfoundry.secrets-config proxy adduser --user ${username} --useRootToken | jq -r '.password')
if [[ -z "${password}" ]]; then
  >&2 echo "Error getting a password for user"
  exit 1
fi

vault_token=$(curl -sS "http://localhost:8200/v1/auth/userpass/login/${username}" -d "{\"password\":\"${password}\"}" | jq -re '.auth.client_token')
if [[ -z "${vault_token}" ]]; then
  >&2 echo "Error retrieving Vault token for user ${username}"
  exit 1
fi

id_token=$(curl -sS -H "Authorization: Bearer ${vault_token}" "http://localhost:8200/v1/identity/oidc/token/${username}" | jq -re '.data.token')
if [[ -z "${id_token}" ]]; then
  >&2 echo "Error retrieving ID token for user ${username}"
  exit 1
fi

echo $id_token
