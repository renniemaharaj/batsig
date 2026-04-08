package cli

import (
	"errors"
	"fmt"
	"io"
	"os"
)

var ErrUsage = errors.New("usage")

func Execute() error {
	if err := ensureAlertsDir(); err != nil {
		return fmt.Errorf("failed to create alerts directory: %w", err)
	}

	args := os.Args[1:]
	if len(args) == 0 {
		return runMonitor()
	}

	switch args[0] {
	case "daemon":
		return runDaemon(args[1:])
	case "install":
		return runInstall()
	case "alert":
		if len(args) < 2 {
			return ErrUsage
		}
		switch args[1] {
		case "-c", "--charging", "-d", "--discharging":
			return runStateAlert(args[1:])
		case "new":
			return runNew(args[2:])
		case "set":
			return runSet(args[2:])
		case "clear":
			return runClear(args[2:])
		case "test":
			return runTest(args[2:])
		case "mv":
			return runMv(args[2:])
		case "list":
			return runList(args[2:])
		default:
			return ErrUsage
		}
	default:
		return ErrUsage
	}
}

func PrintUsage(w io.Writer) {
	fmt.Fprintln(w, "Usage:")
	fmt.Fprintln(w, "  batsig                # monitor battery and send alerts")
	fmt.Fprintln(w, "  batsig daemon         # detach and run the monitor in the background")
	fmt.Fprintln(w, "  batsig install        # install batsig binary, notification sound, and alerts to your local config")
	fmt.Fprintln(w, "  batsig alert new [-c] <percentage> <message>    # create an alert for discharging or charging with -c")
	fmt.Fprintln(w, "  batsig alert set [-c] <percentage> <message>    # create or update an alert")
	fmt.Fprintln(w, "  batsig alert clear [-c] <percentage>            # remove an alert from the selected folder")
	fmt.Fprintln(w, "  batsig alert test [-c] <percentage>             # trigger a specific alert for testing")
	fmt.Fprintln(w, "  batsig alert mv <percentage> <charging|discharging>  # move an alert between folders")
	fmt.Fprintln(w, "  batsig alert list                               # list configured alerts and their messages")
}
