#!/bin/bash -e

ssh "$USER@localhost" -p 8022 "snap install curl"

# check status of core services inside of the emulator
ports=(59880 59881 59882)

for port in "${ports[@]}"
do
  ssh "$USER@localhost" -p 8022 "curl -s http://localhost:$port/api/v2/ping; printf '\n'"
done

