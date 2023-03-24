# Ubuntu Core EdgeX Image Testing

This folder contains scripts for creating an Ubuntu Core OS image, which includes the necessary EdgeX components. The image is designed for testing purposes, and the scripts automate the process of building, booting, and testing the image. 

For details of creating a Ubuntu Core image with EdgeX, please refer to the documentation provided [here](https://docs.edgexfoundry.org/2.3/examples/Ch-OSImageWithEdgeX/#a-create-an-image-with-edgex-components).

### Requirements
Before running the tests, you need to install the following packages:
```
snap install --classic snapcraft
snap install --classic ubuntu-image
snap install lxd
snap install yq
snap install curl
sudo apt update
sudo apt install qemu-kvm
kvm-ok
sudo apt install ovmf
```
Please refer to [Testing Ubuntu Core with QEMU](https://ubuntu.com/core/docs/testing-with-qemu) to learn more about it.

Additionally, you need to log into your snap account by running:
```
snap login
snap keys
```
Note: Follow the instructions [here](https://snapcraft.io/docs/creating-your-developer-account) to create a developer account, if you don't already have one.

After logging in, create a key named `ubuntu-core-edgex-image-test` which will be used later by running:
```
$ snap create-key ubuntu-core-edgex-image-test
# set passphrase
$ snapcraft register-key ubuntu-core-edgex-image-test
```


### Build, Boot and Test the Image
To build the Ubuntu Core image with EdgeX components, run the following command:
```
KEY_NAME=ubuntu-core-edgex-image-test ./build.sh
```
To boot the image in an emulator:
```
./run.sh
```
**After the complete installation, use the wizard to setup the networking and user account.**

Once you see the SSH command and the IP address, open a new terminal to test the image by running:

```
$ USER=<your-username> ./test.sh
# {"apiVersion":"v2","timestamp":"Wed Mar 22 18:05:07 UTC 2023","serviceName":"core-data"}
# {"apiVersion":"v2","timestamp":"Wed Mar 22 18:05:07 UTC 2023","serviceName":"core-metadata"}
# {"apiVersion":"v2","timestamp":"Wed Mar 22 18:05:07 UTC 2023","serviceName":"core-command"}
# {
#         "security": false
# }

# Service                              Startup  Current  Notes
# edgex-device-virtual.device-virtual  enabled  active   -

# 2023-03-22T17:51:27Z edgex-device-virtual.device-virtual[3874]: level=INFO ts=2023-03-22T17:51:27.022824501Z app=device-virtual source=variables.go:377 msg="Variables override of 'Service.StartupMsg' by environment variable: SERVICE_STARTUPMSG=Startup message from gadget!"
```
