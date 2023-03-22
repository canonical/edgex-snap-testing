#!/bin/bash -e

ssh "$USER@localhost" -p 8022 "snap install curl"

# check status of core services inside of the emulator
ports=(59880 59881 59882)

for port in "${ports[@]}"
do
  ssh "$USER@localhost" -p 8022 "curl -s http://localhost:$port/api/v2/ping; printf '\n'"
done

# verify that the security is avaliable as a snap option of edgexfoundry within the emulator
ssh "$USER@localhost" -p 8022 "snap get edgexfoundry security -d; printf '\n'"

# check the status of the device-virtual service within the emulator
ssh "$USER@localhost" -p 8022 "snap services edgex-device-virtual; printf '\n'"

# verify that device-virtual has the startup message set from the gadget within the emulator
ssh "$USER@localhost" -p 8022 'snap logs -n=all edgex-device-virtual | grep "Startup message"; printf '\n''

