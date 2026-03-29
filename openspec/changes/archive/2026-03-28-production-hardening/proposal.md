## Why

The agent-proxy has all P0 and P1 features implemented. Before production deployment, the codebase needs hardening: security review, error handling improvements, graceful shutdown enhancements, and operational observability improvements.

## What Changes

- **Security hardening** — Review and fix potential security issues (token leakage, injection attacks, CORS)
- **Error handling** — Improve error messages, add error codes, better client-facing errors
- **Graceful shutdown** — Ensure in-flight requests complete before shutdown
- **Health checks** — Add deep health checks (database, provider connectivity)
- **Structured logging** — Add request correlation, structured JSON logs
- **Circuit breaker** — Add circuit breaker for failing providers

## Capabilities

### New Capabilities

- `security-hardening`: Security review and fixes
- `error-handling`: Improved error handling and messages
- `health-checks`: Deep health check endpoints
- `circuit-breaker`: Circuit breaker for provider failures

### Modified Capabilities

- None

## Impact

- **Modified files**: Most internal packages for error handling
- **New packages**: `internal/circuitbreaker/`
- **No breaking changes**: API surface unchanged
