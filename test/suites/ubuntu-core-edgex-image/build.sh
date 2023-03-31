#!/bin/bash -ex

# Remove the pc-gadget directory if it already exists
rm -rf pc-gadget
git clone https://github.com/snapcore/pc-gadget.git --branch=22

# Build gadget snap
cd pc-gadget

# Extend the size of disk partitions to have sufficient capacity for EdgeX snaps
yq e -i '(.volumes.pc.structure[] | select(.name=="ubuntu-seed") | .size)="1500M"' gadget.yaml

# Set up default options for snaps
# AZGf0KNnh8aqdkbGATNuRuxnt1GNRKkV (edgexfoundry snap)
# AmKuVTOfsN0uEKsyJG34M8CaMfnIqxc0 (edgex-device-virtual snap) 
yq e -i '.defaults += {
  "AZGf0KNnh8aqdkbGATNuRuxnt1GNRKkV": {
    "app-options": true,
    "security-secret-store": "off"
  },
  "AmKuVTOfsN0uEKsyJG34M8CaMfnIqxc0": {
    "autostart": true,
    "apps": {
      "device-virtual": {
        "config": {
          "edgex-security-secret-store": false
        }
      }
    }
  }
}' gadget.yaml

# Connect edgex-device-virtual's plug (consumer) to 
# edgex-config-provider-example's slot (provider) 
# to override the default configuration files.
yq e -i '.connections += [
          {
            "plug": "AmKuVTOfsN0uEKsyJG34M8CaMfnIqxc0:device-virtual-config", 
            "slot": "WWPGZGi1bImphPwrRfw46aP7YMyZYl6w:device-virtual-config"
          }
        ]
      ' gadget.yaml

snapcraft

# Configure model assertion
cd ../
DEVELOPER_ID=$(snapcraft whoami | grep 'id:' | awk '{print $2}')
TIMESTAMP=$(date -Iseconds --utc)
yq e -i ".authority-id = \"$DEVELOPER_ID\"" model.yaml
yq e -i ".brand-id = \"$DEVELOPER_ID\"" model.yaml
yq e -i ".timestamp = \"$TIMESTAMP\"" model.yaml

# Sign the model assertion
yq eval model.yaml -o=json | snap sign -k $KEY_NAME > model.signed.yaml

# Check the signed model
cat model.signed.yaml
