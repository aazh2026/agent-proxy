## ADDED Requirements

### Requirement: YAML Configuration File
The system SHALL support YAML-based configuration file for all settings.

#### Scenario: Default config location
- **WHEN** system starts without --config flag
- **THEN** system reads configuration from `./agent-proxy.yaml` in current directory

#### Scenario: Custom config location
- **WHEN** system starts with `--config /path/to/config.yaml`
- **THEN** system reads configuration from specified path

#### Scenario: Config file format
- **WHEN** system reads configuration file
- **THEN** file must be valid YAML format

#### Scenario: Missing config file
- **WHEN** specified config file does not exist
- **WHEN** no environment variables provided
- **THEN** system starts with default configuration and warns user

### Requirement: Environment Variable Override
The system SHALL support environment variable overrides for configuration values.

#### Scenario: Environment variable format
- **WHEN** environment variable `AGENT_PROXY_SERVER_PORT=8080` is set
- **THEN** system uses port 8080, overriding config file value

#### Scenario: Nested config override
- **WHEN** environment variable `AGENT_PROXY_AUTH_METHOD=x-user-id` is set
- **THEN** system sets auth.method to "x-user-id"

#### Scenario: Environment variable priority
- **WHEN** both config file and environment variable define same setting
- **THEN** environment variable takes precedence

### Requirement: Command-Line Arguments
The system SHALL support command-line argument overrides for critical settings.

#### Scenario: CLI argument format
- **WHEN** system starts with `--port 8080`
- **THEN** system uses port 8080

#### Scenario: CLI priority
- **WHEN** CLI argument, environment variable, and config file all define same setting
- **THEN** CLI argument takes highest precedence

#### Scenario: Help flag
- **WHEN** system starts with `--help` or `-h`
- **THEN** system displays usage information and exits

#### Scenario: Version flag
- **WHEN** system starts with `--version` or `-v`
- **THEN** system displays version information and exits

### Requirement: Configuration Priority
The system SHALL apply configuration values in defined priority order.

#### Scenario: Priority order
- **WHEN** configuration value defined at multiple levels
- **THEN** priority order: CLI arguments > Environment variables > Config file > Default values

#### Scenario: Partial override
- **WHEN** environment variable overrides only one config field
- **THEN** other fields use config file or default values

### Requirement: Configuration Validation
The system SHALL validate configuration on startup.

#### Scenario: Valid configuration
- **WHEN** configuration passes validation
- **THEN** system starts normally

#### Scenario: Invalid port
- **WHEN** configuration specifies port outside valid range (1-65535)
- **THEN** system fails to start with error message

#### Scenario: Invalid auth method
- **WHEN** configuration specifies unsupported auth method
- **THEN** system fails to start with error message

#### Scenario: Invalid encryption key
- **WHEN** configuration specifies encryption key with insufficient length
- **THEN** system fails to start with error message

#### Scenario: Missing required fields
- **WHEN** configuration missing required fields
- **THEN** system fails to start with validation error listing missing fields

### Requirement: Configuration Hot Reload
The system SHALL support hot-reloading configuration without restart.

#### Scenario: File watcher
- **WHEN** configuration file changes on disk
- **THEN** system detects change and reloads configuration

#### Scenario: Validation on reload
- **WHEN** new configuration is invalid
- **WHEN** reload attempted
- **THEN** system rejects new configuration, keeps running with old configuration

#### Scenario: Atomic reload
- **WHEN** configuration reloads
- **THEN** system applies all changes atomically (no partial state)

#### Scenario: Reload notification
- **WHEN** configuration successfully reloads
- **THEN** system logs reload event

### Requirement: Configuration Sections
The system SHALL organize configuration into logical sections.

#### Scenario: Server configuration
- **WHEN** configuration includes server section
- **THEN** system configures: host, port, read_timeout, write_timeout, max_connections

#### Scenario: Auth configuration
- **WHEN** configuration includes auth section
- **THEN** system configures: method, provider-specific settings

#### Scenario: Token configuration
- **WHEN** configuration includes token section
- **THEN** system configures: encryption_key, storage_path, auto_refresh_threshold

#### Scenario: Routing configuration
- **WHEN** configuration includes routing section
- **THEN** system configures: model_mappings, token_strategy, fallback_rules, retry_policy

#### Scenario: Observability configuration
- **WHEN** configuration includes observability section
- **THEN** system configures: metrics_enabled, log_level, admin_password, lan_access

### Requirement: Configuration Documentation
The system SHALL provide configuration schema and examples.

#### Scenario: Example configuration
- **WHEN** user needs configuration reference
- **THEN** system includes documented example config file with all options

#### Scenario: Schema validation
- **WHEN** configuration file loaded
- **THEN** system validates against schema, reports unknown fields

#### Scenario: Configuration help
- **WHEN** system starts with `--config-help`
- **THEN** system displays all configuration options with descriptions

### Requirement: Sensitive Configuration Handling
The system SHALL handle sensitive configuration values securely.

#### Scenario: Encryption key masking
- **WHEN** configuration displayed in logs or UI
- **THEN** system masks encryption_key value

#### Scenario: Password masking
- **WHEN** configuration displayed
- **THEN** system masks password values

#### Scenario: Token exclusion from config
- **WHEN** tokens stored in configuration
- **THEN** system warns user to use encrypted token storage instead

### Requirement: Configuration Backup and Restore
The system SHALL support configuration backup and restore.

#### Scenario: Configuration export
- **WHEN** user requests config export
- **THEN** system generates YAML file with current configuration (excluding secrets)

#### Scenario: Configuration import
- **WHEN** user imports configuration file
- **WHEN** file is valid
- **THEN** system applies imported configuration

#### Scenario: Secrets exclusion
- **WHEN** configuration exported
- **THEN** sensitive values (passwords, keys) are excluded or masked
