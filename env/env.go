package env

import "os"

const (
	// environment variables
	// used to override defaults
	channel = "CHANNEL"
	snap    = "SNAP"
)

var (
	// global defaults
	Channel = "latest/edge"
	Snap    = ""
)

func init() {
	if v := os.Getenv(channel); v != "" {
		Channel = v
	}

	if v := os.Getenv(snap); v != "" {
		Snap = v
	}
}
