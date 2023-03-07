# Ubuntu Core EdgeX Image Testing

This folder contains scripts for creating an Ubuntu Core OS image, which includes the necessary EdgeX components. For details of how to creat a Ubuntu Core image with EdgeX, please refer to this [documentation](https://docs.edgexfoundry.org/2.3/examples/Ch-OSImageWithEdgeX/#a-create-an-image-with-edgex-components).

### Requirement
To run this test, you will need to install the following snap packages: `snapcraft`, `ubuntu-image`, `lxd`, and `yq`. 
To install them, run the following commands:
```
snap install --classic snapcraft
snap install --classic ubuntu-image
snap install lxd
snap install yq
```

Additionally, you need to log into your snap account by running:
```
snap login
snap keys
```
After logging in, create a key for the `ubuntu-core-edgex-image-test` by running:
```
snap create-key ubuntu-core-edgex-image-test
snapcraft register-key ubuntu-core-edgex-image-test
```

### Run the test
```
go run main.go
```
