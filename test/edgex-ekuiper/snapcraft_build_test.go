package test

import (
	"edgex-snap-testing/env"
	"edgex-snap-testing/utils"
	"testing"
)

func TestSnapBuild(t *testing.T) {

	if env.WorkDir != "" {
		utils.SnapBuild(t, env.WorkDir)
	}
}
