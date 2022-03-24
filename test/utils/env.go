package utils

import "os"

const (
	// environment variables
	// used to override defaults
	platformChannel = "PLATFORM_CHANNEL"
	serviceChannel  = "SERVICE_CHANNEL"
	localSnap       = "LOCAL_SNAP"
)

var (
	// global defaults
	PlatformChannel = "latest/edge"
	ServiceChannel  = "latest/edge"
	LocalSnap       = ""
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
}
