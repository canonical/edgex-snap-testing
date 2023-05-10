#!/bin/bash -e

SERVER_CERT_FILE=server.cert
SERVER_KEY_FILE=server.key
SERVER_CSR_FILE=server.csr
CA_CERT_FILE=ca.cert
CA_KEY_FILE=ca.key

# Generate the Certificate Authority (CA) Private Key
openssl ecparam -name prime256v1 -genkey -noout -out $CA_KEY_FILE
# Generate the Certificate Authority Certificate
openssl req -new -x509 -sha256 -key $CA_KEY_FILE -out $CA_CERT_FILE -subj "/CN=snap-testing-ca"
# Generate the Server Certificate Private Key
openssl ecparam -name prime256v1 -genkey -noout -out $SERVER_KEY_FILE
# Generate the Server Certificate Signing Request
openssl req -new -sha256 -key $SERVER_KEY_FILE -out $SERVER_CSR_FILE -subj "/CN=localhost"
# Generate the Server Certificate
openssl x509 -req -in $SERVER_CSR_FILE -CA $CA_CERT_FILE -CAkey $CA_KEY_FILE -CAcreateserial -out $SERVER_CERT_FILE -days 1000 -sha256

# copy the files to a directory that the snap has permission to see
sudo cp server.cert server.key /var/snap/edgexfoundry/common

sudo edgexfoundry.secrets-config proxy tls \
    --inCert /var/snap/edgexfoundry/common/server.cert \
    --inKey /var/snap/edgexfoundry/common/server.key \
    --targetFolder /var/snap/edgexfoundry/current/nginx

sudo rm server.cert server.key

sudo snap restart --reload edgexfoundry.nginx
