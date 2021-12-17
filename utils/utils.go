package utils

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"os/exec"
	"testing"
)

const ShellToUse = "bash"

func Exec(command string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(ShellToUse, "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func Command(t *testing.T, command string) string {
	stdout, stderr, err := Exec(command)

	assert.Empty(t, err, "Execution error: %s; Stdout: %s;", err, stderr)
	return stdout
}

	fmt.Println(errout)

	return out
