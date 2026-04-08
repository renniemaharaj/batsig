package cli

import "strings"

func runStateAlert(args []string) error {
	if len(args) < 2 {
		return ErrUsage
	}

	charging := false
	switch args[0] {
	case "-c", "--charging":
		charging = true
	case "-d", "--discharging":
		charging = false
	default:
		return ErrUsage
	}

	message := strings.TrimSpace(strings.Join(args[1:], " "))
	if message == "" {
		return ErrUsage
	}

	return createStateAlert(charging, message)
}
