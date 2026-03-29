## ADDED Requirements

### Requirement: Cache Key Generation
The system SHALL generate cache keys from request messages.

#### Scenario: Key from messages
- **WHEN** request contains messages
- **THEN** system generates cache key by normalizing and hashing messages

#### Scenario: Key includes model
- **WHEN** generating cache key
- **THEN** key includes model name to prevent cross-model cache hits

### Requirement: Cache Storage
The system SHALL store responses in memory with LRU eviction.

#### Scenario: Cache miss
- **WHEN** cache key not found
- **THEN** system forwards request to provider and caches response

#### Scenario: Cache hit
- **WHEN** cache key found and not expired
- **THEN** system returns cached response without calling provider

#### Scenario: Cache eviction
- **WHEN** cache is full
- **THEN** system evicts least recently used entry

### Requirement: Cache TTL
The system SHALL support configurable time-to-live for cached responses.

#### Scenario: TTL expiration
- **WHEN** cached response exceeds TTL
- **THEN** system treats as cache miss and refreshes

#### Scenario: TTL configuration
- **WHEN** config specifies cache.ttl_seconds
- **THEN** system uses configured TTL

### Requirement: Cache Bypass
The system SHALL support bypassing cache for specific requests.

#### Scenario: Bypass header
- **WHEN** request includes X-Cache-Bypass header
- **THEN** system skips cache and calls provider directly

#### Scenario: No-cache directive
- **WHEN** request includes Cache-Control: no-cache
- **THEN** system skips cache

### Requirement: Cache Statistics
The system SHALL track cache hit/miss statistics.

#### Scenario: Hit tracking
- **WHEN** cache hit occurs
- **THEN** system increments hit counter

#### Scenario: Miss tracking
- **WHEN** cache miss occurs
- **THEN** system increments miss counter

#### Scenario: Stats endpoint
- **WHEN** GET /cache/stats
- **THEN** system returns cache statistics (hits, misses, size)

### Requirement: Cache Management
The system SHALL provide cache management endpoints.

#### Scenario: Clear cache
- **WHEN** POST /cache/clear
- **THEN** system clears all cached responses

#### Scenario: Invalidate by model
- **WHEN** POST /cache/invalidate with model parameter
- **THEN** system clears cache for specified model
