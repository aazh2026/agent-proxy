# Agent Proxy

A local, zero-dependency LLM authentication proxy that provides a unified OpenAI-compatible endpoint while handling authentication, token security, and multi-provider routing transparently.

## Features

- **OpenAI-Compatible API**: Full support for `/v1/chat/completions` and `/v1/embeddings` endpoints
- **Multi-Provider Support**: OpenAI, Anthropic Claude, Google Gemini
- **Secure Token Management**: AES-256-GCM encrypted storage, tokens never leave proxy boundary
- **Multiple Authentication Methods**: X-User-ID header, local users, OIDC, OAuth2
- **Provider OAuth Support**: OpenAI browser OAuth flow + auto token persistence
- **Intelligent Routing**: Model-based routing, load balancing, failover
- **Built-in Observability**: Real-time metrics, request logging, Web UI dashboard
- **Single Binary**: Zero dependencies, cross-platform deployment

## User Guide

This section is a concise step-by-step guide to using Agent Proxy for the first time.

1. Install and start:
   - Download prebuilt binary from releases or build from source (`make build`).
   - Run `./agent-proxy` (default port 4000) or `./agent-proxy --config /path/to/agent-proxy.yaml`.

2. Verify service status:
   - `curl http://localhost:4000/health`
   - `curl http://localhost:4000/metrics`

3. Register provider tokens (using `X-User-ID` or your auth method):
   - OpenAI:
     ```bash
     curl -X POST http://localhost:4000/tokens \
       -H "Content-Type: application/json" \
       -H "X-User-ID: alice" \
       -d '{"provider":"openai","type":"api_key","access_token":"sk-..."}'
     ```
   - Anthropic:
     ```bash
     curl -X POST http://localhost:4000/tokens \
       -H "Content-Type: application/json" \
       -H "X-User-ID: alice" \
       -d '{"provider":"anthropic","type":"api_key","access_token":"sk-..."}'
     ```

4. Send requests through the proxy:
   - `gpt-*` models route to OpenAI, `claude-*` to Anthropic, `gemini-*` to Google.
   - Example chat completion:
     ```bash
     curl -X POST http://localhost:4000/v1/chat/completions \
       -H "Content-Type: application/json" \
       -H "X-User-ID: alice" \
       -d '{"model":"gpt-4","messages":[{"role":"user","content":"Hello"}]}'
     ```

5. Manage tokens:
   - List: `curl -H "X-User-ID: alice" http://localhost:4000/tokens`
   - Delete: `curl -X DELETE "http://localhost:4000/tokens/delete?token_id=<id>" -H "X-User-ID: alice"`

6. Explore the dashboard:
   - Visit `http://localhost:4000/` for Web UI and real-time metrics.

7. Customize config:
   - Copy `agent-proxy.example.yaml`, edit values (ports, providers, auth, logging), and restart.

> Tip: Use `agent-proxy --help` for all CLI flags and `agent-proxy.example.yaml` for full config options.

## Quick Start

### Download

Download the latest release for your platform from the releases page.

### Run

```bash
# Start with default configuration
./agent-proxy

# Start with custom configuration
./agent-proxy --config /path/to/agent-proxy.yaml

# Start on custom port
./agent-proxy --port 8080
```

### Build + Start Script

A helper script is available for local development:

```bash
./scripts/build-and-start.sh --config /path/to/agent-proxy.yaml --port 4000
```

This builds `agent-proxy` and immediately starts the server using the same process.


### Verify

```bash
# Check health
curl http://localhost:4000/health

# View metrics
curl http://localhost:4000/metrics
```

## Provider Setup

### 1. Add Your API Tokens

Before making requests, add your LLM provider API tokens:

```bash
# Add OpenAI token
curl -X POST http://localhost:4000/tokens \
  -H "Content-Type: application/json" \
  -H "X-User-ID: alice" \
  -d '{
    "provider": "openai",
    "type": "api_key",
    "access_token": "sk-your-openai-key"
  }'

# Add Anthropic token
curl -X POST http://localhost:4000/tokens \
  -H "Content-Type: application/json" \
  -H "X-User-ID: alice" \
  -d '{
    "provider": "anthropic",
    "type": "api_key",
    "access_token": "sk-ant-your-anthropic-key"
  }'

# Add Google Gemini token
curl -X POST http://localhost:4000/tokens \
  -H "Content-Type: application/json" \
  -H "X-User-ID: alice" \
  -d '{
    "provider": "google",
    "type": "api_key",
    "access_token": "your-google-api-key"
  }'
```

### 2. Make Requests

Tokens are automatically routed based on model name:

- `gpt-*` models → OpenAI
- `claude-*` models → Anthropic
- `gemini-*` models → Google

## Configuration

Create a `agent-proxy.yaml` file:

```yaml
server:
  host: "127.0.0.1"
  port: 4000

auth:
  method: "x-user-id"

token:
  encryption_key: ""  # Leave empty to auto-generate
  storage_path: "agent-proxy.db"

providers:
  openai:
    enabled: true
    base_url: "https://api.openai.com/v1"
    timeout_seconds: 60
  anthropic:
    enabled: true
    base_url: "https://api.anthropic.com/v1"
    timeout_seconds: 60
  google:
    enabled: true
    base_url: "https://generativelanguage.googleapis.com/v1beta"
    timeout_seconds: 60

logging:
  level: "info"
```

See `agent-proxy.example.yaml` for all configuration options.

## API Usage

### Chat Completions

```bash
# OpenAI GPT-4
curl -X POST http://localhost:4000/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "X-User-ID: alice" \
  -d '{
    "model": "gpt-4",
    "messages": [{"role": "user", "content": "Hello!"}]
  }'

# Anthropic Claude
curl -X POST http://localhost:4000/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "X-User-ID: alice" \
  -d '{
    "model": "claude-3-opus-20240229",
    "messages": [{"role": "user", "content": "Hello!"}]
  }'

# Google Gemini
curl -X POST http://localhost:4000/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "X-User-ID: alice" \
  -d '{
    "model": "gemini-pro",
    "messages": [{"role": "user", "content": "Hello!"}]
  }'
```

### Streaming

```bash
curl -X POST http://localhost:4000/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "X-User-ID: alice" \
  -d '{
    "model": "gpt-4",
    "messages": [{"role": "user", "content": "Hello!"}],
    "stream": true
  }'
```

### Embeddings

```bash
curl -X POST http://localhost:4000/v1/embeddings \
  -H "Content-Type: application/json" \
  -H "X-User-ID: alice" \
  -d '{
    "model": "text-embedding-ada-002",
    "input": "Hello world"
  }'
```

## Token Management

### Add Token

```bash
curl -X POST http://localhost:4000/tokens \
  -H "Content-Type: application/json" \
  -H "X-User-ID: alice" \
  -d '{
    "provider": "openai",
    "type": "api_key",
    "access_token": "sk-..."
  }'
```

### List Tokens

```bash
curl http://localhost:4000/tokens \
  -H "X-User-ID: alice"
```

### Delete Token

```bash
curl -X DELETE "http://localhost:4000/tokens/delete?token_id=tk_..." \
  -H "X-User-ID: alice"
```

## Authentication

### X-User-ID Header (Default)

```bash
curl -H "X-User-ID: alice" http://localhost:4000/v1/chat/completions
```

### Local Users

```bash
# Login
curl -X POST http://localhost:4000/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "alice", "password": "password"}'

# Use session token
curl -H "Authorization: Bearer <session_token>" http://localhost:4000/v1/chat/completions
```

### OpenAI OAuth (Browser flow)

1. Configure `openai.oauth_client_id`, `openai.oauth_client_secret`, and `openai.oauth_redirect_uri` in `agent-proxy.yaml`.
2. Start browser flow:

```bash
# Must include user identity (X-User-ID or authenticated session)
curl -i -H "X-User-ID: alice" http://localhost:4000/auth/openai/login
```

3. Complete OAuth consent in browser; OpenAI redirects back to `/auth/openai/callback`.
4. Proxy stores encrypted OpenAI access token and refresh token for user `alice` in the token store.
5. Subsequent calls to `/v1/chat/completions` with model `gpt-*` use the persisted OpenAI token automatically.


## Dashboard

Access the admin dashboard at `http://localhost:4000/`.

- In the dashboard, the new Configuration panel lets you view and edit current config JSON, save it to disk, and trigger a reload via the config watcher.


## Building from Source

```bash
# Clone repository
git clone https://github.com/openclaw/agent-proxy.git
cd agent-proxy

# Build
make build

# Run
./agent-proxy

# Cross-compile
make build-all
```

## License

MIT License
