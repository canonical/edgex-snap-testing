#!/bin/bash -ex

if [[ -z "$KEY_NAME" ]]; then
  echo "Required input 'KEY_NAME' is unset. Exiting..."
  exit 1
fi

# Configure model assertion
DEVELOPER_ID=$(snapcraft whoami | grep 'id:' | awk '{print $2}')
TIMESTAMP=$(date -Iseconds --utc)
yq e -i ".authority-id = \"$DEVELOPER_ID\"" model.yaml
yq e -i ".brand-id = \"$DEVELOPER_ID\"" model.yaml
yq e -i ".timestamp = \"$TIMESTAMP\"" model.yaml

# Sign the model assertion
yq eval model.yaml -o=json | snap sign -k $KEY_NAME > model.signed.yaml

# Check the signed model
cat model.signed.yaml

# clean and build the image
sudo rm -f pc.img
ubuntu-image snap model.signed.yaml --validation=enforce --snap pc-gadget/pc_*_amd64.snap

# check the image file
file pc.img

