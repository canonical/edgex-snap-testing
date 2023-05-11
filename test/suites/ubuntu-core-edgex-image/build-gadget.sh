#!/bin/bash -ex

# Remove the pc-gadget directory if it already exists
rm -rf pc-gadget
git clone https://github.com/snapcore/pc-gadget.git --branch=22 --depth 1 

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
    "security": false,
    "autostart": true
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
} ' gadget.yaml


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
# Add kernel options for extended logging and a debugging shell
cp -r ../kernel-options/ .
yq e -i '.parts += { 
    "kernel-options": 
      { 
        "source": "kernel-options/",
        "plugin": "dump" 
      }
    }' snapcraft.yaml


# Build gadget snap
snapcraft

