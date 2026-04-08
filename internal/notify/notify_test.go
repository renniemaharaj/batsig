package notify

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestAudioSoundSpec(t *testing.T) {
	t.Setenv(soundEnv, "ding")
	t.Setenv(audioEnv, "1")

	value, ok := audioSoundSpec()
	if !ok || value != "ding" {
		t.Fatalf("expected sound env to win, got %q ok=%v", value, ok)
	}
}

func TestAudioSoundSpecFallback(t *testing.T) {
	t.Setenv(soundEnv, "")
	t.Setenv(audioEnv, "yes")

	value, ok := audioSoundSpec()
	if !ok || value != defaultSoundEvent {
		t.Fatalf("expected default sound event, got %q ok=%v", value, ok)
	}
}

func TestAudioSoundSpecDisabled(t *testing.T) {
	t.Setenv(soundEnv, "")
	t.Setenv(audioEnv, "0")

	_, ok := audioSoundSpec()
	if ok {
		t.Fatal("expected audio disabled")
	}
}

func TestAudioSoundSpecUsesConfigFile(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)

	path := filepath.Join(dir, "batsig", defaultSoundFileName)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("dummy"), 0o644); err != nil {
		t.Fatal(err)
	}

	value, ok := audioSoundSpec()
	if !ok || value != path {
		t.Fatalf("expected config sound path %q, got %q ok=%v", path, value, ok)
	}
}

func TestSend_NoAudio(t *testing.T) {
	t.Setenv(soundEnv, "")
	t.Setenv(audioEnv, "0")

	savedNotify := notifyCommand
	defer func() { notifyCommand = savedNotify }()
	notifyCommand = func(ctx context.Context, name string, args ...string) *exec.Cmd {
		return exec.CommandContext(ctx, "true")
	}

	if err := Send(context.Background(), "Test", "Body"); err != nil {
		t.Fatalf("Send() returned error: %v", err)
	}
}

func TestSend_WithAudioEvent(t *testing.T) {
	t.Setenv(soundEnv, "message")

	savedNotify := notifyCommand
	savedLookup := lookupPath
	defer func() {
		notifyCommand = savedNotify
		lookupPath = savedLookup
	}()

	lookupPath = func(file string) (string, error) {
		if file == "canberra-gtk-play" {
			return "true", nil
		}
		return "", os.ErrNotExist
	}

	notifyCommand = func(ctx context.Context, name string, args ...string) *exec.Cmd {
		return exec.CommandContext(ctx, "true")
	}

	if err := Send(context.Background(), "Test", "Body"); err != nil {
		t.Fatalf("Send() returned error: %v", err)
	}
}
