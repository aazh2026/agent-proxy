## ADDED Requirements

### Requirement: Request Logging
The system SHALL log all requests for debugging and audit.

#### Scenario: Request logged
- **WHEN** request completes (success or failure)
- **THEN** system creates log entry with request details

#### Scenario: Log entry fields
- **WHEN** log entry created
- **THEN** entry includes: request_id, timestamp, user_id, model, provider, status_code, latency_ms, error_message

### Requirement: Ring Buffer Storage
The system SHALL use ring buffer for log storage.

#### Scenario: Ring buffer capacity
- **WHEN** system stores logs
- **THEN** system uses ring buffer with configurable capacity (default 1000 entries)

#### Scenario: Old entries overwritten
- **WHEN** ring buffer is full
- **THEN** oldest entries are overwritten

### Requirement: Log Querying
The system SHALL support querying recent logs.

#### Scenario: Query all logs
- **WHEN** admin queries logs
- **THEN** system returns most recent log entries

#### Scenario: Query by user
- **WHEN** admin queries logs for specific user
- **THEN** system returns logs for that user only

#### Scenario: Query by model
- **WHEN** admin queries logs for specific model
- **THEN** system returns logs for that model only

### Requirement: Token Masking in Logs
The system SHALL mask sensitive data in logs.

#### Scenario: Token masked
- **WHEN** log entry contains token reference
- **THEN** system logs only masked token (first 4 + last 4 chars)
