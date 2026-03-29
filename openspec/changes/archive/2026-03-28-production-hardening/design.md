## Context

The agent-proxy has all P0 and P1 features implemented. The codebase needs production hardening before deployment: security review, error handling improvements, and operational improvements.

## Goals / Non-Goals

**Goals:**
- Fix security issues (token leakage, injection)
- Improve error handling
- Add circuit breaker for provider failures
- Add deep health checks
- Improve logging

**Non-Goals:**
- New features
- Performance optimization (already benchmarked)
- UI improvements

## Decisions

### Decision 1: Circuit Breaker Pattern
**Rationale:**
- Prevent cascading failures
- Fast failure when provider is down
- Automatic recovery when provider comes back

### Decision 2: Structured JSON Logging
**Rationale:**
- Better for log aggregation systems
- Easier to parse and query
- Request correlation via request ID

## Risks / Trade-offs

### Risk 1: Breaking Changes
**Risk:** Error handling changes may affect clients
**Mitigation:**
- Keep error format compatible with OpenAI
- Add new fields, don't remove existing ones
