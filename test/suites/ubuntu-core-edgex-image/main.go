package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	scripts := []string{"setup.sh", "build.sh", "test.sh"}

	for _, s := range scripts {
		cmd := exec.Command("bash", s)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to execute %s: %s\n", s, err)
			os.Exit(1)
		}
	}
}
