package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExec(t *testing.T) {

	t.Run("one command", func(t *testing.T) {
		stdout, stderr := Exec(t, `echo "hi"`)
		assert.Empty(t, stderr)
		assert.Equal(t, "hi\n", stdout)
	})

	t.Run("exit after slow command", func(t *testing.T) {
		start := time.Now()
		stdout, _ := Exec(t, `echo "hi" && sleep 0.1 && echo "hi2"`)
		// must return after 100ms±50ms
		require.WithinDuration(t,
			start.Add(100*time.Millisecond),
			time.Now(),
			50*time.Millisecond)
		require.Equal(t, "hi\nhi2\n", stdout)
	})

	t.Run("bad command", func(t *testing.T) {
		testingFatal = true
		t.Cleanup(func() {
			testingFatal = false
		})

		stdout, stderr := Exec(t, `bad_command`)
		assert.Empty(t, stdout)
		assert.Contains(t, stderr, "not found")
	})

	t.Run("print to stderr", func(t *testing.T) {
		stdout, stderr := Exec(t, `echo "failing" >&2`)
		assert.Empty(t, stdout)
		assert.Equal(t, "failing\n", stderr)
	})

	t.Run("stderr then stdout", func(t *testing.T) {
		stdout, stderr := Exec(t, `echo "failing" >&2; echo "succeeding"`)
		assert.Equal(t, "failing\n", stderr)
		assert.Equal(t, "succeeding\n", stdout)
	})

}
