package upower

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type BatteryStatus struct {
	Percentage int
	State      string
}

func (s BatteryStatus) IsDischarging() bool {
	return strings.EqualFold(s.State, "discharging")
}

func Status() (BatteryStatus, error) {
	path, err := BatteryDevice()
	if err != nil {
		return BatteryStatus{}, err
	}

	output, err := exec.Command("upower", "-i", path).CombinedOutput()
	if err != nil {
		return BatteryStatus{}, fmt.Errorf("upower -i failed: %w: %s", err, string(output))
	}

	var status BatteryStatus
	foundPercentage := false
	for _, line := range strings.Split(string(output), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "percentage:") {
			parts := strings.Fields(line)
			if len(parts) < 2 {
				continue
			}
			pctString := strings.TrimSuffix(parts[1], "%")
			status.Percentage, err = parseInt(pctString)
			if err != nil {
				return BatteryStatus{}, err
			}
			foundPercentage = true
			continue
		}
		if strings.HasPrefix(line, "state:") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				status.State = strings.ToLower(parts[1])
			}
		}
	}

	if !foundPercentage {
		return BatteryStatus{}, fmt.Errorf("percentage not found in upower output")
	}

	return status, nil
}

func Percentage() (int, error) {
	status, err := Status()
	if err != nil {
		return 0, err
	}
	return status.Percentage, nil
}

func BatteryDevice() (string, error) {
	output, err := exec.Command("upower", "-e").CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("upower -e failed: %w: %s", err, string(output))
	}

	for _, line := range strings.Split(string(output), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.Contains(strings.ToLower(line), "battery") {
			return line, nil
		}
	}

	return "", fmt.Errorf("no battery device found")
}

func parseInt(value string) (int, error) {
	return strconv.Atoi(value)
}
