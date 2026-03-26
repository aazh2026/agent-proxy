## ADDED Requirements

### Requirement: X-User-ID Header Authentication
The system SHALL support user identification via X-User-ID request header for zero-configuration authentication.

#### Scenario: Valid X-User-ID header
- **WHEN** client sends request with `X-User-ID: alice` header
- **THEN** system identifies request as user "alice" and proceeds with request processing

#### Scenario: Missing X-User-ID header
- **WHEN** client sends request without X-User-ID header
- **THEN** system identifies request as user "default" and proceeds with request processing

#### Scenario: X-User-ID format validation
- **WHEN** client sends request with X-User-ID containing invalid characters (not alphanumeric, underscore, or hyphen)
- **WHEN** validation enabled
- **THEN** system returns 400 Bad Request indicating invalid user ID format

#### Scenario: X-User-ID length validation
- **WHEN** client sends request with X-User-ID longer than 64 characters
- **THEN** system returns 400 Bad Request indicating user ID exceeds maximum length

#### Scenario: Allowed user IDs whitelist
- **WHEN** system configured with `allowed_user_ids` whitelist
- **WHEN** client sends X-User-ID not in whitelist
- **THEN** system returns 401 Unauthorized indicating user not in allowed list

### Requirement: Local User Authentication
The system SHALL support local user database authentication with username/password and API key modes.

#### Scenario: Username/password login
- **WHEN** client sends POST to `/auth/login` with valid username and password
- **THEN** system returns session token in response with expiration time

#### Scenario: Invalid credentials
- **WHEN** client sends POST to `/auth/login` with incorrect password
- **THEN** system returns 401 Unauthorized indicating invalid credentials

#### Scenario: API key authentication
- **WHEN** client sends request with `Authorization: Bearer {user_api_key}` header
- **THEN** system validates API key and identifies associated user

#### Scenario: Disabled user login attempt
- **WHEN** client attempts login with disabled user account
- **THEN** system returns 401 Unauthorized indicating account is disabled

#### Scenario: User password storage
- **WHEN** system stores user password
- **THEN** system stores bcrypt hash, never plaintext password

### Requirement: OIDC Authentication
The system SHALL support OpenID Connect authentication for third-party identity providers (Google, Azure AD, Okta).

#### Scenario: OIDC login flow initiation
- **WHEN** client sends GET to `/auth/oidc/login`
- **THEN** system redirects to configured OIDC provider's authorization endpoint with required scopes

#### Scenario: OIDC callback processing
- **WHEN** OIDC provider redirects back to `/auth/oidc/callback` with valid authorization code
- **WHEN** system exchanges code for tokens and validates ID token
- **THEN** system creates session and returns session token

#### Scenario: Domain restriction enforcement
- **WHEN** system configured with `allowed_domains: ["example.com"]`
- **WHEN** OIDC user has email "user@other.com"
- **THEN** system returns 403 Forbidden indicating domain not allowed

#### Scenario: Invalid OIDC token
- **WHEN** OIDC provider returns invalid or expired authorization code
- **THEN** system returns 401 Unauthorized indicating authentication failed

### Requirement: OAuth2 Authentication
The system SHALL support OAuth2 authentication for third-party platforms (GitHub, GitLab).

#### Scenario: OAuth2 login flow
- **WHEN** client sends GET to `/auth/oauth2/login`
- **THEN** system redirects to configured OAuth2 provider's authorization endpoint with configured scopes

#### Scenario: OAuth2 callback processing
- **WHEN** OAuth2 provider redirects back with valid authorization code
- **WHEN** system exchanges code for access token and fetches user info
- **THEN** system creates session and returns session token

#### Scenario: Organization restriction
- **WHEN** system configured with `allowed_organizations: ["myorg"]`
- **WHEN** user not in organization
- **THEN** system returns 403 Forbidden indicating organization not allowed

### Requirement: Session Token Management
The system SHALL support session-based authentication with configurable expiration.

#### Scenario: Session creation
- **WHEN** user successfully authenticates via any method
- **THEN** system creates session with unique session_id, user_id, expires_at, created_at fields

#### Scenario: Session validation
- **WHEN** client sends request with `Authorization: Bearer {session_token}`
- **THEN** system validates session is not expired and identifies associated user

#### Scenario: Expired session rejection
- **WHEN** client sends request with expired session token
- **THEN** system returns 401 Unauthorized indicating session expired

#### Scenario: Session cleanup
- **WHEN** session expires
- **THEN** system automatically removes expired session from storage

#### Scenario: Session storage options
- **WHEN** system stores sessions
- **THEN** system supports both in-memory and SQLite storage backends

### Requirement: Authentication Configuration
The system SHALL support configurable authentication method selection with single active method at runtime.

#### Scenario: Method selection via config
- **WHEN** configuration specifies `auth.method: "x-user-id"`
- **WHEN** system processes request
- **THEN** system uses X-User-ID header authentication only

#### Scenario: Invalid auth method
- **WHEN** configuration specifies unsupported auth method
- **THEN** system fails to start with error indicating invalid auth method

#### Scenario: Auth method hot reload
- **WHEN** configuration file changes auth method
- **THEN** system reloads configuration and uses new auth method for subsequent requests
