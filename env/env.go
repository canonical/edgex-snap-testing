package env

import "os"

const (
	// environment variables
	// used to override defaults
	channel             = "CHANNEL"
	snapcraftProjectDir = "SNAPCRAFT_PROJECT_DIR"
)

var (
	// global defaults
	Channel             = "latest/edge"
	SnapcraftProjectDir = ""
)

func init() {
	if v := os.Getenv(channel); v != "" {
		Channel = v
	}

	if v := os.Getenv(snapcraftProjectDir); v != "" {
		SnapcraftProjectDir = v
	}
}
