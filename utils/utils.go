package utils

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
)

const ShellToUse = "bash"

func Exit(command string) (error, string, string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(ShellToUse, "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return err, stdout.String(), stderr.String()
}

func Command(command string) string {
	err, out, errout := Exit(command)
	if err != nil {
		log.Printf("[LOG] `%s`, error: %v\n", command, err)
	}
	fmt.Println(out)
	fmt.Println(errout)

	return out
}
