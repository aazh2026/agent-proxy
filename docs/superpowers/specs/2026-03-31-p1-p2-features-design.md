# P1/P2 Features Design - Semantic Caching, Cost-Aware Routing, and Analytics

**Date**: 2026-03-31  
**Status**: Approved  
**Author**: Engineering Team  

---

## 1. Semantic Caching

### 1.1 Overview

Semantic caching reduces LLM costs and latency by detecting and reusing responses from semantically similar requests. Unlike exact-match caching, semantic caching uses embeddings to find similar prompts even when the exact wording differs.

### 1.2 Architecture

```
Client Request
      ↓
┌─────────────────────────┐
│  Generate Embedding      │
│  (using ada-002)        │
└─────────────────────────┘
      ↓
┌─────────────────────────┐
│  Similarity Search      │
│  (threshold-based)      │
└─────────────────────────┘
      ↓
   ┌────┴────┐
   ↓          ↓
HIT          MISS
   ↓          ↓
Cached   ┌─────────────┐
Response │ Call LLM   │
         │ Provider   │
         └─────────────┘
```

### 1.3 Configuration

```yaml
cache:
  semantic:
    enabled: true
    similarity_threshold: 0.95    # 0.0-1.0, higher = stricter match
    embedding_model: "text-embedding-3-small"
    ttl_seconds: 3600
    max_size: 1000
    max_response_size: 1048576   # 1MB
    providers:                   # Which providers support semantic cache
      - openai
      - anthropic
      - google
```

### 1.4 Data Structures

**Cache Entry**:
```go
type SemanticCacheEntry struct {
    Key           string    `json:"key"`            // SHA256 hash of normalized prompt
    Model         string    `json:"model"`          // Original model requested
    Provider      string    `json:"provider"`       // Provider used
    PromptHash    string    `json:"prompt_hash"`    // Embedding reference
    Response      []byte    `json:"response"`       // Cached response
    TokensUsed    int       `json:"tokens_used"`     // Original token count (for cost)
    CreatedAt     time.Time `json:"created_at"`
    HitCount      int       `json:"hit_count"`
}
```

**Embedding Store**:
```go
type EmbeddingStore struct {
    PromptHash string    `json:"prompt_hash"`   // Reference key
    Embedding  []float32 `json:"embedding"`      // Vector embedding
    Normalized string    `json:"normalized"`    // Normalized prompt text
}
```

### 1.5 Implementation Details

1. **Embedding Generation**:
   - Use provider's embedding API (OpenAI ada-2 by default)
   - Normalize prompt: lowercase, trim whitespace, remove extra spaces
   - Cache embedding alongside response

2. **Similarity Calculation**:
   - Cosine similarity between prompt embeddings
   - Configurable threshold (default 0.95)
   - For performance: approximate nearest neighbor (AN) for large caches

3. **Cache Invalidation**:
   - TTL-based (configurable)
   - Manual invalidation by model prefix
   - Auto-invalidate when model updates

4. **Storage**:
   - Embeddings in memory for fast lookup
   - Responses in SQLite with encryption
   - Separate storage from exact-match cache

### 1.6 API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/cache/semantic` | GET | Get semantic cache stats |
| `/cache/semantic` | DELETE | Clear semantic cache |
| `/cache/semantic/:key` | DELETE | Invalidate specific entry |

### 1.7 Constraints

- Only cache `stream: false` requests (streaming responses are non-deterministic)
- Exclude requests with `temperature > 0` or `seed` parameter
- Maximum response size configurable to prevent memory bloat

---

## 2. Cost-Aware Routing

### 2.1 Overview

Cost-aware routing automatically selects the optimal LLM provider/model based on configured cost priorities, ensuring cost efficiency while meeting quality requirements.

### 2.2 Configuration

```yaml
routing:
  cost_strategy: "quality-first"  # quality-first, cost-first, latency-first, balanced
  
  # Cost per 1M tokens (USD)
  cost_matrix:
    gpt-4o:
      input: 5.00
      output: 15.00
    gpt-4o-mini:
      input: 0.30
      output: 1.20
    gpt-3.5-turbo:
      input: 0.50
      output: 1.50
    claude-3-opus:
      input: 15.00
      output: 75.00
    claude-3.5-sonnet:
      input: 3.00
      output: 15.00
    gemini-1.5-pro:
      input: 1.25
      output: 5.00

  # Fallback chains (tried in order)
  fallback_chains:
    - [gpt-4o, gpt-4o-mini, gpt-3.5-turbo]
    - [claude-3.5-sonnet, claude-3-haiku]
    - [gemini-1.5-pro, gemini-1.5-flash]

  # User cost preferences (override global)
  user_preferences:
    alice:
      strategy: "cost-first"
    bob:
      strategy: "balanced"
```

### 2.2 Strategy Definitions

| Strategy | Behavior |
|----------|----------|
| `quality-first` | Always use highest quality model available |
| `cost-first` | Use cheapest model that meets minimum quality threshold |
| `latency-first` | Prefer fastest responding model |
| `balanced` | Weighted score: 40% cost, 30% latency, 30% quality |

### 2.3 Request Flow

```
Client Request (model: gpt-4o)
         ↓
   Check user preference
         ↓
   Apply cost strategy
         ↓
   Lookup cost matrix
         ↓
   Select optimal model
         ↓
   Execute with fallback chain
```

### 2.4 Cost Calculation

```go
func CalculateRequestCost(promptTokens, completionTokens int, costMatrix map[string]CostConfig) float64 {
    inputCost := float64(promptTokens) / 1_000_000 * costMatrix[model].Input
    outputCost := float64(completionTokens) / 1_000_000 * costMatrix[model].Output
    return inputCost + outputCost
}
```

### 2.5 Implementation

- Cost matrix loaded from config at startup
- Per-request cost tracking in metrics
- User-level override via `X-Cost-Strategy` header or config

---

## 3. Per-User Cost Quotas

### 3.1 Overview

Track and control LLM usage costs per user with configurable quotas, alerts, and enforcement actions.

### 3.2 Configuration

```yaml
quota:
  enabled: true
  
  # Global default quota
  default:
    period: "monthly"    # daily, weekly, monthly
    tokens_limit: 1000000     # 1M tokens
    cost_limit: 100.00        # USD
    requests_limit: 10000

  # Actions when quota exceeded
  exceeded_action: "block"    # block, warn, fallback
  
  # Alert thresholds (send notification)
  alerts:
    - threshold: 80    # % of quota used
      action: "warn"
    - threshold: 95
      action: "warn"

  # Per-user overrides
  user_quotas:
    alice:
      tokens_limit: 5000000
      cost_limit: 500.00
    team-engineering:
      tokens_limit: 10000000
      cost_limit: 1000.00
```

### 3.3 Data Model

```sql
-- User quota tracking
CREATE TABLE user_quotas (
    id INTEGER PRIMARY KEY,
    user_id TEXT NOT NULL,
    period_type TEXT NOT NULL,  -- daily, weekly, monthly
    period_start INTEGER NOT NULL,  -- Unix timestamp
    tokens_used INTEGER DEFAULT 0,
    cost_used REAL DEFAULT 0.0,
    requests_used INTEGER DEFAULT 0,
    tokens_limit INTEGER NOT NULL,
    cost_limit REAL NOT NULL,
    requests_limit INTEGER NOT NULL,
    created_at INTEGER,
    updated_at INTEGER,
    UNIQUE(user_id, period_type, period_start)
);

-- Quota alerts log
CREATE TABLE quota_alerts (
    id INTEGER PRIMARY KEY,
    user_id TEXT NOT NULL,
    period_type TEXT NOT NULL,
    threshold INTEGER NOT NULL,
    alert_type TEXT NOT NULL,  -- warning, exceeded
    sent_at INTEGER NOT NULL
);
```

### 3.4 Quota Enforcement

| Action | Behavior |
|--------|----------|
| `block` | Return 429 Too Many Requests with quota info |
| `warn` | Allow request but add warning header `X-Quota-Warn: approaching/exceeded` |
| `fallback` | Auto-switch to cheaper model (requires cost_strategy: cost-first) |

### 3.5 API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/quota` | GET | Get current user's quota status |
| `/quota/:user_id` | GET | Admin: get specific user's quota (requires admin auth) |
| `/quota/:user_id` | PUT | Admin: update user's quota |
| `/quota/reset/:user_id` | POST | Admin: manually reset quota |

### 3.6 Dashboard Integration

- Show per-user cost trends (graph)
- Display quota usage bar
- Alert indicators when approaching limit
- Export quota reports (CSV)

---

## 4. A/B Testing Framework

### 4.1 Overview

Split traffic between different models/providers to compare performance and enable gradual rollouts.

### 4.2 Configuration

```yaml
ab_testing:
  enabled: true
  
  experiments:
    - id: "gpt-4o-vs-sonnet"
      description: "Compare GPT-4o vs Claude Sonnet for coding tasks"
      traffic_split:
        gpt-4o: 50
        claude-3.5-sonnet: 50
      # Models must be in cost_matrix
      models:
        - gpt-4o
        - claude-3.5-sonnet
      sticky: true    # Same user always gets same variant
      metrics:
        - latency
        - success_rate
        - cost
        - user_satisfaction  # via feedback endpoint
      
    - id: "fast-model-rollout"
      description: "Gradual rollout of gpt-4o-mini"
      traffic_split:
        gpt-4o: 90
        gpt-4o-mini: 10
      models:
        - gpt-4o
        - gpt-4o-mini
      sticky: true
```

### 4.3 Traffic Assignment

```go
func AssignVariant(userID, experimentID string, split map[string]int) string {
    // Sticky assignment using deterministic hash
    hash := md5.Sum([]byte(userID + experimentID))
    sum := 0
    for _, b := range hash {
        sum += int(b)
    }
    
    cumulative := 0
    for model, percent := range split {
        cumulative += percent
        if sum%100 < cumulative {
            return model
        }
    }
    return ""  // No assignment
}
```

### 4.4 Metrics Collection

```sql
CREATE TABLE experiment_metrics (
    id INTEGER PRIMARY KEY,
    experiment_id TEXT NOT NULL,
    variant TEXT NOT NULL,
    user_id TEXT NOT NULL,
    model TEXT NOT NULL,
    latency_ms INTEGER,
    tokens_used INTEGER,
    cost REAL,
    success BOOLEAN,
    timestamp INTEGER NOT NULL
);
```

### 4.5 Statistical Analysis

- Calculate: conversion rate, mean latency, cost per request
- Chi-square test for significance (p < 0.05)
- Dashboard shows comparison charts

### 4.6 API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/experiments` | GET | List active experiments |
| `/experiments/:id` | GET | Get experiment results |
| `/experiments/:id/start` | POST | Start experiment |
| `/experiments/:id/stop` | POST | Stop experiment |
| `/feedback` | POST | Submit user feedback for A/B test |

---

## 5. Request Parameter Override

### 5.1 Overview

Allow administrators to enforce or default certain parameters per user, model, or globally.

### 5.2 Configuration

```yaml
request_overrides:
  # Global defaults (applied to all requests)
  defaults:
    temperature: 0.7
    top_p: 0.9
  
  # Forced overrides (cannot be overridden by client)
  forced:
    # Per model
    gpt-4o:
      temperature: 0.5
      max_tokens: 4096
    # Per user
    alice:
      max_tokens: 2048
  
  # Allow client override
  allowed_override:
    - temperature
    - top_p
    - max_tokens
  
  # Block client override
  blocked_override:
    - response_format  # Force JSON mode for certain users
```

---

## 6. Implementation Priority

| Priority | Feature | Complexity | Estimated Effort |
|----------|---------|------------|------------------|
| P1 | Cost Matrix + Routing | Medium | 1 week |
| P1 | Per-User Quotas | Medium | 1 week |
| P1 | Basic Semantic Cache | Medium | 1-2 weeks |
| P2 | A/B Testing | High | 2 weeks |
| P2 | Parameter Override | Low | 3-5 days |

---

## 7. Backward Compatibility

- All new features disabled by default
- Existing configs continue to work without changes
- Feature flags in config to enable incrementally

---

## 8. Security Considerations

- Cost matrix is read-only (no API exposure)
- Quota data encrypted at rest
- A/B assignment uses HMAC for tamper-proof sticky sessions
- Admin endpoints require admin password authentication
