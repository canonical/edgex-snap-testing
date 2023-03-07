#!/bin/bash -e

# check status of core and device services inside of the emulator
ports=(59880 59881 59882 59900)

for port in "${ports[@]}"
do
  ssh "$USER@localhost" -p 8022 "curl -s http://localhost:$port/api/v2/ping; printf '\n'"
done

# # make core services listen to all interfaces
# ssh "$USER@localhost" -p 8022 "snap set edgexfoundry app-options=true; printf '\n'"

# services=(core-data core-metadata core-command)

# for service in "${services[@]}"
# do
#   ssh "$USER@localhost" -p 8022 "snap set edgexfoundry apps.$service.config.service-serverbindaddr='0.0.0.0'; snap restart edgexfoundry.$service"
# done

# sleep 5


# # check status of core and device services outside of the emulator
# ports=(59880 59881 59882)
# for port in "${ports[@]}"
# do
#   curl http://localhost:$port/api/v2/ping; printf '\n'
# done

