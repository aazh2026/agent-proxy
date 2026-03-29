## ADDED Requirements

### Requirement: OAuth2 Login Flow
The system SHALL support OAuth2 login with GitHub.

#### Scenario: Login initiation
- **WHEN** user visits `/auth/oauth2/login`
- **THEN** system redirects to GitHub authorization endpoint

#### Scenario: Callback processing
- **WHEN** GitHub redirects back to `/auth/oauth2/callback` with code
- **WHEN** system exchanges code for access token
- **THEN** system fetches user info, creates session

#### Scenario: Invalid callback
- **WHEN** GitHub returns error or invalid code
- **THEN** system returns 401 error

### Requirement: Organization Restriction
The system SHALL restrict login to allowed organizations.

#### Scenario: Allowed organization
- **WHEN** user is member of allowed organization
- **THEN** system allows login

#### Scenario: Denied organization
- **WHEN** user not in allowed organization
- **THEN** system returns 403 Forbidden

### Requirement: OAuth2 Configuration
The system SHALL support configurable OAuth2 settings.

#### Scenario: Client ID configuration
- **WHEN** config specifies `auth.oauth2.client_id`
- **THEN** system uses configured client ID

#### Scenario: Client secret configuration
- **WHEN** config specifies `auth.oauth2.client_secret`
- **THEN** system uses configured client secret

#### Scenario: Scopes configuration
- **WHEN** config specifies `auth.oauth2.scopes`
- **THEN** system requests configured scopes
