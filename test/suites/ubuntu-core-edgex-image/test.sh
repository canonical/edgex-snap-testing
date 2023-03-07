ssh mengyiw@localhost -p 8022 'snap start edgex-device-virtual; snap install curl; curl http://localhost:59881/api/v2/device/all'

curl http://localhost:59881/api/v2/ping