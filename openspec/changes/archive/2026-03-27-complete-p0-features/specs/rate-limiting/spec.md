## ADDED Requirements

### Requirement: Per-User Rate Limiting
The system SHALL enforce rate limits per user to prevent abuse.

#### Scenario: User within rate limit
- **WHEN** user makes request within rate limit (e.g., 60 requests/minute)
- **THEN** request is processed normally

#### Scenario: User exceeds rate limit
- **WHEN** user exceeds configured rate limit
- **THEN** system returns 429 Too Many Requests error

#### Scenario: Rate limit reset
- **WHEN** rate limit window expires
- **THEN** user's request count resets

### Requirement: Per-IP Rate Limiting
The system SHALL enforce rate limits per IP address.

#### Scenario: IP within rate limit
- **WHEN** IP makes requests within limit
- **THEN** requests are processed normally

#### Scenario: IP exceeds rate limit
- **WHEN** IP exceeds configured rate limit
- **THEN** system returns 429 Too Many Requests error

### Requirement: Global Rate Limiting
The system SHALL enforce global rate limits across all users.

#### Scenario: Global within limit
- **WHEN** total requests within global limit
- **THEN** requests are processed normally

#### Scenario: Global exceeds limit
- **WHEN** total requests exceed global limit
- **THEN** system returns 429 Too Many Requests error

### Requirement: Rate Limit Headers
The system SHALL include rate limit information in response headers.

#### Scenario: Rate limit headers present
- **WHEN** system processes request
- **THEN** response includes X-RateLimit-Limit, X-RateLimit-Remaining, X-RateLimit-Reset headers
