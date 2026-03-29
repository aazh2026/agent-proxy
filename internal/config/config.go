package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server    ServerConfig    `yaml:"server"`
	Auth      AuthConfig      `yaml:"auth"`
	Token     TokenConfig     `yaml:"token"`
	Routing   RoutingConfig   `yaml:"routing"`
	Providers ProvidersConfig `yaml:"providers"`
	Admin     AdminConfig     `yaml:"admin"`
	Logging   LoggingConfig   `yaml:"logging"`
}

type ServerConfig struct {
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	TLSEnabled  bool   `yaml:"tls_enabled"`
	TLSCertFile string `yaml:"tls_cert_file"`
	TLSKeyFile  string `yaml:"tls_key_file"`
	ReadTimeout int    `yaml:"read_timeout"`
	IdleTimeout int    `yaml:"idle_timeout"`
	MaxConns    int    `yaml:"max_connections"`
}

type AuthConfig struct {
	Method         string   `yaml:"method"`
	AllowedUserIDs []string `yaml:"allowed_user_ids"`
	AllowedDomains []string `yaml:"allowed_domains"`
}

type TokenConfig struct {
	EncryptionKey      string `yaml:"encryption_key"`
	StoragePath        string `yaml:"storage_path"`
	AutoRefreshMinutes int    `yaml:"auto_refresh_minutes"`
}

type RoutingConfig struct {
	TokenStrategy string                  `yaml:"token_strategy"`
	ModelMappings map[string]ModelMapping `yaml:"model_mappings"`
	RetryPolicy   RetryPolicyConfig       `yaml:"retry_policy"`
}

type ModelMapping struct {
	Provider string `yaml:"provider"`
	Model    string `yaml:"model"`
}

type RetryPolicyConfig struct {
	MaxRetries int `yaml:"max_retries"`
	InitialMs  int `yaml:"initial_ms"`
	MaxDelayMs int `yaml:"max_delay_ms"`
}

type AdminConfig struct {
	Password  string `yaml:"password"`
	LANAccess bool   `yaml:"lan_access"`
}

type LoggingConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

type ProvidersConfig struct {
	OpenAI    ProviderConfig `yaml:"openai"`
	Anthropic ProviderConfig `yaml:"anthropic"`
	Google    ProviderConfig `yaml:"google"`
}

type ProviderConfig struct {
	Enabled           bool   `yaml:"enabled"`
	BaseURL           string `yaml:"base_url"`
	Timeout           int    `yaml:"timeout_seconds"`
	MaxConns          int    `yaml:"max_connections"`
	OAuthClientID     string `yaml:"oauth_client_id"`
	OAuthClientSecret string `yaml:"oauth_client_secret"`
	OAuthRedirectURI  string `yaml:"oauth_redirect_uri"`
}

func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Host:        "127.0.0.1",
			Port:        4000,
			ReadTimeout: 30,
			IdleTimeout: 120,
			MaxConns:    1000,
		},
		Auth: AuthConfig{
			Method: "x-user-id",
		},
		Token: TokenConfig{
			StoragePath:        "agent-proxy.db",
			AutoRefreshMinutes: 30,
		},
		Routing: RoutingConfig{
			TokenStrategy: "round-robin",
			RetryPolicy: RetryPolicyConfig{
				MaxRetries: 3,
				InitialMs:  100,
				MaxDelayMs: 5000,
			},
		},
		Logging: LoggingConfig{
			Level:  "info",
			Format: "text",
		},
		Providers: ProvidersConfig{
			OpenAI: ProviderConfig{
				Enabled:  true,
				BaseURL:  "https://api.openai.com/v1",
				Timeout:  60,
				MaxConns: 100,
			},
			Anthropic: ProviderConfig{
				Enabled:  true,
				BaseURL:  "https://api.anthropic.com/v1",
				Timeout:  60,
				MaxConns: 100,
			},
			Google: ProviderConfig{
				Enabled:  true,
				BaseURL:  "https://generativelanguage.googleapis.com/v1beta",
				Timeout:  60,
				MaxConns: 100,
			},
		},
	}
}

func Load(configPath string) (*Config, error) {
	config := DefaultConfig()

	if err := loadFromFile(config, configPath); err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to load config file: %w", err)
		}
	}

	applyEnvOverrides(config)
	applyCLIOverrides(config)

	if err := validate(config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

func loadFromFile(config *Config, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(data, config); err != nil {
		return fmt.Errorf("failed to parse YAML: %w", err)
	}

	return nil
}

func applyEnvOverrides(config *Config) {
	if v := os.Getenv("AGENT_PROXY_SERVER_HOST"); v != "" {
		config.Server.Host = v
	}
	if v := os.Getenv("AGENT_PROXY_SERVER_PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			config.Server.Port = port
		}
	}
	if v := os.Getenv("AGENT_PROXY_AUTH_METHOD"); v != "" {
		config.Auth.Method = v
	}
	if v := os.Getenv("AGENT_PROXY_TOKEN_ENCRYPTION_KEY"); v != "" {
		config.Token.EncryptionKey = v
	}
	if v := os.Getenv("AGENT_PROXY_TOKEN_STORAGE_PATH"); v != "" {
		config.Token.StoragePath = v
	}
	if v := os.Getenv("AGENT_PROXY_ADMIN_PASSWORD"); v != "" {
		config.Admin.Password = v
	}
	if v := os.Getenv("AGENT_PROXY_LOG_LEVEL"); v != "" {
		config.Logging.Level = v
	}
}

func applyCLIOverrides(config *Config) {
	for i, arg := range os.Args {
		if i == 0 {
			continue
		}

		switch {
		case arg == "--port" && i+1 < len(os.Args):
			if port, err := strconv.Atoi(os.Args[i+1]); err == nil {
				config.Server.Port = port
			}
		case strings.HasPrefix(arg, "--port="):
			if port, err := strconv.Atoi(strings.TrimPrefix(arg, "--port=")); err == nil {
				config.Server.Port = port
			}
		case arg == "--host" && i+1 < len(os.Args):
			config.Server.Host = os.Args[i+1]
		case strings.HasPrefix(arg, "--host="):
			config.Server.Host = strings.TrimPrefix(arg, "--host=")
		case arg == "--auth-method" && i+1 < len(os.Args):
			config.Auth.Method = os.Args[i+1]
		case strings.HasPrefix(arg, "--auth-method="):
			config.Auth.Method = strings.TrimPrefix(arg, "--auth-method=")
		case arg == "--log-level" && i+1 < len(os.Args):
			config.Logging.Level = os.Args[i+1]
		case strings.HasPrefix(arg, "--log-level="):
			config.Logging.Level = strings.TrimPrefix(arg, "--log-level=")
		}
	}
}

func validate(config *Config) error {
	if config.Server.Port < 1 || config.Server.Port > 65535 {
		return fmt.Errorf("invalid port: %d", config.Server.Port)
	}

	validAuthMethods := map[string]bool{
		"x-user-id": true,
		"local":     true,
		"oidc":      true,
		"oauth2":    true,
		"session":   true,
	}
	if !validAuthMethods[config.Auth.Method] {
		return fmt.Errorf("invalid auth method: %s", config.Auth.Method)
	}

	validStrategies := map[string]bool{
		"round-robin": true,
		"weighted":    true,
		"priority":    true,
	}
	if !validStrategies[config.Routing.TokenStrategy] {
		return fmt.Errorf("invalid token strategy: %s", config.Routing.TokenStrategy)
	}

	validLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	if !validLevels[config.Logging.Level] {
		return fmt.Errorf("invalid log level: %s", config.Logging.Level)
	}

	return nil
}

// Validate checks if the provided config values are valid.
func Validate(config *Config) error {
	return validate(config)
}

// Save writes the configuration to the specified path in YAML format.
func Save(configPath string, config *Config) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, data, 0644)
}
