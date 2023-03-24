#!/bin/bash -ex

rm -rf pc-gadget
git clone https://github.com/snapcore/pc-gadget.git --branch=22

# build gadget snap
cd pc-gadget
# extend the size of disk partitions to have sufficient capacity for EdgeX snaps
yq e '(.volumes.pc.structure[] | select(.name=="ubuntu-seed") | .size)="1500M"' gadget.yaml -i

# setup default options for snaps
# AZGf0KNnh8aqdkbGATNuRuxnt1GNRKkV (edgexfoundry snap)
# AmKuVTOfsN0uEKsyJG34M8CaMfnIqxc0 (edgex-device-virtual snap) 
yq e '.defaults += {
  "AZGf0KNnh8aqdkbGATNuRuxnt1GNRKkV": {
    "app-options": true,
    "security": false
  },
  "AmKuVTOfsN0uEKsyJG34M8CaMfnIqxc0": {
    "autostart": true,
    "app-options": true,
    "apps": {
      "device-virtual": {
        "config": {
          "service-startupmsg": "Startup message from gadget!",
          "edgex-security-secret-store": false
        }
      }
    }
  }
}' gadget.yaml -i

snapcraft

# configure model assertion
cd ../
DEVELOPER_ID=$(snapcraft whoami | grep 'id:' | awk '{print $2}')
TIMESTAMP=$(date -Iseconds --utc)
yq e -i ".authority-id = \"$DEVELOPER_ID\"" model.yaml
yq e -i ".brand-id = \"$DEVELOPER_ID\"" model.yaml
yq e -i ".timestamp = \"$TIMESTAMP\"" model.yaml

# sign the model assertion
yq eval model.yaml -o=json | snap sign -k $KEY_NAME > model.signed.yaml

# check the signed model
cat model.signed.yaml
