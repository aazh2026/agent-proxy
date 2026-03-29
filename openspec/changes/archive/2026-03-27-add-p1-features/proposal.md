## Why

The agent-proxy has core P0 features complete (96%) and most P1 features implemented (85%). Remaining gaps include: OIDC/OAuth2 authentication (currently placeholders), token auto-refresh, model access control enforcement, and performance validation. Completing these brings P1 completion to 100% and enables enterprise/production deployment.

## What Changes

- **OIDC authentication** — Implement Google OIDC login flow with domain restrictions
- **OAuth2 authentication** — Implement GitHub/GitLab OAuth2 login with organization restrictions
- **Token auto-refresh** — Automatic refresh of expiring OAuth tokens
- **Model access control** — Per-user model restrictions enforcement
- **Performance benchmarks** — Validate PRD targets (<5ms latency, >10k QPS)

## Capabilities

### New Capabilities

- `oidc-auth`: OpenID Connect authentication with Google, domain restrictions
- `oauth2-auth`: OAuth2 authentication with GitHub/GitLab, organization restrictions
- `token-auto-refresh`: Automatic token refresh before expiration
- `model-access-control`: Per-user model access restrictions

### Modified Capabilities

- None (implementation changes only, no requirement changes)

## Impact

- **Modified files**: `internal/auth/oidc.go`, `internal/auth/oauth2.go`, `internal/token/refresh.go`
- **New dependencies**: `golang.org/x/oauth2` for OAuth2 flows
- **Configuration changes**: New OIDC/OAuth2 provider settings
- **No breaking changes**: API surface unchanged
