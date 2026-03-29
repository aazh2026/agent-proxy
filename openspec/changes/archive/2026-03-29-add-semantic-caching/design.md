## Context

LLM API calls are expensive ($0.01-0.06 per 1K tokens) and slow (1-10 seconds). Many requests are semantically similar but produce identical responses. Caching can reduce costs and improve latency.

## Goals / Non-Goals

**Goals:**
- Cache LLM responses based on semantic similarity
- Reduce API costs by 30-50%
- Improve response latency for cached queries
- Transparent to clients (no code changes needed)

**Non-Goals:**
- Exact string matching only (semantic similarity)
- Distributed caching (single-instance only)
- Cache persistence across restarts (in-memory only in MVP)

## Decisions

### Decision 1: Hash-Based Similarity
**Rationale:**
- Simple and fast (no ML models needed)
- Normalize messages (lowercase, trim whitespace) then hash
- Configurable similarity threshold

### Decision 2: In-Memory Cache
**Rationale:**
- Fastest access (no network overhead)
- Simple implementation
- LRU eviction for memory management

### Decision 3: Cache Per Model
**Rationale:**
- Different models have different responses
- Separate cache namespaces per model
- Prevent cross-model cache pollution

## Risks / Trade-offs

### Risk 1: Stale Cache
**Risk:** Cached responses may become outdated
**Mitigation:**
- Configurable TTL (default 1 hour)
- Cache invalidation API
- Bypass header for fresh responses

### Risk 2: Memory Usage
**Risk:** Large responses consume memory
**Mitigation:**
- LRU eviction policy
- Configurable max cache size
- Response size limit for caching
