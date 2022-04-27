package utils

import (
	"os"
	"strconv"
)

const (
	// environment variables
	// used to override defaults
	platformChannel = "PLATFORM_CHANNEL" // edgexfoundry channel when testing other snaps (has default)
	serviceChannel  = "SERVICE_CHANNEL"  // channel of the snap to be tested (has default)
	localSnap       = "LOCAL_SNAP"       // path to local snap to be tested instead of downloading from a channel
	fullConfigTest  = "FULL_CONFIG_TEST" // toggle full config tests (has default)
)

var (
	// global defaults
	PlatformChannel = "latest/edge"
	ServiceChannel  = "latest/edge"
	LocalSnap       = ""
	FullConfigTest  = false
)

func init() {
	if v := os.Getenv(platformChannel); v != "" {
		PlatformChannel = v
	}

	if v := os.Getenv(serviceChannel); v != "" {
		ServiceChannel = v
	}

	if v := os.Getenv(localSnap); v != "" {
		LocalSnap = v
	}

	if v := os.Getenv(fullConfigTest); v != "" {
		var err error
		FullConfigTest, err = strconv.ParseBool(v)
		if err != nil {
			panic(err)
		}
	}
}
