package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

const alertsDirName = "batsig/alerts"

func configDir() string {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return xdg
	}
	if home := os.Getenv("HOME"); home != "" {
		return filepath.Join(home, ".config")
	}
	return ".config"
}

func alertsDir() string {
	return filepath.Join(configDir(), alertsDirName)
}

func ensureAlertsDir() error {
	return os.MkdirAll(alertsDir(), 0o755)
}

func parsePercentage(value string) (int, error) {
	pct, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}
	if pct < 0 || pct > 100 {
		return 0, fmt.Errorf("percentage must be between 0 and 100")
	}
	return pct, nil
}

func createAlert(pct int, message string) error {
	alertPath := filepath.Join(alertsDir(), strconv.Itoa(pct))
	return os.WriteFile(alertPath, []byte(message), 0o644)
}

func removeAlert(pct int) error {
	alertPath := filepath.Join(alertsDir(), strconv.Itoa(pct))
	return os.Remove(alertPath)
}

func readAlertMessage(pct int) (string, error) {
	alertPath := filepath.Join(alertsDir(), strconv.Itoa(pct))
	data, err := os.ReadFile(alertPath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
