## ADDED Requirements

### Requirement: User Token Isolation
The system SHALL enforce strict token isolation between users.

#### Scenario: Token ownership
- **WHEN** user creates token
- **THEN** token is associated with user's user_id

#### Scenario: Token access restriction
- **WHEN** user A requests their tokens
- **THEN** system returns only tokens where user_id matches user A

#### Scenario: Cross-user token prevention
- **WHEN** user A makes request
- **THEN** system only considers tokens where user_id matches user A, never user B's tokens

#### Scenario: Token deletion isolation
- **WHEN** user A deletes token
- **THEN** system only deletes if token user_id matches user A

### Requirement: Model Access Control
The system SHALL support per-user model access restrictions.

#### Scenario: Allowed models configuration
- **WHEN** configuration specifies user "alice" can access ["gpt-4", "gpt-3.5-turbo"]
- **WHEN** alice requests "gpt-4"
- **THEN** system allows request

#### Scenario: Denied model access
- **WHEN** configuration specifies user "alice" can access ["gpt-4"]
- **WHEN** alice requests "claude-3-opus"
- **THEN** system returns 403 Forbidden error

#### Scenario: Default allowed models
- **WHEN** configuration specifies global `default_allowed_models: ["gpt-3.5-turbo"]`
- **WHEN** user has no specific allowed_models configured
- **THEN** system uses default_allowed_models

#### Scenario: User override global
- **WHEN** user has specific allowed_models configured
- **THEN** user's list overrides global default_allowed_models

### Requirement: Usage Quotas
The system SHALL enforce per-user usage quotas.

#### Scenario: Request count quota
- **WHEN** configuration specifies user quota: 1000 requests per day
- **WHEN** user exceeds 1000 requests in 24 hours
- **THEN** system returns 429 error indicating quota exceeded

#### Scenario: Token consumption quota
- **WHEN** configuration specifies user quota: 1M tokens per month
- **WHEN** user exceeds token consumption limit
- **THEN** system returns 429 error indicating token quota exceeded

#### Scenario: Cost quota
- **WHEN** configuration specifies user quota: $100 per month
- **WHEN** estimated cost exceeds limit
- **THEN** system returns 429 error indicating cost quota exceeded

#### Scenario: Quota period reset
- **WHEN** quota period expires (daily/monthly)
- **THEN** system resets usage counters for new period

#### Scenario: Quota warning threshold
- **WHEN** user usage reaches 80% of quota (configurable threshold)
- **THEN** system logs warning (notification system future enhancement)

### Requirement: Rate Limiting
The system SHALL support per-user and per-IP rate limiting.

#### Scenario: User rate limit
- **WHEN** configuration specifies 10 requests per minute per user
- **WHEN** user exceeds rate within 1 minute window
- **THEN** system returns 429 Too Many Requests error

#### Scenario: IP rate limit
- **WHEN** configuration specifies 100 requests per minute per IP
- **WHEN** IP exceeds rate within 1 minute window
- **THEN** system returns 429 Too Many Requests error

#### Scenario: Global rate limit
- **WHEN** configuration specifies 10000 requests per minute globally
- **WHEN** total requests exceed rate
- **THEN** system returns 429 Too Many Requests error

#### Scenario: Rate limit headers
- **WHEN** system processes request
- **THEN** system includes rate limit headers: X-RateLimit-Limit, X-RateLimit-Remaining, X-RateLimit-Reset

### Requirement: Usage Tracking
The system SHALL track per-user usage statistics.

#### Scenario: Request counting
- **WHEN** user makes request
- **THEN** system increments user's request counter

#### Scenario: Token counting
- **WHEN** request completes with token usage
- **THEN** system adds prompt_tokens, completion_tokens, total_tokens to user's counters

#### Scenario: Cost estimation
- **WHEN** request completes
- **THEN** system calculates estimated cost based on model pricing and token usage

#### Scenario: Usage persistence
- **WHEN** usage statistics updated
- **THEN** system persists to SQLite database

### Requirement: User Management
The system SHALL support user account management for local authentication.

#### Scenario: User creation
- **WHEN** admin creates user with username and password
- **THEN** system creates user record with bcrypt-hashed password

#### Scenario: User deletion
- **WHEN** admin deletes user
- **THEN** system removes user record and optionally associated tokens

#### Scenario: User enable/disable
- **WHEN** admin disables user
- **THEN** system rejects authentication attempts for disabled user

#### Scenario: Password reset
- **WHEN** admin resets user password
- **THEN** system updates password hash, invalidates existing sessions

### Requirement: Quota Configuration
The system SHALL support flexible quota configuration at global and user levels.

#### Scenario: Global default quotas
- **WHEN** configuration specifies global quota limits
- **THEN** system applies to all users without specific quotas

#### Scenario: User-specific quotas
- **WHEN** configuration specifies user-specific quotas
- **THEN** user quotas override global defaults

#### Scenario: Quota hot reload
- **WHEN** quota configuration changes
- **THEN** system applies new quotas to subsequent requests

#### Scenario: Unlimited quota
- **WHEN** quota set to 0 or -1
- **THEN** system treats as unlimited (no quota enforcement)
