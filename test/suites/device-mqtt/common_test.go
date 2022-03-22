package test

import (
	"edgex-snap-testing/test/common-tests"
	"testing"
)

func TestVersion(t *testing.T) {
	common.Version(t, thisSnap)
}
