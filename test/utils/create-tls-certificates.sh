#!/bin/bash -e

SERVER_CERT_FILE=./tmp/server.cert
SERVER_KEY_FILE=./tmp/server.key
SERVER_CSR_FILE=./tmp/server.csr
CA_CERT_FILE=./tmp/ca.cert
CA_KEY_FILE=./tmp/ca.key

mkdir tmp

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
sudo cp $SERVER_CERT_FILE $SERVER_KEY_FILE /var/snap/edgexfoundry/common

sudo edgexfoundry.secrets-config proxy tls \
    --inCert /var/snap/edgexfoundry/common/server.cert \
    --inKey /var/snap/edgexfoundry/common/server.key \
    --targetFolder /var/snap/edgexfoundry/current/nginx

sudo rm $SERVER_CERT_FILE $SERVER_KEY_FILE

sudo snap restart --reload edgexfoundry.nginx
