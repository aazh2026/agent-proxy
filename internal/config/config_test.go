package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Server.Host != "127.0.0.1" {
		t.Errorf("Expected default host 127.0.0.1, got %s", cfg.Server.Host)
	}

	if cfg.Server.Port != 4000 {
		t.Errorf("Expected default port 4000, got %d", cfg.Server.Port)
	}

	if cfg.Auth.Method != "x-user-id" {
		t.Errorf("Expected default auth method x-user-id, got %s", cfg.Auth.Method)
	}

	if cfg.Token.AutoRefreshMinutes != 30 {
		t.Errorf("Expected default auto refresh 30 minutes, got %d", cfg.Token.AutoRefreshMinutes)
	}

	if cfg.Logging.Level != "info" {
		t.Errorf("Expected default log level info, got %s", cfg.Logging.Level)
	}
}

func TestLoad_NoFile(t *testing.T) {
	// Test with non-existent config file
	cfg, err := Load("/non/existent/config.yaml")
	if err != nil {
		t.Fatalf("Load should not fail for non-existent file: %v", err)
	}

	if cfg.Server.Host != "127.0.0.1" {
		t.Errorf("Expected default host 127.0.0.1 when file missing, got %s", cfg.Server.Host)
	}
}

func TestLoad_YAMLFile(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Write test config
	configContent := `
server:
  host: "0.0.0.0"
  port: 8080
auth:
  method: "local"
token:
  encryption_key: "test-key-32-bytes-long!!"
  auto_refresh_minutes: 10
logging:
  level: "debug"
`

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.Server.Host != "0.0.0.0" {
		t.Errorf("Expected host 0.0.0.0, got %s", cfg.Server.Host)
	}

	if cfg.Server.Port != 8080 {
		t.Errorf("Expected port 8080, got %d", cfg.Server.Port)
	}

	if cfg.Auth.Method != "local" {
		t.Errorf("Expected auth method local, got %s", cfg.Auth.Method)
	}

	if cfg.Token.EncryptionKey != "test-key-32-bytes-long!!" {
		t.Errorf("Expected encryption key test-key-32-bytes-long!!, got %s", cfg.Token.EncryptionKey)
	}

	if cfg.Token.AutoRefreshMinutes != 10 {
		t.Errorf("Expected auto refresh 10 minutes, got %d", cfg.Token.AutoRefreshMinutes)
	}

	if cfg.Logging.Level != "debug" {
		t.Errorf("Expected log level debug, got %s", cfg.Logging.Level)
	}
}

func TestLoad_EnvOverride(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Write test config with port 8080
	configContent := `
server:
  port: 8080
`

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Set environment variable for port 9090
	os.Setenv("AGENT_PROXY_SERVER_PORT", "9090")
	defer os.Unsetenv("AGENT_PROXY_SERVER_PORT")

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.Server.Port != 9090 {
		t.Errorf("Expected port 9090 from env var, got %d", cfg.Server.Port)
	}
}

func TestLoad_CLIOverride(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Write test config with port 8080
	configContent := `
server:
  port: 8080
`

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Temporarily replace os.Args to simulate CLI args
	oldArgs := os.Args
	os.Args = []string{"agent-proxy", "--port", "7070"}
	defer func() { os.Args = oldArgs }()

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.Server.Port != 7070 {
		t.Errorf("Expected port 7070 from CLI, got %d", cfg.Server.Port)
	}
}

func TestLoad_Validation_InvalidPort(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Write test config with invalid port
	configContent := `
server:
  port: 70000
`

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	_, err := Load(configPath)
	if err == nil {
		t.Error("Expected error for invalid port 70000, got nil")
	}
}

func TestLoad_Validation_InvalidAuthMethod(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Write test config with invalid auth method
	configContent := `
auth:
  method: "invalid-method"
`

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	_, err := Load(configPath)
	if err == nil {
		t.Error("Expected error for invalid auth method, got nil")
	}
}

func TestLoad_Validation_InvalidLogLevel(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Write test config with invalid log level
	configContent := `
logging:
  level: "invalid-level"
`

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	_, err := Load(configPath)
	if err == nil {
		t.Error("Expected error for invalid log level, got nil")
	}
}
