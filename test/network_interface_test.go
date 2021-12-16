package test

import (
	"MQTT-test-project/utils"
)

func (suite *MqttTestSuite) TestNetworkInterface() {
	out := utils.Command("sudo lsof -nPi :59982 | grep \\*:59982")
	suite.Empty(out, "Listening on 0.0.0.0 (not allowed)")

	out = utils.Command("sudo lsof -nPi :59982 | grep 127.0.0.1:59982 ")
	suite.NotEmpty(out, "Service is not listening on port 59982")
}
