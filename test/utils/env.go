package utils

import (
	"os"
	"strconv"
)

const (
	// environment variables
	// used to override defaults
	platformChannel = "PLATFORM_CHANNEL"
	serviceChannel  = "SERVICE_CHANNEL"
	localSnap       = "LOCAL_SNAP"
	fullConfigTest  = "FULL_CONFIG_TEST"
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
		FullConfigTest, _ = strconv.ParseBool(v)
	}
}
