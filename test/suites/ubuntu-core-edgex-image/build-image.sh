#!/bin/bash -ex

if [[ -z "$KEY_NAME" ]]; then
  echo "Required input 'KEY_NAME' is unset. Exiting..."
  exit 1
fi

MODEL=model-test.yaml
SIGNED_MODEL=model-test.signed.yaml
cp model.yaml $MODEL

# Configure model assertion
DEVELOPER_ID=$(snapcraft whoami | grep 'id:' | awk '{print $2}')
TIMESTAMP=$(date -Iseconds --utc)
yq e -i ".authority-id = \"$DEVELOPER_ID\"" $MODEL
yq e -i ".brand-id = \"$DEVELOPER_ID\"" $MODEL
yq e -i ".timestamp = \"$TIMESTAMP\"" $MODEL

# Sign the model assertion
yq eval $MODEL -o=json | snap sign -k $KEY_NAME > $SIGNED_MODEL

# Check the signed model
cat $SIGNED_MODEL

# clean and build the image
ubuntu-image snap $SIGNED_MODEL --validation=enforce --snap pc-gadget/pc_*_amd64.snap

# check the image file
file pc.img

