package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"batsig/internal/notify"
	"batsig/internal/upower"
)

const pollInterval = 30 * time.Second

func runMonitor() error {
	ctx := context.Background()
	fired := map[int]bool{}

	for {
		status, err := upower.Status()
		if err != nil {
			fmt.Fprintf(os.Stderr, "warning: unable to read battery status: %v\n", err)
		} else {
			checkAlerts(ctx, status, fired)
		}
		time.Sleep(pollInterval)
	}
}

func checkAlerts(ctx context.Context, status upower.BatteryStatus, fired map[int]bool) {
	entries, err := os.ReadDir(alertsDir())
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: unable to read alerts directory: %v\n", err)
		return
	}

	activeThresholds := []int{}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		pct, err := parsePercentage(entry.Name())
		if err != nil {
			continue
		}

		activeThresholds = append(activeThresholds, pct)
	}

	candidate := selectLowestAboveOrEqualThreshold(activeThresholds, status.Percentage)
	if candidate >= 0 && status.IsDischarging() {
		if !fired[candidate] {
			message, err := os.ReadFile(filepath.Join(alertsDir(), strconv.Itoa(candidate)))
			if err != nil {
				fmt.Fprintf(os.Stderr, "warning: unable to read alert %d: %v\n", candidate, err)
			} else {
				err = notify.Send(ctx, fmt.Sprintf("Battery %d%%", candidate), strings.TrimSpace(string(message)))
				if err != nil {
					fmt.Fprintf(os.Stderr, "warning: notify-send failed: %v\n", err)
				} else {
					fired[candidate] = true
				}
			}
		}
	}
}

func selectLowestAboveOrEqualThreshold(thresholds []int, current int) int {
	sort.Ints(thresholds)
	for _, pct := range thresholds {
		if pct >= current {
			return pct
		}
	}
	return -1
}
