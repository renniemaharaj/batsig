package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParsePercentage(t *testing.T) {
	cases := []struct {
		input string
		want  int
		ok    bool
	}{
		{"0", 0, true},
		{"50", 50, true},
		{"100", 100, true},
		{"-1", 0, false},
		{"101", 0, false},
		{"foo", 0, false},
	}

	for _, tc := range cases {
		got, err := parsePercentage(tc.input)
		if (err == nil) != tc.ok {
			t.Errorf("parsePercentage(%q) ok=%v, err=%v", tc.input, tc.ok, err)
			continue
		}
		if tc.ok && got != tc.want {
			t.Errorf("parsePercentage(%q) = %d, want %d", tc.input, got, tc.want)
		}
	}
}

func TestAlertFileLifecycle(t *testing.T) {
	dir := t.TempDir()
	origWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.Chdir(origWd)
	}()

	if err := os.Setenv("XDG_CONFIG_HOME", dir); err != nil {
		t.Fatal(err)
	}

	if err := ensureAlertsDir(); err != nil {
		t.Fatal(err)
	}

	if err := createAlert(80, "charge test"); err != nil {
		t.Fatal(err)
	}

	message, err := readAlertMessage(80)
	if err != nil {
		t.Fatal(err)
	}
	if message != "charge test" {
		t.Fatalf("unexpected alert message: %q", message)
	}

	if err := removeAlert(80); err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat(filepath.Join(dir, "batsig", "alerts", "80"))
	if err == nil || !os.IsNotExist(err) {
		t.Fatalf("expected removed alert file, got err=%v", err)
	}
}

func TestExecuteUnknownCommand(t *testing.T) {
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	os.Args = []string{"batsig", "unknown"}
	if err := Execute(); err != ErrUsage {
		t.Fatalf("Execute() = %v, want ErrUsage", err)
	}
}

func TestExecuteSetCommand(t *testing.T) {
	dir := t.TempDir()
	origWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.Chdir(origWd)
	}()

	if err := os.Setenv("XDG_CONFIG_HOME", dir); err != nil {
		t.Fatal(err)
	}

	os.Args = []string{"batsig", "alert", "set", "50", "override test"}
	if err := Execute(); err != nil {
		t.Fatalf("Execute() = %v", err)
	}

	message, err := readAlertMessage(50)
	if err != nil {
		t.Fatal(err)
	}
	if message != "override test" {
		t.Fatalf("unexpected alert message: %q", message)
	}
}
