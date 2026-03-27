## Context

The agent-proxy has core functionality working (LLM provider integration, authentication, token management) but is missing critical P0 features for production deployment. The existing code has implementations for rate limiting, quota enforcement, logging, and config hot reload, but they are not wired into the request flow.

**Current State:**
- `internal/routing/ratelimit.go` — Rate limiter exists but not called
- `internal/routing/quota.go` — Quota tracker exists but not called
- `internal/observability/logger.go` — Request logger exists but not integrated
- `internal/config/watcher.go` — Config watcher exists but not started

**Constraints:**
- All code already exists — this is integration work, not new implementation
- Must not break existing functionality
- Must maintain backward compatibility

## Goals / Non-Goals

**Goals:**
- Wire rate limiting into request middleware
- Wire quota enforcement into request pipeline
- Integrate request logging into handlers
- Start config hot reload on startup
- Add performance measurement

**Non-Goals:**
- Rewrite existing implementations
- Add new features beyond PRD P0 requirements
- Performance optimization beyond measurement

## Decisions

### Decision 1: Middleware Integration Pattern
**Rationale:**
- Rate limiting should be a middleware (before auth)
- Quota enforcement should be in the pipeline (after auth)
- Request logging should be in handlers (after response)

### Decision 2: Minimal Changes
**Rationale:**
- All code exists — just wire it together
- No new packages needed
- No API changes required

## Risks / Trade-offs

### Risk 1: Performance Impact
**Risk:** Additional middleware may increase latency
**Mitigation:**
- Rate limiter uses in-memory token bucket (fast)
- Quota check is simple counter comparison
- Logging is async (ring buffer)

### Risk 2: Configuration Complexity
**Risk:** Hot reload may cause unexpected behavior
**Mitigation:**
- Validate config before applying
- Graceful fallback on validation failure
