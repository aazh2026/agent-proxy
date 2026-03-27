## Context

The agent-proxy currently has authentication, token management, and configuration infrastructure in place, but cannot actually proxy requests to LLM providers. All API endpoints (`/v1/chat/completions`, `/v1/embeddings`) return placeholder responses.

**Current State:**
- Chat handler returns hardcoded "This is a placeholder response"
- Provider interfaces exist (`internal/provider/`) but have no HTTP client integration
- Token resolver exists but is not connected to request flow
- Model routing registry exists but is not called during request processing

**Constraints:**
- Must maintain OpenAI-compatible API surface (no breaking changes)
- Must support streaming (SSE) with zero buffering
- Must inject provider-specific authentication headers
- Must handle provider errors gracefully and translate to OpenAI format

**Stakeholders:** LLM application developers using agent-proxy as a local gateway

## Goals / Non-Goals

**Goals:**
- Implement HTTP client layer for outbound requests to LLM providers
- Wire provider adapters to actual API endpoints (OpenAI, Anthropic, Google)
- Create request forwarding pipeline: Client → Auth → Router → Token → Provider → Response
- Implement bidirectional streaming proxy for real-time LLM output
- Add provider configuration for base URLs and timeouts

**Non-Goals:**
- Semantic caching or response caching (future enhancement)
- Cost-aware routing or A/B testing (future enhancement)
- Plugin architecture for third-party providers (future enhancement)
- Performance optimization beyond basic connection pooling

## Decisions

### Decision 1: Standard Library HTTP Client
**Rationale:**
- Zero external dependencies aligns with PRD requirement
- `net/http` client is production-ready and well-tested
- Connection pooling built-in with `http.Transport`
- Sufficient for expected load (10k QPS target)

**Alternatives considered:**
- `fasthttp`: Higher performance but less compatible, breaks `context.Context` support
- `github.com/hashicorp/go-cleanhttp`: Adds dependency for marginal benefit

**Implementation:**
```go
type Client struct {
    httpClient *http.Client
    baseURL    string
    apiKey     string
}
```

### Decision 2: Request Pipeline Architecture
**Rationale:**
- Clear separation of concerns (auth → routing → provider → response)
- Each stage is testable independently
- Easy to add middleware (logging, metrics, rate limiting)

**Pipeline stages:**
1. Authentication middleware (already exists)
2. Request validation (already exists)
3. Provider resolution via model routing
4. Token selection and decryption
5. HTTP request construction and forwarding
6. Response transformation to OpenAI format
7. Streaming proxy (if requested)

### Decision 3: Streaming via io.Pipe
**Rationale:**
- Zero-copy streaming from provider to client
- Context cancellation propagation
- No buffering of entire response in memory

**Implementation:**
```go
func proxyStream(w http.ResponseWriter, resp *http.Response) {
    defer resp.Body.Close()
    flusher := w.(http.Flusher)
    buf := make([]byte, 4096)
    for {
        n, err := resp.Body.Read(buf)
        if n > 0 {
            w.Write(buf[:n])
            flusher.Flush()
        }
        if err != nil {
            break
        }
    }
}
```

### Decision 4: Provider Configuration in YAML
**Rationale:**
- Consistent with existing configuration approach
- Supports hot reload via file watcher
- Environment variable override for containerized deployments

**Configuration structure:**
```yaml
providers:
  openai:
    base_url: "https://api.openai.com/v1"
    timeout: 30s
  anthropic:
    base_url: "https://api.anthropic.com/v1"
    timeout: 30s
  google:
    base_url: "https://generativelanguage.googleapis.com/v1beta"
    timeout: 30s
```

## Risks / Trade-offs

### Risk 1: Provider API Changes
**Risk:** LLM providers may change API formats, breaking integration
**Mitigation:**
- Version detection per provider
- Comprehensive integration tests
- Clear error messages when provider returns unexpected format

### Risk 2: Streaming Timeout
**Risk:** Long-running streams may timeout or consume excessive resources
**Mitigation:**
- Configurable per-provider timeout
- Context cancellation propagation
- Connection limits via `http.Transport`

### Risk 3: Memory Usage During Streaming
**Risk:** Large responses could cause memory pressure
**Mitigation:**
- Zero-copy streaming via `io.Copy`
- No response buffering
- Configurable buffer size

## Migration Plan

### Phase 1: HTTP Client Layer
1. Create `internal/httpclient/` package
2. Implement provider-specific clients (OpenAI, Anthropic, Google)
3. Add connection pooling and timeout configuration

### Phase 2: Request Pipeline
1. Create `internal/pipeline/` package
2. Wire authentication → routing → provider → response
3. Integrate token resolver for authentication injection

### Phase 3: Streaming Proxy
1. Implement bidirectional streaming in pipeline
2. Add context cancellation handling
3. Test with real provider endpoints

### Rollback Strategy
- Feature flag: `providers.enabled: true/false`
- Graceful degradation to placeholder responses if providers unavailable
- Existing API surface unchanged (no breaking changes)

## Open Questions

1. **Retry strategy for provider errors**: Exponential backoff with jitter?
   - Recommendation: Yes, with configurable max retries (default 3)

2. **Connection pool size**: Should we limit concurrent connections per provider?
   - Recommendation: Yes, default 100 per provider, configurable

3. **Request timeout**: Should we have separate timeouts for connect vs read?
   - Recommendation: Yes, 10s connect, 60s read (configurable)
