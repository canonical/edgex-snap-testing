package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	scripts := []string{"setup.sh", "build.sh", "run.sh"}

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

	// todo: find a way to execute test.sh after compeleting Ubuntu Core Setup in run.sh script
	// cmd := exec.Command("bash", "test.sh")
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	// if err := cmd.Run(); err != nil {
    // 	fmt.Fprintf(os.Stderr, "Failed to execute test.sh: %s\n", err)
   	// 	os.Exit(1)
	// }
}
