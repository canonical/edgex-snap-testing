package test

import (
	"edgex-snap-testing/test/utils"
	"testing"
)

func TestVersion(t *testing.T) {
	utils.RequireSnapSemver(t, appRfidLlrpSnap)
}
