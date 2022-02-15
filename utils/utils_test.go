package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommand(t *testing.T) {
	t.Run("one command", func(t *testing.T) {
		stdout, stderr := RunCommand(t, `echo "hi"`)
		assert.Empty(t, stderr)
		require.Equal(t, "hi\n", stdout)
	})

	t.Run("exit after slow command", func(t *testing.T) {
		start := time.Now()
		stdout, _ := RunCommand(t, `echo "hi" && sleep 1 && echo "hi2"`)
		// must return after 1s±200ms
		require.WithinDuration(t,
			start.Add(1*time.Second),
			time.Now(),
			200*time.Millisecond)
		require.Equal(t, "hi\nhi2\n", stdout)
	})

	t.Run("multiple commands", func(t *testing.T) {
		stdout, _ := RunCommand(t,
			`echo "hi"`,
			`echo "hi2"`,
		)
		require.Equal(t, "hi\nhi2\n", stdout)
	})

	t.Run("bad command", func(t *testing.T) {
		_, stderr := RunCommand(t, `bad_command`)
		require.NotEmpty(t, stderr)
	})

	t.Run("redirect stdout to stderr", func(t *testing.T) {
		stdout, stderr := RunCommand(t, `echo "hello" >&2`)
		assert.Empty(t, stdout)
		require.Contains(t, stderr, "hello")
	})

	t.Run("bad command, redirects stderr to stdout", func(t *testing.T) {
		// Do not pass t which raises the error because we want to
		// validate the error handling
		stdout, stderr := RunCommand(nil, `bad_command 2>&1`)
		assert.Empty(t, stderr)
		require.Contains(t, stdout, "not found")
	})

	t.Run("bad and good commands", func(t *testing.T) {
		// Do not pass t which raises the error because we want to
		// validate the error handling
		t.Run("bad+good", func(t *testing.T) {
			stdout, stderr := RunCommand(nil,
				`bad_command`,
				`echo 'hi'`,
			)
			require.Contains(t, stderr, "not found")
			assert.Contains(t, stdout, "hi")
		})
		t.Run("good+bad", func(t *testing.T) {
			stdout, stderr := RunCommand(nil,
				`echo 'hi'`,
				`bad_command`,
			)
			require.Contains(t, stderr, "not found")
			assert.Contains(t, stdout, "hi")
		})
	})
}
