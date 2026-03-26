## Context

LLM application developers currently face significant integration challenges:
- **Multi-provider complexity**: Each LLM provider (OpenAI, Anthropic, Google) has different API formats, authentication methods, and error codes
- **Security risks**: API keys scattered across applications, config files, and environment variables create leakage vectors
- **Operational blindness**: No centralized visibility into token usage, costs, or request patterns across applications
- **Development friction**: Switching providers or adding new ones requires code changes throughout the application stack

**Current state**: Developers either hardcode provider-specific implementations or use provider-specific SDKs, creating tight coupling and security vulnerabilities.

**Constraints**:
- Must run locally with zero external dependencies (no databases, no Docker required)
- Must achieve <5ms proxy latency to be transparent to applications
- Must support streaming (SSE) with zero buffering
- Must be a single binary for cross-platform deployment

**Stakeholders**: LLM application developers, DevOps teams, security teams, enterprise IT departments

## Goals / Non-Goals

**Goals:**
- Provide a unified OpenAI-compatible API endpoint that requires zero application code changes
- Implement secure token residency—tokens never leave the proxy boundary
- Support multiple authentication providers with pluggable architecture
- Deliver production-grade observability without external dependencies
- Achieve <5ms proxy latency and 10k+ QPS on commodity hardware
- Enable multi-user isolation with per-user token and quota management

**Non-Goals:**
- Agent orchestration or workflow management
- Prompt engineering or template management
- LLM model training, fine-tuning, or inference hosting
- Cloud-hosted SaaS service (local-only deployment)
- Complex RBAC permissions (user-level isolation only in MVP)
- Container orchestration (single binary, not Docker-first)

## Decisions

### Decision 1: Go as Primary Implementation Language
**Rationale**:
- Excellent standard library for HTTP servers and concurrency
- Single binary compilation for all platforms (cross-compilation built-in)
- Mature SQLite driver ecosystem (go-sqlite3)
- Fast development iteration for MVP
- Large developer community for maintenance

**Alternatives considered**:
- Rust: Better performance characteristics but slower development velocity for MVP
- Node.js: Good ecosystem but higher memory footprint and startup time
- Python: Easy to prototype but poor performance for proxy workloads

**Trade-off**: Go's GC introduces ~1ms P99 latency variance vs Rust's zero-GC guarantee, but development speed is 2-3x faster for MVP phase.

### Decision 2: SQLite for Persistent Storage
**Rationale**:
- Zero external dependencies (embedded database)
- Single file storage with ACID guarantees
- Excellent performance for read-heavy workloads (token lookups)
- Built-in encryption support via SQLCipher or application-level AES-256-GCM
- Battle-tested reliability

**Alternatives considered**:
- File-based JSON/YAML: No transactional guarantees, race conditions on concurrent writes
- Embedded PostgreSQL: External dependency, more complex setup
- In-memory only: Data loss on restart, unacceptable for token management

### Decision 3: AES-256-GCM for Token Encryption
**Rationale**:
- Industry-standard authenticated encryption
- Provides both confidentiality and integrity
- Hardware acceleration available on modern CPUs
- NIST-approved, suitable for enterprise security requirements

**Implementation**: Tokens encrypted at rest in SQLite, decrypted only during request injection, immediately cleared from memory after use.

### Decision 4: Data Plane / Control Plane Separation
**Rationale**:
- Request forwarding (data plane) must be isolated from management operations (control plane)
- Control plane failures (metrics, logs, UI) must not impact request forwarding
- Different performance requirements: data plane needs <5ms latency, control plane can tolerate higher latency

**Architecture**:
- Data plane: HTTP server, auth middleware, token injection, response forwarding
- Control plane: Separate goroutine pool for metrics aggregation, log writing, Web UI serving
- Communication: Lock-free channels for metrics, async writes for logs

### Decision 5: Provider Adapter Plugin Architecture
**Rationale**:
- Each LLM provider has unique API formats, authentication methods, and streaming implementations
- New providers should be addable without modifying core proxy logic
- Testing providers in isolation is easier with plugin architecture

**Design**:
- `Provider` interface with methods: `TransformRequest()`, `TransformResponse()`, `InjectAuth()`, `HandleStreaming()`
- Registry pattern for provider discovery
- Configuration-driven provider selection

### Decision 6: Hot-Reload Configuration
**Rationale**:
- Production deployments need configuration changes without downtime
- Token rotation, user management, and routing rules change frequently
- File-watching is simpler and more reliable than API-based configuration

**Implementation**:
- YAML configuration file with file-system watcher
- Atomic configuration swap (new config validated before activation)
- Environment variable override for containerized deployments

## Risks / Trade-offs

### Risk 1: Token Leakage via Memory Dumps
**Risk**: Application crashes could produce core dumps containing decrypted tokens
**Mitigation**:
- Immediate memory clearing after token use (zero-overwrite pattern)
- Disable core dumps via `ulimit -c 0` in production
- Security audit of all error paths

### Risk 2: Streaming Performance Under Load
**Risk**: High-concurrency streaming could exhaust file descriptors or goroutine pools
**Mitigation**:
- Connection pooling with configurable limits
- Backpressure mechanisms (slow client detection)
- Graceful degradation (reject new connections at capacity)
- Load testing at 2x expected capacity before release

### Risk 3: Provider API Breaking Changes
**Risk**: LLM providers may change API formats, breaking adapter layer
**Mitigation**:
- Version detection for each provider
- Comprehensive integration tests with real provider endpoints
- Fallback to error with clear diagnostic message

### Risk 4: SQLite Write Contention
**Risk**: High-throughput token operations could cause SQLite lock contention
**Mitigation**:
- Write-Ahead Logging (WAL) mode for concurrent reads
- In-memory token cache with periodic SQLite sync
- Batched writes for usage statistics

### Risk 5: Single Binary Size
**Risk**: Embedding SQLite and crypto libraries could exceed 50MB target
**Mitigation**:
- UPX compression for release binaries
- Stripped binaries (no debug symbols)
- Link-time optimization (LTO)

## Migration Plan

### Phase 1: MVP (Weeks 1-2)
1. Core HTTP server with OpenAI-compatible endpoints
2. X-User-ID and local user authentication
3. Token encryption and storage (SQLite)
4. OpenAI provider adapter
5. Basic configuration (YAML)

### Phase 2: Multi-Provider (Weeks 3-4)
1. Anthropic Claude adapter
2. Google Gemini adapter
3. OIDC/OAuth2 authentication
4. Web UI dashboard
5. Multi-token load balancing

### Phase 3: Production Hardening (Weeks 5-6)
1. Failover and fallback routing
2. Usage quotas and rate limiting
3. Advanced observability
4. Security audit and penetration testing
5. Performance optimization

### Rollback Strategy
- Each release maintains backward-compatible configuration format
- Binary includes `--version` flag for deployment verification
- Configuration validation on startup (fail-fast on invalid config)
- Graceful shutdown with request completion (no dropped connections)

## Open Questions

1. **Token refresh strategy**: Should we implement proactive refresh (before expiry) or reactive refresh (on auth failure)?
   - Recommendation: Proactive with configurable threshold (default: 30 minutes before expiry)

2. **Metrics retention**: How long should we retain request logs and usage statistics?
   - Recommendation: Ring buffer for request logs (1000 entries), SQLite for usage stats (configurable retention)

3. **Web UI authentication**: Should the admin UI require separate authentication from API requests?
   - Recommendation: Optional password protection, default to localhost-only access

4. **Configuration encryption**: Should we encrypt the configuration file containing provider credentials?
   - Recommendation: Support encrypted config with `--config-key` flag, but default to plaintext for simplicity

5. **Multi-instance coordination**: Should we support multiple proxy instances sharing token state?
   - Recommendation: Not in MVP, document as future enhancement with Redis/etcd backend
