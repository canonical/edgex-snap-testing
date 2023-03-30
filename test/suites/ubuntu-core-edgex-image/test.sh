#!/bin/bash -e

SSH_USER="$USER"
SSH_PORT="8022"

printf "# Install curl on host:\n"
snap install curl; echo

printf "# Check status of core services inside of the emulator:\n"
ports=(59880 59881 59882)

for port in "${ports[@]}"
do
  ssh "$SSH_USER@localhost" -p $SSH_PORT "curl -s http://localhost:$port/api/v2/ping"; echo
done
echo

printf "# Verify that the security is avaliable as a snap option of edgexfoundry within the emulator:\n"
ssh "$SSH_USER@localhost" -p $SSH_PORT "snap get edgexfoundry security -d"; echo

printf "# Check the status of the device-virtual service within the emulator:\n"
ssh "$SSH_USER@localhost" -p $SSH_PORT "snap services edgex-device-virtual"; echo


printf "# Access the service endpoints via API Gateway outside of the emulator:\n"
curl --insecure --show-err https://localhost:8443/core-data/api/v2/ping; echo
echo

printf "# List snaps and check edgex-config-provider-example is in the list:\n"
ssh "$SSH_USER@localhost" -p $SSH_PORT 'snap list'; echo

printf "# Verify that Device Virtual only has one profile, as configured in the config provider:\n"
ssh "$SSH_USER@localhost" -p $SSH_PORT 'curl --silent http://localhost:59881/api/v2/deviceprofile/all' | jq '.totalCount'; echo 

printf "# Verify that Device Virtual has the startup message set from the provider:\n"
ssh "$SSH_USER@localhost" -p $SSH_PORT 'snap logs -n=all edgex-device-virtual | grep "CONFIG BY EXAMPLE PROVIDER"'; echo

printf "# Query the metadata of Device Virtual from your host machine:\n"
curl --insecure --silent --show-err https://localhost:8443/core-data/api/v2/reading/all; echo
