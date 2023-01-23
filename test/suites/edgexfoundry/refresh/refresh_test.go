package test

import (
	"edgex-snap-testing/test/utils"
	"testing"
)

const (
	platformSnap = "edgexfoundry"
)

func TestCommon(t *testing.T) {
	utils.TestRefresh(t, platformSnap)
}
