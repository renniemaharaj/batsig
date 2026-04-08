package cli

import (
	"fmt"
	"strings"
)

func runNew(args []string) error {
	charging, args, _, err := parseAlertFlags(args)
	if err != nil {
		return err
	}
	if len(args) < 2 {
		return ErrUsage
	}

	pct, err := parsePercentage(args[0])
	if err != nil {
		return fmt.Errorf("invalid percentage: %w", err)
	}

	return createAlert(pct, strings.Join(args[1:], " "), charging)
}
