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

	entries, err := os.ReadDir(alertsDir())
	if err != nil {
		return fmt.Errorf("unable to read alerts directory: %w", err)
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
		message, err := os.ReadFile(filepath.Join(alertsDir(), strconv.Itoa(pct)))
		if err != nil {
			return fmt.Errorf("unable to read alert %d: %w", pct, err)
		}

		fmt.Printf("%d%%: %s\n", pct, strings.TrimSpace(string(message)))
	}

	return nil
}
