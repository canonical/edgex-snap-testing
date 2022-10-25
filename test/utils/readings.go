package utils

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
	"time"
)

type EventCount struct {
	Count int `json:"Count"`
}

// TestDeviceVirtualReading waits for device-virtual to produce readings by querying core-data
// up to a maximun number
func TestDeviceVirtualReading(t *testing.T) {
	t.Run("query readings", func(t *testing.T) {
		var eventCount EventCount

		// wait device-virtual to produce readings with maximum 60 seconds
		for i := 1; ; i++ {
			time.Sleep(1 * time.Second)
			resp, err := http.Get("http://localhost:59880/api/v2/event/count")
			if err != nil {
				fmt.Print(err)
				return
			}
			defer resp.Body.Close()

			if err = json.NewDecoder(resp.Body).Decode(&eventCount); err != nil {
				t.Fatal(err)
			}

			t.Logf("waiting for device-virtual to produce readings, current retry count: %d/60\n", i)

			if i <= 60 && eventCount.Count > 0 {
				t.Logf("device-virtual is producing readings now, readings queried from core-data")
				break
			}

			if i > 60 && eventCount.Count <= 0 {
				t.Logf("waiting for device-virtual to produce readings, reached maximum retry count of 60")
				break
			}
		}
		require.Greaterf(t, eventCount.Count, 0, "No device-virtual reading in core-data")
	})
}
