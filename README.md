# EdgeX Snap Tests
Test scripts for EdgeX Foundry snaps.

## Usage
Test all:
```bash
go test -v ./tests/...
```

Test one, e.g.:
```bash
go test -v ./tests/device-mqtt
```

Test the testing utils:
```bash
go test -v ./utils
```

### Override behavior
Use environment variables, as defined in [env/env.go](./env/env.go)