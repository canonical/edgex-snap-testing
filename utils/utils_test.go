package utils

import (
	"testing"
)

func TestUtils(t *testing.T) {
	t.Run("one command", func(t *testing.T) {
		Command(`echo "hi" && sleep 1 && echo "hi2"`)
	})
	t.Run("multiple commands", func(t *testing.T) {
		Command(
			`echo "hi" && sleep 1 && echo "hi2"`,
			`echo "hi3"`,
		)
	})

	t.Run("bad command", func(t *testing.T) {
		stdout, stderr, err := Command(`bad_command`)
		CommandLog(t, stdout, stderr, err)
	})

	t.Run("bad command, redirects stdout in stderr", func(t *testing.T) {
		stdout, stderr, err := Command(`bad_command >&2`)
		CommandLog(t, stdout, stderr, err)
	})

	t.Run("bad command, redirects stderr to stdout", func(t *testing.T) {
		stdout, stderr, err := Command(`bad_command 2>&1`)
		CommandLog(t, stdout, stderr, err)
	})

	// expect exit code 1
	// got exit code 0
	t.Run("single bad command", func(t *testing.T) {
		Command(`bad_command && echo $?`)
	})

	// expect exit code 1
	// got exit code 0
	t.Run("multiple commands include bad command", func(t *testing.T) {
		stdout, stderr, err := Command("bad_command",
			"echo $?")
		CommandLog(t, stdout, stderr, err)
	})
}

func TestLogger(t *testing.T) {
	t.Run("bad command", func(t *testing.T) {
		stdout, stderr, err := Command(`bad_command`)
		CommandLog(t, stdout, stderr, err)
	})
	t.Run("good command", func(t *testing.T) {
		stdout, stderr, err := Command(`ls`)
		CommandLog(t, stdout, stderr, err)
	})
}
