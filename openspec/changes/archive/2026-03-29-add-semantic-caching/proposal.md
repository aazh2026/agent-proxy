## Why

LLM API calls are expensive and slow. Many requests are semantically similar (e.g., "What is Go?" vs "Explain Go programming language") but produce identical responses. Semantic caching can reduce API costs by 30-50% and improve response latency from seconds to milliseconds for cached queries.

## What Changes

- **Semantic cache** — Cache LLM responses based on message similarity
- **Cache key generation** — Generate cache keys from message content
- **Similarity matching** — Find cached responses for semantically similar queries
- **Cache management** — TTL, eviction, invalidation, statistics
- **Cache bypass** — Option to skip cache for specific requests

## Capabilities

### New Capabilities

- `semantic-cache`: Cache LLM responses with semantic similarity matching

### Modified Capabilities

- None

## Impact

- **New package**: `internal/cache/`
- **Modified files**: `internal/api/chat.go`, `internal/api/embeddings.go`
- **New dependencies**: None (hash-based similarity)
- **Configuration**: New cache settings in config.yaml
- **No breaking changes**: Cache is transparent to clients
