#!/bin/bash -e
sudo apt-get update
sudo apt install qemu-kvm
sudo apt install ovmf

# start lxd for snapcraft usage
snap start lxd
