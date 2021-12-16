package test

import (
	"MQTT-test-project/utils"
	"github.com/stretchr/testify/suite"
	"testing"
)

type MqttTestSuite struct {
	suite.Suite
}

func (suite *MqttTestSuite) SetupSuite() {
	utils.Command("sudo snap install edgexfoundry --channel=2.1/stable")
	utils.Command("sudo snap install edgex-device-mqtt --channel=2.1/stable")
}

func (suite *MqttTestSuite) BeforeTest(_, _ string) {
	utils.Command("sudo snap start --enable edgex-device-mqtt.device-mqtt")
}

func (suite *MqttTestSuite) AfterTest(_, _ string) {
	utils.Command("sudo snap stop --disable edgex-device-mqtt.device-mqtt")
}

func (suite *MqttTestSuite) TearDownSuite() {
	utils.Command("sudo snap remove --purge edgex-device-mqtt")
	utils.Command("sudo snap remove --purge edgexfoundry")
}

func TestMqttTestSuite(t *testing.T) {
	suite.Run(t, new(MqttTestSuite))
}
