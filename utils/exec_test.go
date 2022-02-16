package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunCommand(t *testing.T) {

	t.Run("one command", func(t *testing.T) {
		stdout, stderr := RunCommand(t, `echo "hi"`)
		assert.Empty(t, stderr)
		assert.Equal(t, "hi\n", stdout)
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
		assert.Equal(t, "hi\nhi2\n", stdout)
	})

	t.Run("bad command", func(t *testing.T) {
		testingFatal = true
		t.Cleanup(func() {
			testingFatal = false
		})

		stdout, stderr := RunCommand(t, `bad_command`)
		assert.Empty(t, stdout)
		assert.Contains(t, stderr, "not found")
	})

	t.Run("print to stderr", func(t *testing.T) {
		stdout, stderr := RunCommand(t, `echo "failing" >&2`)
		assert.Empty(t, stdout)
		assert.Equal(t, "failing\n", stderr)
	})

	t.Run("stderr then stdout", func(t *testing.T) {
		stdout, stderr := RunCommand(t, `echo "failing" >&2; echo "succeeding"`)
		assert.Equal(t, "failing\n", stderr)
		assert.Equal(t, "succeeding\n", stdout)
	})

	t.Run("bad then good", func(t *testing.T) {
		testingFatal = true
		t.Cleanup(func() {
			testingFatal = false
		})

		stdout, stderr := RunCommand(t,
			`bad_command`, // it must stop after this
			`echo 'hi'`,
		)
		assert.Contains(t, stderr, "not found")
		assert.NotContains(t, stdout, "hi")
	})

	t.Run("good then bad", func(t *testing.T) {
		testingFatal = true
		t.Cleanup(func() {
			testingFatal = false
		})

		stdout, stderr := RunCommand(t,
			`echo 'hi'`,
			`bad_command`,
		)
		assert.Contains(t, stdout, "hi")
		assert.Contains(t, stderr, "not found")
	})
}
