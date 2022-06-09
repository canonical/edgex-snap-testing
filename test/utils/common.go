package utils

import (
	"testing"
)

func TestCommon(t *testing.T, params Params) {
	TestConfig(t, params.Snap, params.Config)
	TestNet(t, params.Snap, params.Net)
	TestPackaging(t, params.Snap, params.Packaging)
}
