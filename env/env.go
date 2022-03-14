package env

import "os"

const (
	// environment variables
	// used to override defaults
	channel = "CHANNEL"
	snap    = "SNAP"
	workDir = "WORKDIR"
)

var (
	// global defaults
	Channel = "latest/edge"
	Snap    = ""
	WorkDir = ""
)

func init() {
	if v := os.Getenv(channel); v != "" {
		Channel = v
	}

	if v := os.Getenv(snap); v != "" {
		Snap = v
	}
	if v := os.Getenv(workDir); v != "" {
		WorkDir = v
	}
}
