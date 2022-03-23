# EdgeX Snap Tests
Test scripts for EdgeX Foundry snaps.

## Test manually
Test all:
```bash
go test -v ./test/suites/...
```

Test one, e.g.:
```bash
go test -v ./test/suites/device-mqtt
```

Test the testing utils:
```bash
go test -v ./test/utils/...
```

### Override behavior
Use environment variables, as defined in [env/env.go](./env/env.go)

## Test using Github Actions
This project includes two Github Actions that can be used in workflows to test snaps:
* [build](./build): Checkout code, build the snap, upload snap as build artifact
* [test](./test): Download snap, run smoke tests

A workflow that uses both the actions from `main` branch may look as follows:

`.github/workflows/test-snap.yml`
```yaml
name: Snap Testing

on:
  pull_request:
    branches: [ main ]
  # allow manual trigger
  workflow_dispatch:

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - name: Build and upload snap
        id: build
        uses: canonical/edgex-snap-testing/build@v2
    outputs:
      snap: ${{steps.build.outputs.snap}}

  test:
    needs: build
    runs-on: ubuntu-latest
    steps:
      - name: Download and test snap
        uses: canonical/edgex-snap-testing/test@v2
        with:
          name: device-mqtt
          snap: ${{needs.build.outputs.snap}}
```