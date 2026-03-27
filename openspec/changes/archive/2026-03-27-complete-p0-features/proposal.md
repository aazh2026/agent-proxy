## Why

The agent-proxy has core LLM provider integration working but is missing critical P0 features from the PRD: rate limiting, quota enforcement, usage tracking, request logging integration, and hot reload. These are blocking production deployment and PRD compliance (currently at 50% completion).

## What Changes

- **Rate limiting integration** — Wire existing `routing/ratelimit.go` into request middleware
- **Quota enforcement integration** — Wire existing `routing/quota.go` into request pipeline
- **Usage statistics tracking** — Implement data collection for existing `usage_stats` database schema
- **Request logging integration** — Wire existing `observability/logger.go` into request handlers
- **Hot reload integration** — Wire existing `config/watcher.go` to reload configuration at runtime
- **Performance benchmarks** — Add latency and throughput measurement

## Capabilities

### New Capabilities

- `rate-limiting`: Per-user, per-IP, and global rate limiting enforced in request middleware
- `quota-enforcement`: Request count, token consumption, and cost quota enforcement per user
- `usage-tracking`: Real-time usage statistics collection and persistence
- `request-logging`: Request/response logging with ring buffer storage
- `hot-reload`: Configuration hot reload without service restart

### Modified Capabilities

- None (all code exists, just needs integration)

## Impact

- **Modified files**: `cmd/agent-proxy/main.go`, `internal/api/chat.go`, `internal/api/embeddings.go`
- **New middleware**: Rate limiting, quota enforcement, request logging
- **No new dependencies**: All code already exists in the codebase
- **No breaking changes**: API surface unchanged
