package utils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

type Secret struct {
	TestManualConnection bool
}

func TestSecret(t *testing.T, consumerName string, conf Secret) {
	t.Run("secrets interface", func(t *testing.T) {
		if conf.TestManualConnection {
			SnapConnect(nil,
				"edgexfoundry:edgex-secretstore-token",
				"edgex-"+consumerName+":edgex-secretstore-token",
			)
		}
		requireCorrectSecret(t, consumerName)
	})
}

func requireCorrectSecret(t *testing.T, consumerName string) {
	providerSecret, stderr := Exec(t, "sudo cat /var/snap/edgexfoundry/current/secrets/"+consumerName+"/secrets-token.json")
	require.Empty(t, stderr)

	consumerSecret, stderr := Exec(t, "sudo cat /var/snap/edgex-"+consumerName+"/current/"+consumerName+"/secrets-token.json")
	require.Empty(t, stderr)

	t.Run("same content on both sides", func(t *testing.T) {
		require.Equal(t, providerSecret, consumerSecret)
	})
}
