# Agent Proxy

A local, zero-dependency LLM authentication proxy that provides a unified OpenAI-compatible endpoint while handling authentication, token security, and multi-provider routing transparently.

## Features

- **OpenAI-Compatible API**: Full support for `/v1/chat/completions` and `/v1/embeddings` endpoints
- **Multi-Provider Support**: OpenAI, Anthropic Claude, Google Gemini
- **Secure Token Management**: AES-256-GCM encrypted storage, tokens never leave proxy boundary
- **Multiple Authentication Methods**: X-User-ID header, local users, OIDC, OAuth2
- **Intelligent Routing**: Model-based routing, load balancing, failover
- **Built-in Observability**: Real-time metrics, request logging, Web UI dashboard
- **Single Binary**: Zero dependencies, cross-platform deployment

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

### Verify

```bash
# Check health
curl http://localhost:4000/health

# View metrics
curl http://localhost:4000/metrics
```

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

logging:
  level: "info"
```

See `agent-proxy.example.yaml` for all configuration options.

## API Usage

### Chat Completions

```bash
curl -X POST http://localhost:4000/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "X-User-ID: alice" \
  -d '{
    "model": "gpt-4",
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

## Dashboard

Access the admin dashboard at `http://localhost:4000/`

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
