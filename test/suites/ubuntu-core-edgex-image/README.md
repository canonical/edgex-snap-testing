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
snap create-key ubuntu-core-edgex-image-test
# set passphrase
snapcraft register-key ubuntu-core-edgex-image-test
```


### Build, Boot and Test the Image
To build the Ubuntu Core image with EdgeX components, run the following command:
```
SNAP_KEY=ubuntu-core-edgex-image-test ./build.sh
```
To boot the image in an emulator:
```
./run.sh
```
Once the networking and user account setup is complete, open a new terminal to test the image by running:
```
USER=<your-username> ./test.sh
# {"apiVersion":"v2","timestamp":"Tue Mar  7 18:18:14 UTC 2023","serviceName":"core-data"}
# {"apiVersion":"v2","timestamp":"Tue Mar  7 18:18:15 UTC 2023","serviceName":"core-metadata"}
# {"apiVersion":"v2","timestamp":"Tue Mar  7 18:18:15 UTC 2023","serviceName":"core-command"}
```
