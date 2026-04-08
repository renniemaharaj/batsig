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

func chargingAlertsDir() string {
	return filepath.Join(alertsDir(), "charging")
}

func stateAlertsDir() string {
	return filepath.Join(alertsDir(), "state")
}

func stateAlertPath(charging bool) string {
	if charging {
		return filepath.Join(stateAlertsDir(), "charging")
	}
	return filepath.Join(stateAlertsDir(), "discharging")
}

func alertDir(charging bool) string {
	if charging {
		return chargingAlertsDir()
	}
	return alertsDir()
}

func ensureAlertsDir() error {
	if err := os.MkdirAll(alertsDir(), 0o755); err != nil {
		return err
	}
	if err := os.MkdirAll(chargingAlertsDir(), 0o755); err != nil {
		return err
	}
	return os.MkdirAll(stateAlertsDir(), 0o755)
}

func parseAlertFlags(args []string) (bool, []string, bool, error) {
	if len(args) == 0 {
		return false, args, false, nil
	}

	switch args[0] {
	case "-c", "--charging":
		return true, args[1:], true, nil
	case "-d", "--discharging":
		return false, args[1:], true, nil
	}

	return false, args, false, nil
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

func createAlert(pct int, message string, charging bool) error {
	alertPath := filepath.Join(alertDir(charging), strconv.Itoa(pct))
	return os.WriteFile(alertPath, []byte(message), 0o644)
}

func removeAlert(pct int, charging bool) error {
	alertPath := filepath.Join(alertDir(charging), strconv.Itoa(pct))
	return os.Remove(alertPath)
}

func readAlertMessage(pct int, charging bool) (string, error) {
	alertPath := filepath.Join(alertDir(charging), strconv.Itoa(pct))
	data, err := os.ReadFile(alertPath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func createStateAlert(charging bool, message string) error {
	if err := os.MkdirAll(stateAlertsDir(), 0o755); err != nil {
		return err
	}
	return os.WriteFile(stateAlertPath(charging), []byte(message), 0o644)
}

func readStateAlert(charging bool) (string, error) {
	data, err := os.ReadFile(stateAlertPath(charging))
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func removeStateAlert(charging bool) error {
	return os.Remove(stateAlertPath(charging))
}

func moveAlert(pct int, charging bool) error {
	dest := filepath.Join(alertDir(charging), strconv.Itoa(pct))
	if _, err := os.Stat(dest); err == nil {
		return nil
	}

	src := filepath.Join(alertDir(!charging), strconv.Itoa(pct))
	if _, err := os.Stat(src); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("alert %d not found", pct)
		}
		return err
	}

	return os.Rename(src, dest)
}
