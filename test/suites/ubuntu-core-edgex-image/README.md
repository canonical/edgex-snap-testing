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
To build the gadget snap, run the following command:
```
./build-gadget.sh
```
To build the Ubuntu Core image:
```
KEY_NAME=ubuntu-core-edgex-image-test ./build-image.sh
```
To boot the Ubuntu Core image with EdgeX components in an emulator:
```
./run.sh
```
**After the complete installation, use the wizard to setup the networking and user account.**

Once you see the SSH command and the IP address, open a new terminal to test the image by running:

```
$ USER=<your-username> ./test.sh
+ SSH_USER=<your-username>
+ SSH_PORT=8022
+ remote_call 'snap install curl'
+ ssh <your-username>@localhost -p 8022 'snap install curl'
curl 8.0.1 from Wouter van Bommel (woutervb) installed
+ ports=(59880 59881 59882)
+ for port in "${ports[@]}"
+ remote_call 'curl -s http://localhost:59880/api/v3/ping'
+ ssh <your-username>@localhost -p 8022 'curl -s http://localhost:59880/api/v3/ping'
{"apiVersion":"v3","timestamp":"Tue Apr  4 07:33:14 UTC 2023","serviceName":"core-data"}+ for port in "${ports[@]}"
+ remote_call 'curl -s http://localhost:59881/api/v3/ping'
+ ssh <your-username>@localhost -p 8022 'curl -s http://localhost:59881/api/v3/ping'
{"apiVersion":"v3","timestamp":"Tue Apr  4 07:33:14 UTC 2023","serviceName":"core-metadata"}+ for port in "${ports[@]}"
+ remote_call 'curl -s http://localhost:59882/api/v3/ping'
+ ssh <your-username>@localhost -p 8022 'curl -s http://localhost:59882/api/v3/ping'
{"apiVersion":"v3","timestamp":"Tue Apr  4 07:33:14 UTC 2023","serviceName":"core-command"}+ remote_call 'snap get edgexfoundry security-secret-store -d'
+ remote_call 'snap get edgexfoundry security -d'
+ ssh <your-username>@localhost -p 8022 'snap get edgexfoundry security -d'
{
        "security": false
}
+ remote_call 'snap services edgex-device-virtual'
+ ssh <your-username>@localhost -p 8022 'snap services edgex-device-virtual'
Service                              Startup  Current  Notes
edgex-device-virtual.device-virtual  enabled  active   -
+ remote_call 'snap list'
+ ssh <your-username>@localhost -p 8022 'snap list'
Name                           Version         Rev    Tracking       Publisher    Notes
core20                         20230404        1879   latest/stable  canonical**  base
core22                         20230404        617    latest/stable  canonical**  base
curl                           8.0.1           1679   latest/stable  woutervb     -
edgex-config-provider-example  2.3             6      latest/stable  farshidtz    -
edgex-device-virtual           3.0.0-dev.45    645    latest/edge    canonical**  -
edgexfoundry                   3.0.0-dev.155   4428   latest/edge    canonical**  -
pc                             22-0.3          x1     -              -            gadget
pc-kernel                      5.15.0-71.78.1  1281   22/stable      canonical**  kernel
snapd                          2.59.2          19122  latest/stable  canonical**  snapd
+ remote_call 'curl --silent http://localhost:59881/api/v3/deviceprofile/all'
+ ssh <your-username>@localhost -p 8022 'curl --silent http://localhost:59881/api/v3/deviceprofile/all'
+ jq .totalCount
0
+ remote_call 'snap logs -n=all edgex-device-virtual | grep "CONFIG BY EXAMPLE PROVIDER"'
+ ssh <your-username>@localhost -p 8022 'snap logs -n=all edgex-device-virtual | grep "CONFIG BY EXAMPLE PROVIDER"'

```
