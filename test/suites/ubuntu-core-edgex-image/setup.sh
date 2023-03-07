#!/bin/bash -e

# source snap.env

# only for debugging
snap remove --purge snapcraft
snap remove --purge multipass
snap remove --purge yq
snap remove --purge ubuntu-image
snap remove --purge qemu-kvm
snap remove --purge ovmf


snap install snapcraft --classic
snap install multipass
snap install yq
snap install ubuntu-image --classic

snap start multipass.multipassd
multipass start

snap login
snap keys

# snap create-key edgex-demo-test
# snapcraft register-key edgex-demo-test

sudo apt install qemu-kvm
kvm-ok
sudo apt install ovmf
