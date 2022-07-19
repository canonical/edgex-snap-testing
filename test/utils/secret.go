package utils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSecret(t *testing.T, consumerName string) {
	t.Run("secrets interface", func(t *testing.T) {
		TestSeededSecretstoreToken(t, consumerName)
	})
}

func TestSeededSecretstoreToken(t *testing.T, consumerName string) {
	t.Run("same content on both sides", func(t *testing.T) {
		providerSecret, _ := Exec(t, "sudo cat /var/snap/edgexfoundry/current/secrets/"+consumerName+"/secrets-token.json")
		consumerSecret, _ := Exec(t, "sudo cat /var/snap/edgex-"+consumerName+"/current/"+consumerName+"/secrets-token.json")
		require.NotEmpty(t, providerSecret)

		require.Equal(t, providerSecret, consumerSecret)
	})
}
