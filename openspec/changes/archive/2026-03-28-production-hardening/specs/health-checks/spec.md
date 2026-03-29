## ADDED Requirements

### Requirement: Basic Health Check
The system SHALL provide basic health check endpoint.

#### Scenario: Health check response
- **WHEN** GET /health
- **THEN** returns 200 with status "healthy"

### Requirement: Deep Health Check
The system SHALL provide deep health check endpoint.

#### Scenario: Database check
- **WHEN** GET /health/ready
- **THEN** checks database connectivity

#### Scenario: Provider check
- **WHEN** GET /health/ready
- **THEN** checks provider connectivity

### Requirement: Liveness Check
The system SHALL provide liveness check endpoint.

#### Scenario: Liveness response
- **WHEN** GET /health/live
- **THEN** returns 200 if process is running
