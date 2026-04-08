package notify

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	soundEnv             = "BATSIG_NOTIFY_SOUND"
	audioEnv             = "BATSIG_NOTIFY_AUDIO"
	defaultSoundEvent    = "message"
	defaultSoundFileName = "notification.wav"
)

var (
	notifyCommand = exec.CommandContext
	lookupPath    = exec.LookPath
	stat          = os.Stat
)

func Send(ctx context.Context, title, body string) error {
	if body == "" {
		body = title
	}

	if err := sendNotification(ctx, title, body); err != nil {
		return err
	}

	soundName, ok := audioSoundSpec()
	if !ok {
		return nil
	}

	if err := playAudio(ctx, soundName, title, body); err != nil {
		return fmt.Errorf("audio notification failed: %w", err)
	}

	return nil
}

func sendNotification(ctx context.Context, title, body string) error {
	cmd := notifyCommand(ctx, "notify-send", title, body)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("notify-send failed: %w", err)
	}
	return nil
}

func audioSoundSpec() (string, bool) {
	if value := strings.TrimSpace(os.Getenv(soundEnv)); value != "" {
		return value, true
	}
	if truthy(os.Getenv(audioEnv)) {
		return defaultSoundEvent, true
	}
	if path, ok := configNotificationSoundPath(); ok {
		return path, true
	}
	return "", false
}

func configNotificationSoundPath() (string, bool) {
	path := filepath.Join(configDir(), defaultSoundFileName)
	if _, err := stat(path); err == nil {
		return path, true
	}
	return "", false
}

func configDir() string {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "batsig")
	}
	if home := os.Getenv("HOME"); home != "" {
		return filepath.Join(home, ".config", "batsig")
	}
	return filepath.Join(".config", "batsig")
}

func truthy(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}

func playAudio(ctx context.Context, soundName, title, body string) error {
	if filepath.IsAbs(soundName) || strings.ContainsRune(soundName, '/') || strings.HasPrefix(soundName, "./") || strings.HasPrefix(soundName, "../") {
		return playAudioFile(ctx, soundName)
	}
	return playAudioEvent(ctx, soundName, title, body)
}

func playAudioEvent(ctx context.Context, eventID, title, body string) error {
	command, err := lookupPath("canberra-gtk-play")
	if err != nil {
		return errors.New("canberra-gtk-play not found")
	}

	cmd := notifyCommand(ctx, command, "--id", eventID, "--description", body)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("canberra-gtk-play failed: %w", err)
	}
	return nil
}

func playAudioFile(ctx context.Context, path string) error {
	if _, err := stat(path); err != nil {
		return fmt.Errorf("sound file not found: %w", err)
	}

	if command, err := lookupPath("paplay"); err == nil {
		cmd := notifyCommand(ctx, command, path)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("paplay failed: %w", err)
		}
		return nil
	}

	if command, err := lookupPath("aplay"); err == nil {
		cmd := notifyCommand(ctx, command, path)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("aplay failed: %w", err)
		}
		return nil
	}

	return errors.New("audio playback command not found; install paplay or aplay")
}
