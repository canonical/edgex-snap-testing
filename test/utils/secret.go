package utils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

type Secret struct {
	TestSecretToken bool
	Snap            string
	App             string
}

// TODO: rename to a more generic name such as TestContentInterfaces
// to allow having both secret seeding and config interfaces inside
func TestSecret(t *testing.T, conf Secret) {
	t.Run("secrets interface", func(t *testing.T) {
		if conf.TestSecretToken {
			TestSeededSecretstoreToken(t, conf.Snap, conf.App)
		} else {
			t.Skip("Security off mode, there is no secret needed.")
		}
	})
}

func TestSeededSecretstoreToken(t *testing.T, consumerName, consumerAppName string) {
	t.Run("same content on both sides", func(t *testing.T) {
		providerSecret, _ := Exec(t, "sudo cat /var/snap/edgexfoundry/current/secrets/"+consumerAppName+"/secrets-token.json")
		consumerSecret, _ := Exec(t, "sudo cat /var/snap/"+consumerName+"/current/"+consumerAppName+"/secrets-token.json")
		require.NotEmpty(t, providerSecret)

		require.Equal(t, providerSecret, consumerSecret)
	})
}
