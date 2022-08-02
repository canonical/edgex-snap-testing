package utils

import (
	"os"
	"strconv"
)

const (
	// environment variables
	// used to override defaults
	platformChannelEnv = "PLATFORM_CHANNEL" // edgexfoundry channel when testing other snaps (has default)
	serviceChannelEnv  = "SERVICE_CHANNEL"  // channel of the snap to be tested (has default)
	localSnapEnv       = "LOCAL_SNAP"       // path to local snap to be tested instead of downloading from a channel
	fullConfigTestEnv  = "FULL_CONFIG_TEST" // toggle full config tests (has default)
)

var (
	// global defaults
	PlatformChannel = "latest/edge"
	ServiceChannel  = "latest/edge"
	LocalSnapPath   = ""
	FullConfigTest  = false
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
}
