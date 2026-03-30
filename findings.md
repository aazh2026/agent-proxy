# Findings: PRD Gap Analysis - Complete Report

## Executive Summary

Based on comprehensive code analysis using parallel explore agents, here's the complete gap assessment against PRD.md requirements:

---

## ✅ ALREADY IMPLEMENTED (PRD Compliant)

| Feature | PRD Section | Status | Evidence |
|---------|------------|--------|----------|
| Transparent Auth (X-User-ID) | 4.1.3 | ✅ COMPLETE | Default auth method is `x-user-id`, Authorization header ignored |
| Token Encryption at Rest | 4.3.2.3 | ✅ COMPLETE | AES-256-GCM encryption in handler.go |
| Token Not in API Responses | 4.3.2.3 | ✅ COMPLETE | TokenResponse excludes access_token |
| Chat Completions API | 4.1.1.1 | ✅ COMPLETE | /v1/chat/completions implemented |
| Multi-Provider Support | 4.5 | ✅ COMPLETE | OpenAI, Anthropic, Google |
| Model Routing | 4.4.1 | ✅ COMPLETE | Prefix-based routing (gpt-*, claude-*, gemini-*) |
| Rate Limiting | 4.7.3 | ✅ COMPLETE | middleware/ratelimit.go |
| Quota Tracking | 4.7.2 | ✅ COMPLETE | routing/quota.go |
| Web UI | 4.8 | ✅ COMPLETE | observability/webui.go |
| Config Hot Reload | 4.9 | ✅ COMPLETE | config/watcher.go |
| Token Auto-Refresh | 4.3.2.2 | ✅ COMPLETE | token/refresh.go |

---

## ⚠️ GAPS IDENTIFIED (Require Implementation)

### Gap 1: Embeddings Routing - CRITICAL BUG
**PRD Reference:** 4.1.1.2

**Issue:** OpenAI embeddings routed to wrong endpoint

**Current Behavior:**
- `ForwardingStage.buildURL` (pipeline/forwarding.go lines 101-103) routes ALL OpenAI requests to `/chat/completions`
- Embeddings SHOULD go to `/embeddings` endpoint

**Code Evidence:**
```go
// pipeline/forwarding.go lines 99-110
func (s *ForwardingStage) buildURL(baseURL, provider, model string, stream bool) string {
    switch provider {
    case "openai":
        return baseURL + "/chat/completions"  // ❌ BUG: embeddings also goes here!
    // ...
    }
}
```

**Missing Anthropic Support:**
- `resolveEmbeddingProvider` (embeddings.go lines 158-168) has no Anthropic case
- Anthropic DOES have embeddings API but it's not wired

**Files to Fix:**
- `internal/pipeline/forwarding.go` - Add embeddings path detection
- `internal/api/embeddings.go` - Add Anthropic provider routing

---

### Gap 2: Token Masking in Logs - SECURITY RISK
**PRD Reference:** 4.3.2.3 (全链路脱敏)

**Issue:** MaskToken function exists but is NEVER USED

**Current Behavior:**
- `crypto.MaskToken` defined (crypto.go lines 109-114) but no callers
- Audit logs could potentially leak tokens

**Code Evidence:**
```go
// internal/crypto/crypto.go - UNUSED!
func MaskToken(token string) string {
    if len(token) <= 8 {
        return "****"
    }
    return token[:4] + "****" + token[len(token)-4:]
}
```

**Files to Fix:**
- `internal/token/audit.go` - Integrate MaskToken
- `internal/observability/logger.go` - Add token masking

---

### Gap 3: Streaming Buffer - MINOR DEVIATION
**PRD Reference:** 4.1.2 (零缓冲)

**Issue:** Scanner uses 1MB buffer (not strictly zero-buffer)

**Current Behavior:**
- `bufio.Scanner` with 1MB buffer (streaming.go line 54)
- Per-line flush but internal buffering exists

**Assessment:** Works correctly but has buffer. PRD wants zero-buffer but this is acceptable for most use cases.

---

### Gap 4: Failover/Fallback Wiring - COMPONENT EXISTS, NOT WIRED
**PRD Reference:** 4.4.3

**Issue:** Components exist but are NOT wired into request flow

**What's Implemented:**
- `FailoverHandler` (routing/failover.go) - ✅ Full implementation
- `FallbackRouter` (routing/fallback.go) - ✅ Full implementation
- `TokenSelector` (routing/selector.go) - ✅ Full implementation

**What's NOT Connected:**
- These are ONLY instantiated in tests, NOT in main.go or handlers
- chat.go does NOT use FailoverHandler or FallbackRouter
- This is a FUTURE ENHANCEMENT needed

**Status:** Not implemented in production code

---

### Gap 5: Local Authenticator - INCOMPLETE
**PRD Reference:** 4.2.2.2

**Issue:** No distinct LocalAuthenticator type

**Current Behavior:**
- Config supports "local" method
- Login creates sessions, but no LocalAuthenticator struct

---

## 📋 COMPLETE GAP MATRIX

| # | Gap | PRD Section | Severity | Status | Fix Complexity |
|---|-----|------------|----------|--------|----------------|
| 1 | OpenAI embeddings wrong endpoint | 4.1.1.2 | CRITICAL | BUG | Medium |
| 2 | Anthropic embeddings missing | 4.1.1.2 | HIGH | Missing | Medium |
| 3 | Token masking in logs unused | 4.3.2.3 | MEDIUM | Security | Low |
| 4 | Streaming has 1MB buffer | 4.1.2 | LOW | Deviation | Low |
| 5 | Failover not verified wired | 4.4.3 | HIGH | Unknown | High |
| 6 | LocalAuthenticator missing | 4.2.2.2 | LOW | Incomplete | Low |

---

## 🎯 RECOMMENDED PRIORITY ORDER

1. **P0 (Critical):** Fix OpenAI embeddings endpoint routing
2. **P0 (Critical):** Add Anthropic embeddings support  
3. **P1 (High):** Integrate MaskToken into logs
4. **P2 (Medium):** Verify failover/fallback wiring
5. **P3 (Low):** Fix streaming buffer (if strictly required)
6. **P3 (Low):** Complete LocalAuthenticator

---

## 📁 KEY FILE REFERENCES

### Authentication
- `/internal/auth/auth.go` - XUserIDAuthenticator (transparent auth)
- `/internal/auth/session_auth.go` - SessionAuthenticator

### Token Security  
- `/internal/token/handler.go` - TokenResponse (no token leak)
- `/internal/token/store.go` - Encrypted storage
- `/internal/crypto/crypto.go` - MaskToken (unused!)

### Embeddings (BUGS)
- `/internal/api/embeddings.go` - Handler, routing logic
- `/internal/pipeline/forwarding.go` - WRONG URL for embeddings!

### Streaming
- `/internal/pipeline/streaming.go` - Has 1MB buffer

### Failover/Fallback
- `/internal/routing/failover.go` - Implemented
- `/internal/routing/fallback.go` - Implemented  
- `/internal/routing/selector.go` - Implemented
