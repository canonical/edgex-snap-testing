package test

import (
	"edgex-snap-testing/test/common"
	"testing"
)

func TestVersion(t *testing.T) {
	common.TestVersion(t, thisSnap)
}
