package env

import "os"

const (
	// environment variables
	// used to override defaults
	channel        = "CHANNEL"
	ekuiperchannel = "EKUIPERCHANNEL"
	snap           = "SNAP"
)

var (
	// global defaults
	Channel        = "latest/edge"
	EKuiperChannel = "1/edge"
	Snap           = ""
)

func init() {
	if v := os.Getenv(channel); v != "" {
		Channel = v
	}

	if v := os.Getenv(ekuiperchannel); v != "" {
		Snap = v
	}

	if v := os.Getenv(snap); v != "" {
		Snap = v
	}

}
