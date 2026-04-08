package cli

import (
	"context"
	"fmt"
	"strings"

	"batsig/internal/notify"
)

func runTest(args []string) error {
	charging, args, flagPresent, err := parseAlertFlags(args)
	if err != nil {
		return err
	}

	if len(args) == 0 {
		if !flagPresent {
			return ErrUsage
		}
		message, err := readStateAlert(charging)
		if err != nil {
			return err
		}
		return notify.Send(context.Background(), "Battery state test", strings.TrimSpace(message))
	}

	if len(args) != 1 {
		return ErrUsage
	}

	pct, err := parsePercentage(args[0])
	if err != nil {
		return fmt.Errorf("invalid percentage: %w", err)
	}

	message, err := readAlertMessage(pct, charging)
	if err != nil {
		return err
	}

	return notify.Send(context.Background(), fmt.Sprintf("Battery %d%% test", pct), strings.TrimSpace(message))
}
