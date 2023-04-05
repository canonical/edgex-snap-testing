package utils

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func LoginTestUser(t *testing.T) (idToken string) {
	// The script path relative to the testing suites
	const loginScriptPath = "../../utils/login-test-user.sh"

	idToken, _, _ = Exec(t, loginScriptPath)
	t.Log("ID Token for 'example' user:", idToken)
	return strings.TrimSpace(idToken)
}

// WaitForReadings waits for readings to appear in core-data
// The readings are produced by device-virtual or another service
func WaitForReadings(t *testing.T, secured bool) {
	const coreDataReadingCountEndpoint = "http://localhost:59880/api/v2/reading/count"

	t.Run("query readings count", func(t *testing.T) {
		var eventCount struct {
			Count int
		}

		var idToken string
		if secured {
			idToken = LoginTestUser(t)
		}

		for i := 1; ; i++ {
			time.Sleep(1 * time.Second)
			req, err := http.NewRequest(http.MethodGet, coreDataReadingCountEndpoint, nil)
			require.NoError(t, err)

			if secured {
				req.Header.Set("Authorization", "Bearer "+idToken)
			}

			client := &http.Client{}
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()
			require.Equal(t, 200, resp.StatusCode, "Unexpected HTTP response")

			require.NoError(t, json.NewDecoder(resp.Body).Decode(&eventCount))

			t.Logf("Waiting for readings in Core Data, current retry: %d/60", i)

			if i <= 60 && eventCount.Count > 0 {
				t.Logf("Found readings in Core Data")
				break
			}

			if i > 60 && eventCount.Count <= 0 {
				t.Logf("Waiting for readings in Core Data: reached maximum 60 retries")
				break
			}
		}
		require.Greaterf(t, eventCount.Count, 0, "No readings in Core Data")
	})
}
