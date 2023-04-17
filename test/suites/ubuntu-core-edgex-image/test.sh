#!/bin/bash -ex

if [[ -z "$USER" ]]; then
  echo "Required input 'USER' is unset. Exiting..."
  exit 1
fi

SSH_USER="$USER"
SSH_PORT="8022"
remote_call() {
  ssh "$SSH_USER@localhost" -p $SSH_PORT "$@"
}

# Install curl on emulator
remote_call "snap install curl"

# Check status of core services inside of the emulator
ports=(59880 59881 59882)

for port in "${ports[@]}"
do
  remote_call "curl -s http://localhost:$port/api/v2/ping"
done

# Verify that the security is avaliable as a snap option of edgexfoundry within the emulator
remote_call "snap get edgexfoundry security-secret-store -d"

# Check the status of the device-virtual service within the emulator
remote_call "snap services edgex-device-virtual"

# List snaps and check edgex-config-provider-example is in the list
remote_call 'snap list'

# Verify that Device Virtual only has one profile, as configured in the config provider
remote_call 'curl --silent http://localhost:59881/api/v2/deviceprofile/all' | jq '.totalCount'

# Verify that Device Virtual has the startup message set from the provider
remote_call 'snap logs -n=all edgex-device-virtual | grep "CONFIG BY EXAMPLE PROVIDER"'

