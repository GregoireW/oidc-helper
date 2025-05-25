package daemon

import (
	"os"
	"os/exec"
)

// LaunchDaemon launches the daemon process with --daemon flag
func LaunchDaemon(pipeName string) error {
	args := []string{"--daemon", "--pipe", pipeName}
	cmd := exec.Command(os.Args[0], args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Start()
}
