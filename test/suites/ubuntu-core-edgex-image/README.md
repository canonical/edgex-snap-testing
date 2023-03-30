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
# Install curl on host:
curl 8.0.1 from Wouter van Bommel (woutervb) installed

# Check status of core services inside of the emulator:
{"apiVersion":"v2","timestamp":"Thu Mar 30 13:36:10 UTC 2023","serviceName":"core-data"}
{"apiVersion":"v2","timestamp":"Thu Mar 30 13:36:10 UTC 2023","serviceName":"core-metadata"}
{"apiVersion":"v2","timestamp":"Thu Mar 30 13:36:10 UTC 2023","serviceName":"core-command"}

# Verify that the security is avaliable as a snap option of edgexfoundry within the emulator:
{
        "security": false
}

# Check the status of the device-virtual service within the emulator:
Service                              Startup  Current  Notes
edgex-device-virtual.device-virtual  enabled  active   -

# Access the service endpoints via API Gateway outside of the emulator:
{"message":"Unauthorized"}

# List snaps and check edgex-config-provider-example is in the list:
Name                           Version         Rev    Tracking       Publisher    Notes
core20                         20230308        1852   latest/stable  canonical**  base
core22                         20230316        583    latest/stable  canonical**  base
curl                           8.0.1           1579   latest/stable  woutervb     -
edgex-config-provider-example  2.3             6      latest/edge    farshidtz    -
edgex-device-virtual           2.3.0           335    latest/stable  canonical**  -
edgexfoundry                   2.3.0+1         4108   latest/stable  canonical**  -
pc                             22-0.3          x1     -              -            gadget
pc-kernel                      5.15.0-69.76.1  1258   22/stable      canonical**  kernel
snapd                          2.58.3          18596  latest/stable  canonical**  snapd

# Verify that Device Virtual only has one profile, as configured in the config provider:
1

# Verify that Device Virtual has the startup message set from the provider:
2023-03-30T12:47:58Z edgex-device-virtual.device-virtual[3965]: level=INFO ts=2023-03-30T12:47:58.514769034Z app=device-virtual source=message.go:55 msg="CONFIG BY EXAMPLE PROVIDER"

# Query the metadata of Device Virtual from your host machine:
{"message":"Unauthorized"}
```
