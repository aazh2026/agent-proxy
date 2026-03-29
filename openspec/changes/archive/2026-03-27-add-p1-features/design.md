## Context

The agent-proxy has placeholder implementations for OIDC and OAuth2 authentication. These need to be implemented with proper OAuth2 flows using the `golang.org/x/oauth2` library. Token auto-refresh needs to be integrated with the existing token management system.

## Goals / Non-Goals

**Goals:**
- Implement OIDC authentication with Google
- Implement OAuth2 authentication with GitHub
- Implement token auto-refresh
- Implement model access control
- Add performance benchmarks

**Non-Goals:**
- Support all OIDC providers (start with Google)
- Support all OAuth2 providers (start with GitHub)
- Complex RBAC permissions

## Decisions

### Decision 1: Use golang.org/x/oauth2 Library
**Rationale:**
- Official Google OAuth2 library
- Well-tested and maintained
- Supports all required OAuth2 flows

### Decision 2: Token Refresh Strategy
**Rationale:**
- Proactive refresh 30 minutes before expiration
- Automatic retry on failure
- Disable token if refresh fails

## Risks / Trade-offs

### Risk 1: OAuth2 Library Dependency
**Risk:** External dependency adds complexity
**Mitigation:**
- Use official, well-maintained library
- Isolate OAuth2 code in auth package
