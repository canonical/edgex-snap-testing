#!/bin/bash -ex

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
 -net user,hostfwd=tcp::8022-:22,hostfwd=tcp::59880-:59880,hostfwd=tcp::8443-:8443


