## ADDED Requirements

### Requirement: OIDC Login Flow
The system SHALL support OIDC login with Google.

#### Scenario: Login initiation
- **WHEN** user visits `/auth/oidc/login`
- **THEN** system redirects to Google authorization endpoint

#### Scenario: Callback processing
- **WHEN** Google redirects back to `/auth/oidc/callback` with code
- **WHEN** system exchanges code for tokens
- **THEN** system creates session and returns session token

#### Scenario: Invalid callback
- **WHEN** Google returns error or invalid code
- **THEN** system returns 401 error

### Requirement: Domain Restriction
The system SHALL restrict login to allowed email domains.

#### Scenario: Allowed domain
- **WHEN** user email domain is in allowed list
- **THEN** system allows login

#### Scenario: Denied domain
- **WHEN** user email domain not in allowed list
- **THEN** system returns 403 Forbidden

### Requirement: OIDC Configuration
The system SHALL support configurable OIDC settings.

#### Scenario: Client ID configuration
- **WHEN** config specifies `auth.oidc.client_id`
- **THEN** system uses configured client ID

#### Scenario: Client secret configuration
- **WHEN** config specifies `auth.oidc.client_secret`
- **THEN** system uses configured client secret
