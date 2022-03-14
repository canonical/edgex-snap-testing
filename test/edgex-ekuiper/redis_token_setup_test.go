package test

import (
	"edgex-snap-testing/utils"
	"testing"
)

func setupSubtestRedisToeknSetUp(t *testing.T) {
	t.Log("[SUBTEST SETUP]")
	utils.Exec(t,
		"sudo snap restart edgex-ekuiper.kuiper")
}
func TestRedisToeknSetUp(t *testing.T) {
	//TODO: how to validate that redis token works, since security-on already works?
}
