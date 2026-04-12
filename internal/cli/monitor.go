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

const (
	pollInterval       = 30 * time.Second
	firedResetInterval = 5 * time.Minute
)

func runMonitor() error {
	ctx := context.Background()
	fired := map[string]bool{}
	var prevChargingState *bool

	pollTicker := time.NewTicker(pollInterval)
	defer pollTicker.Stop()
	resetTicker := time.NewTicker(firedResetInterval)
	defer resetTicker.Stop()

	status, err := upower.Status()
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: unable to read battery status: %v\n", err)
	} else {
		checkAlerts(ctx, status, fired, &prevChargingState)
	}

	for {
		select {
		case <-pollTicker.C:
			status, err := upower.Status()
			if err != nil {
				fmt.Fprintf(os.Stderr, "warning: unable to read battery status: %v\n", err)
			} else {
				checkAlerts(ctx, status, fired, &prevChargingState)
			}
		case <-resetTicker.C:
			fired = map[string]bool{}
		}
	}
}

func checkAlerts(ctx context.Context, status upower.BatteryStatus, fired map[string]bool, prevChargingState **bool) {
	charging := status.IsCharging()
	if prevChargingState == nil || *prevChargingState == nil || **prevChargingState != charging {
		if err := sendStateAlert(ctx, charging); err != nil {
			fmt.Fprintf(os.Stderr, "warning: state alert failed: %v\n", err)
		}
		if *prevChargingState == nil {
			b := charging
			*prevChargingState = &b
		} else {
			**prevChargingState = charging
		}
	}

	if !charging && !status.IsDischarging() {
		return
	}

	entries, err := os.ReadDir(alertDir(charging))
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

	candidate := selectExactThreshold(activeThresholds, status.Percentage)
	if candidate >= 0 {
		key := fmt.Sprintf("%t:%d", charging, candidate)
		if !fired[key] {
			message, err := os.ReadFile(filepath.Join(alertDir(charging), strconv.Itoa(candidate)))
			if err != nil {
				fmt.Fprintf(os.Stderr, "warning: unable to read alert %d: %v\n", candidate, err)
				return
			}

			err = notify.Send(ctx, fmt.Sprintf("Battery %d%%", candidate), strings.TrimSpace(string(message)))
			if err != nil {
				fmt.Fprintf(os.Stderr, "warning: notify-send failed: %v\n", err)
			} else {
				fired[key] = true
			}
		}
	}
}

func sendStateAlert(ctx context.Context, charging bool) error {
	message, err := readStateAlert(charging)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	var title string
	if charging {
		title = "Charging"
	} else {
		title = "Discharging"
	}

	return notify.Send(ctx, title, strings.TrimSpace(message))
}

func selectExactThreshold(thresholds []int, current int) int {
	sort.Ints(thresholds)
	for _, pct := range thresholds {
		if pct == current {
			return pct
		}
	}
	return -1
}
