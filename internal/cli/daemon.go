package cli

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

const daemonEnv = "BATSIG_DAEMON"

func runDaemon(args []string) error {
	if os.Getenv(daemonEnv) == "1" {
		return runMonitor()
	}

	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("unable to locate executable: %w", err)
	}

	devNull, err := os.OpenFile("/dev/null", os.O_RDWR, 0)
	if err != nil {
		return fmt.Errorf("unable to open /dev/null: %w", err)
	}
	defer devNull.Close()

	cmd := exec.Command(exe, "daemon")
	cmd.Env = append(os.Environ(), daemonEnv+"=1")
	cmd.Stdin = devNull
	cmd.Stdout = devNull
	cmd.Stderr = devNull
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start daemon: %w", err)
	}

	return nil
}
