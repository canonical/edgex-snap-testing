#!/bin/bash -ex

rm -rf pc-gadget
git clone https://github.com/snapcore/pc-gadget.git --branch=22

# build gadget snap
cd pc-gadget
yq e '(.volumes.pc.structure[] | select(.name=="ubuntu-seed") | .size)="1500M"' gadget.yaml -i
snapcraft

# configure model assertion
cd ../
DEVELOPER_ID=$(snapcraft whoami | grep 'id:' | awk '{print $2}')
TIMESTAMP=$(date -Iseconds --utc)
yq e -i ".authority-id = \"$DEVELOPER_ID\"" model.yaml
yq e -i ".brand-id = \"$DEVELOPER_ID\"" model.yaml
yq e -i ".timestamp = \"$TIMESTAMP\"" model.yaml

# sign the model assertion
yq eval model.yaml -o=json | snap sign -k $SNAP_KEY > model.signed.yaml

# check the signed model
cat model.signed.yaml
