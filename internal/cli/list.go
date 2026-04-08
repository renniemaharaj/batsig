package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func runList(args []string) error {
	if len(args) != 0 {
		return ErrUsage
	}

	if err := printStateAlert("charging", true); err != nil {
		return err
	}
	if err := printStateAlert("discharging", false); err != nil {
		return err
	}
	if err := printAlertSection("discharging", alertsDir()); err != nil {
		return err
	}
	return printAlertSection("charging", chargingAlertsDir())
}

func printStateAlert(name string, charging bool) error {
	message, err := readStateAlert(charging)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("unable to read %s state alert: %w", name, err)
	}

	fmt.Printf("%s state: %s\n", strings.Title(name), strings.TrimSpace(message))
	return nil
}

func printAlertSection(name, dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("unable to read %s alerts directory: %w", name, err)
	}

	thresholds := []int{}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		pct, err := parsePercentage(entry.Name())
		if err != nil {
			continue
		}

		thresholds = append(thresholds, pct)
	}

	sort.Ints(thresholds)
	for _, pct := range thresholds {
		message, err := os.ReadFile(filepath.Join(dir, strconv.Itoa(pct)))
		if err != nil {
			return fmt.Errorf("unable to read alert %d: %w", pct, err)
		}

		fmt.Printf("%s %d%%: %s\n", strings.Title(name), pct, strings.TrimSpace(string(message)))
	}

	return nil
}
