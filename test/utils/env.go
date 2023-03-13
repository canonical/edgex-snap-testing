package utils

import (
	"os"
	"strconv"
)

const (
	// environment variables
	// used to override defaults
	platformChannelEnv     = "PLATFORM_CHANNEL"      // channel/revision of the edgexfoundry snap (has default)
	serviceChannelEnv      = "SERVICE_CHANNEL"       // channel/revision of the service snap (has default)
	localSnapEnv           = "LOCAL_SNAP"            // path to local snap to be tested instead of downloading from a channel
	fullConfigTestEnv      = "FULL_CONFIG_TEST"      // toggle full config tests (has default)
	skipTeardownRemovalEnv = "SKIP_TEARDOWN_REMOVAL" // skip the removal of snaps during teardown
)

var (
	// global defaults
	PlatformChannel     = "latest/edge"
	ServiceChannel      = "latest/edge"
	LocalSnapPath       = ""
	FullConfigTest      = false
	SkipTeardownRemoval = false
)

func init() {
	if v := os.Getenv(platformChannelEnv); v != "" {
		PlatformChannel = v
	}

	if v := os.Getenv(serviceChannelEnv); v != "" {
		ServiceChannel = v
	}

	if v := os.Getenv(localSnapEnv); v != "" {
		LocalSnapPath = v
	}

	if v := os.Getenv(fullConfigTestEnv); v != "" {
		var err error
		FullConfigTest, err = strconv.ParseBool(v)
		if err != nil {
			panic(err)
		}
	}

	if v := os.Getenv(skipTeardownRemovalEnv); v != "" {
		var err error
		SkipTeardownRemoval, err = strconv.ParseBool(v)
		if err != nil {
			panic(err)
		}
	}
}
