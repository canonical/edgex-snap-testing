package latest

import (
	"MQTT-test-project/utils"
	"fmt"
)

func NetworkInterface() {
	out := utils.Command("sudo lsof -nPi :59982 | grep \\*:59982")
	if out != "" {
		fmt.Printf("listening on 0.0.0.0. It should not allow remote listening\n")
	}

	out = utils.Command("sudo lsof -nPi :59982 | grep 127.0.0.1:59982 ")
	if out != "" {
		fmt.Printf("listening on 127.0.0.1\n")
	}
}
