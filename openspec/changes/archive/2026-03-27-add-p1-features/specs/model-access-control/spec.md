## ADDED Requirements

### Requirement: Per-User Model Access
The system SHALL enforce per-user model access restrictions.

#### Scenario: Allowed model
- **WHEN** user requests model in their allowed list
- **THEN** system processes request

#### Scenario: Denied model
- **WHEN** user requests model not in their allowed list
- **THEN** system returns 403 Forbidden

#### Scenario: No restrictions
- **WHEN** user has no model restrictions configured
- **THEN** system allows all models

### Requirement: Global Default Models
The system SHALL support global default allowed models.

#### Scenario: Global default
- **WHEN** config specifies `routing.default_allowed_models`
- **WHEN** user has no specific restrictions
- **THEN** system uses global default list

#### Scenario: User override
- **WHEN** user has specific allowed models
- **THEN** user list overrides global default

### Requirement: Model Access Configuration
The system SHALL support configurable model access.

#### Scenario: User-specific models
- **WHEN** config specifies `routing.user_models.{user_id}`
- **THEN** system restricts user to specified models

#### Scenario: Model patterns
- **WHEN** config specifies model patterns (e.g., "gpt-*")
- **THEN** system allows all models matching pattern
