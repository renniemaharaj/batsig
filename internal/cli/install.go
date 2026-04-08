package cli

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

const (
	installSoundName  = "notification.wav"
	installAlertsName = "alerts"
	installConfigDir  = "batsig"
	defaultUserBinDir = ".local/bin"
)

func runInstall() error {
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("unable to determine executable path: %w", err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("unable to determine current working directory: %w", err)
	}

	if err := installExecutable(execPath); err != nil {
		return err
	}

	if err := installNotificationSound(cwd, filepath.Dir(execPath)); err != nil {
		return err
	}

	if err := installAlerts(cwd, filepath.Dir(execPath)); err != nil {
		return err
	}

	return nil
}

func installExecutable(src string) error {
	binDir, err := userBinDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(binDir, 0o755); err != nil {
		return fmt.Errorf("unable to create bin directory: %w", err)
	}

	dest := filepath.Join(binDir, filepath.Base(src))
	if same, err := sameFile(src, dest); err == nil && same {
		return nil
	}

	return copyFile(src, dest, 0o755)
}

func installNotificationSound(cwd, execDir string) error {
	src, ok, err := findLocalResource(installSoundName, cwd, execDir)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	destDir := filepath.Join(configDir(), installConfigDir)
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return fmt.Errorf("unable to create config directory: %w", err)
	}

	dest := filepath.Join(destDir, installSoundName)
	if same, err := sameFile(src, dest); err == nil && same {
		return nil
	}

	return copyFile(src, dest, 0o644)
}

func installAlerts(cwd, execDir string) error {
	srcDir, ok, err := findLocalDirectory(installAlertsName, cwd, execDir)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	destDir := alertsDir()
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return fmt.Errorf("unable to create alerts directory: %w", err)
	}

	return copyAlertTree(srcDir, destDir)
}

func copyAlertTree(srcDir, destDir string) error {
	return filepath.WalkDir(srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(destDir, relPath)
		if d.IsDir() {
			return os.MkdirAll(destPath, 0o755)
		}

		if destInfo, err := os.Stat(destPath); err == nil && destInfo.Mode().IsDir() {
			return fmt.Errorf("destination path %q is a directory", destPath)
		}

		if same, err := sameFile(path, destPath); err == nil && same {
			return nil
		}
		return copyFile(path, destPath, 0o644)
	})
}

func userBinDir() (string, error) {
	if xdg := os.Getenv("XDG_BIN_HOME"); xdg != "" {
		return xdg, nil
	}
	if home := os.Getenv("HOME"); home != "" {
		return filepath.Join(home, defaultUserBinDir), nil
	}
	return "", fmt.Errorf("HOME not set")
}

func findLocalResource(name, cwd, execDir string) (string, bool, error) {
	paths := []string{
		filepath.Join(cwd, name),
		filepath.Join(cwd, "assets", name),
		filepath.Join(execDir, name),
		filepath.Join(execDir, "assets", name),
	}

	for _, path := range paths {
		if info, err := os.Stat(path); err == nil && !info.IsDir() {
			return path, true, nil
		}
	}

	return "", false, nil
}

func findLocalDirectory(name, cwd, execDir string) (string, bool, error) {
	paths := []string{
		filepath.Join(cwd, name),
		filepath.Join(cwd, "assets", name),
		filepath.Join(execDir, name),
		filepath.Join(execDir, "assets", name),
	}

	for _, path := range paths {
		if info, err := os.Stat(path); err == nil && info.IsDir() {
			return path, true, nil
		}
	}

	return "", false, nil
}

func copyFile(src, dest string, perm os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("unable to open source file %q: %w", src, err)
	}
	defer in.Close()

	out, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, perm)
	if err != nil {
		return fmt.Errorf("unable to create destination file %q: %w", dest, err)
	}
	defer func() {
		_ = out.Close()
	}()

	if _, err := io.Copy(out, in); err != nil {
		return fmt.Errorf("unable to copy file %q to %q: %w", src, dest, err)
	}

	if err := out.Sync(); err != nil {
		return fmt.Errorf("unable to sync destination file %q: %w", dest, err)
	}

	if err := os.Chmod(dest, perm); err != nil {
		return fmt.Errorf("unable to set permission on %q: %w", dest, err)
	}

	return nil
}

func sameFile(a, b string) (bool, error) {
	aInfo, err := os.Stat(a)
	if err != nil {
		return false, err
	}
	bInfo, err := os.Stat(b)
	if err != nil {
		return false, err
	}
	if aInfo.Size() != bInfo.Size() {
		return false, nil
	}

	aData, err := os.ReadFile(a)
	if err != nil {
		return false, err
	}
	bData, err := os.ReadFile(b)
	if err != nil {
		return false, err
	}
	return string(aData) == string(bData), nil
}
