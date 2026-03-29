package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestWatcher_StartStop(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	if err := os.WriteFile(configPath, []byte(`server:\n  port: 8080\n`), 0644); err != nil {
		t.Fatalf("Failed to write initial config: %v", err)
	}

	watcher := NewConfigWatcher(configPath, func(newConfig *Config) {})

	if err := watcher.Start(); err != nil {
		t.Fatalf("Failed to start watcher: %v", err)
	}

	time.Sleep(1 * time.Second)

	watcher.Stop()
}

func TestWatcher_HotReload(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	if err := os.WriteFile(configPath, []byte(`server:\n  port: 8080\n`), 0644); err != nil {
		t.Fatalf("Failed to write initial config: %v", err)
	}

	var reloaded bool
	watcher := NewConfigWatcher(configPath, func(newConfig *Config) {
		if newConfig.Server.Port == 9090 {
			reloaded = true
		}
	})

	if err := watcher.Start(); err != nil {
		t.Fatalf("Failed to start watcher: %v", err)
	}

	time.Sleep(2 * time.Second)

	if err := os.WriteFile(configPath, []byte(`server:\n  port: 9090\n`), 0644); err != nil {
		t.Fatalf("Failed to write updated config: %v", err)
	}

	time.Sleep(5 * time.Second)

	watcher.Stop()

	if !reloaded {
		t.Error("Expected watcher callback to be triggered with new port 9090")
	}
}
