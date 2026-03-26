## ADDED Requirements

### Requirement: Token Data Structure
The system SHALL store tokens with comprehensive metadata for lifecycle management and access control.

#### Scenario: Token fields
- **WHEN** system stores a token
- **THEN** token includes: token_id (unique), user_id (owner), provider (openai/anthropic/gemini), type (api_key/oauth), access_token (encrypted), refresh_token (encrypted, optional), expires_at, created_at, updated_at, status (enabled/disabled), priority, allowed_models

#### Scenario: Token ID generation
- **WHEN** system creates new token
- **THEN** system generates unique token_id using UUID v4 or equivalent

#### Scenario: Token timestamps
- **WHEN** system creates or updates token
- **THEN** system sets created_at on creation, updates updated_at on every modification

### Requirement: Token CRUD Operations
The system SHALL support complete token lifecycle management: create, read, update, delete.

#### Scenario: Token creation
- **WHEN** user submits new token with valid provider, access_token, and optional metadata
- **THEN** system encrypts token, stores in database, returns token_id

#### Scenario: Token retrieval by user
- **WHEN** user requests their tokens
- **THEN** system returns list of tokens owned by user with metadata (access_token masked)

#### Scenario: Token update
- **WHEN** user updates token metadata (status, priority, allowed_models)
- **THEN** system updates token record and sets updated_at timestamp

#### Scenario: Token deletion
- **WHEN** user deletes token
- **THEN** system removes token from database permanently

#### Scenario: Cross-user token isolation
- **WHEN** user A requests tokens
- **THEN** system returns only tokens where user_id matches user A, never user B's tokens

### Requirement: Token Encryption at Rest
The system SHALL encrypt all token values using AES-256-GCM before persisting to storage.

#### Scenario: Token encryption on write
- **WHEN** system stores token
- **THEN** system encrypts access_token and refresh_token using AES-256-GCM with configured encryption key

#### Scenario: Token decryption on read
- **WHEN** system needs to use token for request
- **WHEN** system decrypts token in memory
- **THEN** system uses for request and immediately clears plaintext from memory

#### Scenario: Encryption key configuration
- **WHEN** system starts with `encryption_key` in configuration
- **THEN** system uses provided key for AES-256-GCM encryption/decryption

#### Scenario: Missing encryption key
- **WHEN** system starts without encryption_key configured
- **THEN** system generates random key and warns user to persist it for data recovery

#### Scenario: Invalid encryption key
- **WHEN** system attempts to decrypt with wrong key
- **THEN** system returns error indicating decryption failed, token unusable

### Requirement: Token Residency Security
The system SHALL enforce token residency—tokens never leave proxy boundary under any circumstances.

#### Scenario: Token never in response
- **WHEN** system returns any response to client (success, error, metadata)
- **THEN** response never contains token value in body, headers, or cookies

#### Scenario: Token never in logs
- **WHEN** system writes logs
- **THEN** logs never contain full token value, only masked version (first 4 + last 4 chars)

#### Scenario: Token never in metrics
- **WHEN** system exposes metrics
- **THEN** metrics never contain token values

#### Scenario: Token memory clearing
- **WHEN** system decrypts token for request injection
- **THEN** system immediately overwrites plaintext memory with zeros after use

#### Scenario: Token in crash dumps
- **WHEN** system crashes
- **THEN** core dumps or error reports never contain decrypted token values (memory clearing before dump)

### Requirement: Token Status Management
The system SHALL support enabling and disabling tokens without deletion.

#### Scenario: Disabled token skipped
- **WHEN** token status is "disabled"
- **WHEN** system selects token for request
- **THEN** system skips disabled token and selects next available token

#### Scenario: Enable token
- **WHEN** user enables previously disabled token
- **THEN** system updates status to "enabled", token becomes available for requests

#### Scenario: Disable all tokens error
- **WHEN** user has no enabled tokens for required provider
- **WHEN** request requires that provider
- **THEN** system returns error indicating no available tokens

### Requirement: Token Model Permissions
The system SHALL support per-token model access control.

#### Scenario: Token with allowed_models
- **WHEN** token has `allowed_models: ["gpt-4", "gpt-3.5-turbo"]`
- **WHEN** request is for model "gpt-4"
- **THEN** system allows token to be used

#### Scenario: Token model restriction
- **WHEN** token has `allowed_models: ["gpt-4"]`
- **WHEN** request is for model "gpt-3.5-turbo"
- **THEN** system skips this token, selects another token or returns error

#### Scenario: Empty allowed_models
- **WHEN** token has `allowed_models: []` (empty array)
- **THEN** system allows token to be used for any model from matching provider

### Requirement: Token Auto-Refresh
The system SHALL support automatic refresh of OAuth tokens before expiration.

#### Scenario: Proactive refresh trigger
- **WHEN** OAuth token expires_at is within configurable threshold (default 30 minutes)
- **WHEN** system processes request needing this token
- **THEN** system automatically refreshes token using refresh_token before using it

#### Scenario: Refresh success
- **WHEN** token refresh succeeds
- **WHEN** provider returns new access_token
- **THEN** system updates token record with new access_token and new expires_at

#### Scenario: Refresh failure
- **WHEN** token refresh fails (invalid refresh_token, provider error)
- **THEN** system disables token, logs error, selects next available token

#### Scenario: Refresh retry
- **WHEN** token refresh fails with transient error
- **THEN** system retries refresh up to configured max_retries before disabling

### Requirement: Token Storage Backend
The system SHALL persist tokens in SQLite database with encryption.

#### Scenario: SQLite persistence
- **WHEN** system stores tokens
- **THEN** tokens persisted in SQLite file at configured path (default: `~/.agent-proxy/tokens.db`)

#### Scenario: Database migration
- **WHEN** system starts with newer version
- **WHEN** database schema needs update
- **THEN** system automatically migrates schema without data loss
