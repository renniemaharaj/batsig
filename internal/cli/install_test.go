package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindLocalResource(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "notification.wav")
	if err := os.WriteFile(file, []byte("data"), 0o644); err != nil {
		t.Fatal(err)
	}

	src, ok, err := findLocalResource("notification.wav", dir, "/tmp")
	if err != nil {
		t.Fatal(err)
	}
	if !ok || src != file {
		t.Fatalf("expected %q, got %q ok=%v", file, src, ok)
	}
}

func TestFindLocalDirectory(t *testing.T) {
	dir := t.TempDir()
	subdir := filepath.Join(dir, "alerts")
	if err := os.MkdirAll(subdir, 0o755); err != nil {
		t.Fatal(err)
	}

	src, ok, err := findLocalDirectory("alerts", dir, "/tmp")
	if err != nil {
		t.Fatal(err)
	}
	if !ok || src != subdir {
		t.Fatalf("expected %q, got %q ok=%v", subdir, src, ok)
	}
}
