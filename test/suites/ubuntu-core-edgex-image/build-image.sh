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

# test ubuntu core with QEMU
sudo qemu-system-x86_64 \
 -smp 4 \
 -m 4096 \
 -drive file=/usr/share/OVMF/OVMF_CODE.fd,if=pflash,format=raw,unit=0,readonly=on \
 -drive file=pc.img,cache=none,format=raw,id=disk1,if=none \
 -device virtio-blk-pci,drive=disk1,bootindex=1 \
 -machine accel=kvm \
 -serial mon:stdio \
 -net nic,model=virtio \
 -net user,hostfwd=tcp::8022-:22

