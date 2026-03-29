## ADDED Requirements

### Requirement: Circuit Breaker States
The system SHALL implement circuit breaker with three states.

#### Scenario: Closed state
- **WHEN** requests succeed
- **THEN** circuit is closed, requests pass through

#### Scenario: Open state
- **WHEN** failure threshold exceeded
- **THEN** circuit opens, requests fail fast

#### Scenario: Half-open state
- **WHEN** timeout expires
- **THEN** circuit allows one test request

### Requirement: Failure Threshold
The system SHALL open circuit after failure threshold.

#### Scenario: Threshold exceeded
- **WHEN** consecutive failures exceed threshold
- **THEN** circuit opens

### Requirement: Recovery Timeout
The system SHALL attempt recovery after timeout.

#### Scenario: Recovery attempt
- **WHEN** timeout expires
- **THEN** circuit enters half-open state

#### Scenario: Recovery success
- **WHEN** test request succeeds
- **THEN** circuit closes

#### Scenario: Recovery failure
- **WHEN** test request fails
- **THEN** circuit opens again
