package cli

import "fmt"

func runClear(args []string) error {
	charging, args, flagPresent, err := parseAlertFlags(args)
	if err != nil {
		return err
	}
	if len(args) == 0 {
		if !flagPresent {
			return ErrUsage
		}
		return removeStateAlert(charging)
	}
	if len(args) != 1 {
		return ErrUsage
	}

	pct, err := parsePercentage(args[0])
	if err != nil {
		return fmt.Errorf("invalid percentage: %w", err)
	}

	return removeAlert(pct, charging)
}
