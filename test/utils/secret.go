package utils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

type Secret struct {
	TestSecretsInterface bool
}

// TODO: rename to a more generic name such as TestContentInterfaces
// to allow having both secret seeding and config interfaces inside
func TestSecret(t *testing.T, providerPath, consumerName, consumerSecretPath string, conf Secret) {
	t.Run("secrets interface", func(t *testing.T) {
		if conf.TestSecretsInterface {
			TestSeededSecretstoreToken(t, providerPath, consumerName, consumerSecretPath)
		} else {
			t.Skip("Security off mode, there is no secret needed.")
		}
	})
}

func TestSeededSecretstoreToken(t *testing.T, providerPath, consumerName, consumerPath string) {
	t.Run("same content on both sides", func(t *testing.T) {
		providerSecret, _ := Exec(t, "sudo cat /var/snap/edgexfoundry/current/secrets/"+providerPath+"/secrets-token.json")
		consumerSecret, _ := Exec(t, "sudo cat /var/snap/"+consumerName+"/current/"+consumerPath+"/secrets-token.json")
		require.NotEmpty(t, providerSecret)

		require.Equal(t, providerSecret, consumerSecret)
	})
}
