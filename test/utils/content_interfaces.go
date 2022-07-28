package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type ContentInterfaces struct {
	TestSecretstoreToken bool
	TestConfigProvider   bool
	Snap                 string
	App                  string
}

func TestContentInterfaces(t *testing.T, conf ContentInterfaces) {
	t.Run("content interfaces", func(t *testing.T) {
		if conf.TestSecretstoreToken {
			TestSeedSecretstoreToken(t, conf.Snap, conf.App)
		}
		if conf.TestConfigProvider {
			t.Skip("TODO")
		}
	})
}

func TestSeedSecretstoreToken(t *testing.T, consumerName, consumerAppName string) {
	t.Run("seed secretstore token", func(t *testing.T) {
		providerSecret, _, _ := Exec(t, "sudo cat /var/snap/edgexfoundry/current/secrets/"+consumerAppName+"/secrets-token.json")
		consumerSecret, _, _ := Exec(t, "sudo cat /var/snap/"+consumerName+"/current/"+consumerAppName+"/secrets-token.json")
		require.NotEmpty(t, providerSecret)

		require.Equal(t, providerSecret, consumerSecret)
	})
}
