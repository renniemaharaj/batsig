package cli

import "fmt"

func runMv(args []string) error {
	if len(args) != 2 {
		return ErrUsage
	}

	pct, err := parsePercentage(args[0])
	if err != nil {
		return fmt.Errorf("invalid percentage: %w", err)
	}

	target := args[1]
	switch target {
	case "charging":
		return moveAlert(pct, true)
	case "discharging":
		return moveAlert(pct, false)
	default:
		return ErrUsage
	}
}
