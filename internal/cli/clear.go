package cli

import "fmt"

func runClear(args []string) error {
	if len(args) != 1 {
		return ErrUsage
	}

	pct, err := parsePercentage(args[0])
	if err != nil {
		return fmt.Errorf("invalid percentage: %w", err)
	}

	return removeAlert(pct)
}
