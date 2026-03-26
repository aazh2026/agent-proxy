## Why

LLM application developers face three critical pain points: managing API keys across multiple providers (OpenAI, Anthropic, Google), implementing complex authentication flows, and lacking observability into token usage and costs. This creates security risks from scattered API keys, development friction from multi-provider code changes, and operational blind spots with no centralized visibility.

**Agent Proxy** solves this by providing a local, zero-dependency proxy that presents a unified OpenAI-compatible endpoint while handling authentication, token security, and multi-provider routing transparently. Applications only need to change their `base_url`—no code modifications required.

## What Changes

- **New local proxy service** (`agent-proxy`) — single binary, zero external dependencies, runs on configurable port
- **OpenAI-compatible API** — full support for `/v1/chat/completions` and `/v1/embeddings` endpoints with streaming (SSE)
- **Multi-provider authentication** — configurable auth providers: X-User-ID header, local user database, Google OIDC, OAuth2, Session Token
- **Encrypted token management** — AES-256-GCM encrypted storage, tokens never leave proxy boundary, never logged in plaintext
- **Provider adapter layer** — protocol translation for OpenAI, Anthropic Claude, Google Gemini with plugin architecture
- **Intelligent routing** — model-based routing, multi-token load balancing (round-robin, weighted, priority), automatic failover
- **Multi-user isolation** — per-user token isolation, model access control, usage quotas, rate limiting
- **Built-in observability** — real-time metrics dashboard, request logs, token management UI, usage statistics
- **Configuration system** — YAML config with environment variable override, hot reload support

## Capabilities

### New Capabilities

- `openai-api`: OpenAI-compatible REST API endpoints with streaming support, parameter validation, and error standardization
- `authentication`: Multi-provider identity authentication system (X-User-ID, local users, OIDC, OAuth2, session tokens)
- `token-management`: Encrypted token storage, lifecycle management, auto-refresh, and residency security enforcement
- `provider-adapters`: LLM provider protocol translation (OpenAI, Anthropic, Gemini) with token injection and response normalization
- `model-routing`: Model name-based routing, multi-token load balancing, failover, and cross-provider fallback
- `user-isolation`: Per-user token isolation, model access control, usage quotas, and rate limiting
- `observability`: Real-time metrics, request logging, token management UI, and usage analytics dashboard
- `configuration`: YAML-based configuration with environment variable override, hot reload, and CLI parameters

### Modified Capabilities

- None (greenfield project)

## Impact

- **New codebase**: Go implementation with zero external dependencies (single binary delivery)
- **New API surface**: OpenAI-compatible endpoints at `http://localhost:{port}`
- **New data storage**: SQLite for encrypted tokens and usage statistics
- **New dependencies**: SQLite driver, encryption libraries (AES-256-GCM), HTTP server framework
- **Target platforms**: Linux, macOS, Windows (x86_64, arm64)
- **Performance targets**: <5ms proxy latency, 10k+ QPS, <50MB memory, <100ms startup
