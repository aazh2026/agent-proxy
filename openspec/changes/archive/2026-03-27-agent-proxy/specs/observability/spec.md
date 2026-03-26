## ADDED Requirements

### Requirement: Web Dashboard Access
The system SHALL provide a web-based admin dashboard at `/admin` path.

#### Scenario: Dashboard accessibility
- **WHEN** user navigates to `http://localhost:{port}/admin`
- **THEN** system serves web dashboard interface

#### Scenario: Localhost-only access
- **WHEN** system starts with default configuration
- **THEN** dashboard only accessible from localhost (127.0.0.1)

#### Scenario: LAN access configuration
- **WHEN** configuration enables `admin.lan_access: true`
- **THEN** dashboard accessible from LAN addresses

#### Scenario: Dashboard authentication
- **WHEN** configuration specifies `admin.password`
- **THEN** dashboard requires password authentication before access

### Requirement: Real-Time Metrics Dashboard
The system SHALL display real-time performance metrics on dashboard.

#### Scenario: QPS display
- **WHEN** dashboard loads
- **THEN** system displays current queries per second with 1-second refresh

#### Scenario: Latency metrics
- **WHEN** dashboard loads
- **THEN** system displays average latency and P99 latency

#### Scenario: Success rate
- **WHEN** dashboard loads
- **THEN** system displays request success rate percentage

#### Scenario: Error rate
- **WHEN** dashboard loads
- **THEN** system displays error rate by error type

#### Scenario: Provider distribution
- **WHEN** dashboard loads
- **THEN** system displays request distribution across providers

#### Scenario: Time range selection
- **WHEN** user selects time range (1min, 5min, 1hour)
- **THEN** dashboard updates metrics for selected time window

### Requirement: Request Logging
The system SHALL log recent requests for debugging and audit.

#### Scenario: Request log entry
- **WHEN** request completes (success or failure)
- **WHEN** logging enabled
- **THEN** system creates log entry with: request_id, timestamp, user_id, model, provider, status_code, latency_ms, error_message

#### Scenario: Ring buffer storage
- **WHEN** system stores request logs
- **THEN** system uses ring buffer with configurable capacity (default 1000 entries)

#### Scenario: Log query
- **WHEN** user queries request logs
- **WHEN** filters applied (user, model, status, time range)
- **THEN** system returns filtered log entries

#### Scenario: Token masking in logs
- **WHEN** log entry contains token reference
- **THEN** system logs only masked token (first 4 + last 4 chars)

#### Scenario: Log persistence
- **WHEN** system restarts
- **THEN** request logs cleared (ring buffer in-memory only)

### Requirement: User Usage Statistics
The system SHALL provide per-user usage statistics on dashboard.

#### Scenario: User usage display
- **WHEN** admin views user usage page
- **THEN** system displays: request count, success/failure ratio, average latency, token consumption, estimated cost

#### Scenario: Usage filtering
- **WHEN** admin filters by time range, provider, or model
- **THEN** system updates usage statistics for filtered criteria

#### Scenario: Usage export
- **WHEN** admin exports usage data
- **WHEN** export format selected (CSV, JSON)
- **THEN** system generates downloadable file with usage statistics

### Requirement: Token Management UI
The system SHALL provide visual interface for token management.

#### Scenario: Token list display
- **WHEN** admin views token management page
- **THEN** system displays list of tokens with metadata (provider, status, created date)

#### Scenario: Token creation
- **WHEN** admin creates token via UI
- **WHEN** form submitted with provider, access_token, and metadata
- **THEN** system creates token and updates list

#### Scenario: Token editing
- **WHEN** admin edits token via UI
- **WHEN** changes saved
- **THEN** system updates token metadata

#### Scenario: Token enable/disable
- **WHEN** admin toggles token status
- **THEN** system updates token status immediately

#### Scenario: Token deletion
- **WHEN** admin deletes token via UI
- **WHEN** deletion confirmed
- **THEN** system removes token from database

#### Scenario: Token masking
- **WHEN** UI displays token
- **THEN** system shows masked value (first 4 + last 4 chars), not full token

### Requirement: Configuration Management UI
The system SHALL provide visual interface for configuration management.

#### Scenario: Configuration display
- **WHEN** admin views configuration page
- **THEN** system displays current configuration in editable format

#### Scenario: Configuration editing
- **WHEN** admin edits configuration via UI
- **WHEN** changes saved and validated
- **THEN** system applies configuration without restart

#### Scenario: Configuration validation
- **WHEN** admin saves invalid configuration
- **THEN** system displays validation errors, does not apply changes

#### Scenario: Configuration backup
- **WHEN** admin exports configuration
- **THEN** system generates downloadable YAML file

### Requirement: Metrics Collection
The system SHALL collect and aggregate performance metrics.

#### Scenario: Request metrics
- **WHEN** request completes
- **THEN** system records: latency, status code, provider, model, user

#### Scenario: Metric aggregation
- **WHEN** metrics collected
- **THEN** system aggregates into time buckets (1min, 5min, 1hour)

#### Scenario: Metric storage
- **WHEN** system stores metrics
- **THEN** metrics stored in-memory, reset on restart

#### Scenario: Prometheus metrics endpoint
- **WHEN** system configured with `metrics.prometheus_enabled: true`
- **WHEN** request to `/metrics`
- **THEN** system returns metrics in Prometheus exposition format

### Requirement: Dashboard Security
The system SHALL protect dashboard access with authentication when configured.

#### Scenario: Password protection
- **WHEN** configuration specifies `admin.password`
- **WHEN** user accesses dashboard
- **THEN** system requires password entry before showing dashboard

#### Scenario: Session management
- **WHEN** user authenticates to dashboard
- **THEN** system creates session with configurable timeout

#### Scenario: Failed authentication
- **WHEN** user enters incorrect password
- **THEN** system denies access, logs attempt

#### Scenario: No authentication required
- **WHEN** configuration does not specify admin.password
- **THEN** dashboard accessible without authentication (localhost-only by default)
