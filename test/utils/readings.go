package utils

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
	"time"
)

// TestDeviceVirtualReading waits for device-virtual to produce readings by querying core-data
// up to a maximun number
func TestDeviceVirtualReading(t *testing.T) {
	t.Run("query readings", func(t *testing.T) {
		var count int

		// wait device-virtual producing readings with maximum 60 seconds
		for i := 1; ; i++ {
			time.Sleep(1 * time.Second)
			req, err := http.NewRequest("GET", "http://localhost:59880/api/v2/event/count", nil)
			if err != nil {
				fmt.Print(err)
				return
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				fmt.Print(err)
				return
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Print(err)
				return
			}

			mapContainer := make(map[string]json.RawMessage)
			err = json.Unmarshal(body, &mapContainer)
			if err != nil {
				fmt.Print(err)
				return
			}

			c := mapContainer["Count"]
			count, _ = strconv.Atoi(string(c))

			t.Logf("waiting for device-virtual produce readings, current retry count: %d/60\n", i)

			if i <= 60 && count > 0 {
				t.Logf("device-virtual is producing readings now, readings queried from core-data")
				break
			}

			if i > 60 && count <= 0 {
				t.Logf("waiting for device-virtual produce readings, reached maximum retry count of 60")
				break
			}
		}
		require.Greaterf(t, count, 0, "No device-virtual reading in core-data")
	})
}
