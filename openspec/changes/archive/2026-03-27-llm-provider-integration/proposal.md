## Why

The current agent-proxy implementation has foundational infrastructure (authentication, token management, configuration) but **cannot actually proxy requests to LLM providers**. All API endpoints return placeholder responses instead of forwarding requests to OpenAI, Anthropic, or Google. This blocks the core value proposition of the product.

**Gap Analysis Results:**
- P0 features: 42% complete (5/12 missing critical components)
- No HTTP client layer exists for outbound provider requests
- Provider adapters are interfaces only, not connected to real APIs
- Streaming returns hardcoded chunks, not real LLM output

## What Changes

- **New HTTP client layer** — connection pooling, timeouts, retry logic for provider communication
- **Real provider integration** — OpenAI, Anthropic Claude, Google Gemini API endpoints
- **Request forwarding pipeline** — Client → Router → Token Resolver → Provider → Response Transform → Client
- **Bidirectional streaming proxy** — Real-time streaming from upstream providers to clients
- **Provider configuration** — Base URLs, API versions, timeout settings per provider
- **Token injection middleware** — Automatic authentication header injection per provider

## Capabilities

### New Capabilities

- `http-client`: HTTP client layer with connection pooling, configurable timeouts, and retry logic
- `provider-integration`: Real LLM provider API integration (OpenAI, Anthropic, Google)
- `request-pipeline`: Request forwarding pipeline connecting auth → routing → provider → response
- `streaming-proxy`: Bidirectional SSE streaming proxy with context cancellation

### Modified Capabilities

- `openai-api`: Handlers now forward to real providers instead of returning placeholders
- `provider-adapters`: Adapters now make actual HTTP calls to provider APIs
- `model-routing`: Registry now integrated into request flow for provider selection

## Impact

- **Modified files**: `cmd/agent-proxy/main.go`, `internal/api/chat.go`, `internal/api/embeddings.go`
- **New packages**: `internal/httpclient/`, `internal/pipeline/`
- **New dependencies**: None (uses standard library `net/http`)
- **Configuration changes**: New provider endpoint configuration sections
- **Breaking changes**: None (API surface unchanged, only behavior changes)
