package notify

import (
	"context"
	"fmt"
	"os/exec"
)

func Send(ctx context.Context, title, body string) error {
	if body == "" {
		body = title
	}

	cmd := exec.CommandContext(ctx, "notify-send", title, body)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("notify-send failed: %w", err)
	}

	return nil
}
