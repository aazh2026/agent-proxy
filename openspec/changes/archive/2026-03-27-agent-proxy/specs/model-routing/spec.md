## ADDED Requirements

### Requirement: Model Name Routing
The system SHALL route requests to appropriate provider based on model name prefix matching.

#### Scenario: Default routing rules
- **WHEN** request model starts with "gpt-"
- **THEN** system routes to OpenAI provider

- **WHEN** request model starts with "claude-"
- **THEN** system routes to Anthropic Claude provider

- **WHEN** request model starts with "gemini-"
- **THEN** system routes to Google Gemini provider

#### Scenario: Custom model alias
- **WHEN** configuration defines alias `my-gpt-4: { provider: "openai", model: "gpt-4" }`
- **WHEN** request uses model "my-gpt-4"
- **THEN** system routes to OpenAI provider with model "gpt-4"

#### Scenario: Custom routing priority
- **WHEN** custom routing rules defined in configuration
- **THEN** custom rules take precedence over default prefix rules

#### Scenario: No matching route
- **WHEN** request model matches no routing rule
- **THEN** system returns 404 error indicating model not found

### Requirement: Token Selection Strategy
The system SHALL support multiple strategies for selecting tokens when multiple tokens available for same provider.

#### Scenario: Round-robin selection
- **WHEN** configuration specifies `token_strategy: "round-robin"`
- **WHEN** multiple enabled tokens available for provider
- **THEN** system distributes requests evenly across tokens in circular order

#### Scenario: Weighted selection
- **WHEN** configuration specifies `token_strategy: "weighted"`
- **WHEN** tokens have priority weights configured
- **THEN** system selects tokens proportionally to their weights

#### Scenario: Priority selection
- **WHEN** configuration specifies `token_strategy: "priority"`
- **WHEN** multiple tokens available with different priorities
- **THEN** system always selects highest priority enabled token

#### Scenario: Fallback selection
- **WHEN** highest priority token fails or is disabled
- **WHEN** other tokens available
- **THEN** system falls back to next priority token

### Requirement: Multi-Token Load Balancing
The system SHALL distribute requests across multiple tokens for same provider.

#### Scenario: Load distribution
- **WHEN** user has 3 enabled tokens for OpenAI
- **WHEN** load balancing enabled
- **THEN** system distributes requests across all 3 tokens

#### Scenario: Token health tracking
- **WHEN** token returns error
- **THEN** system marks token as unhealthy for configurable cooldown period

#### Scenario: Token recovery
- **WHEN** cooldown period expires
- **THEN** system marks token as healthy and includes in load balancing

### Requirement: Failover Within Provider
The system SHALL automatically retry failed requests with alternative tokens from same provider.

#### Scenario: Automatic token failover
- **WHEN** request with token A fails with retryable error (429, 5xx)
- **WHEN** user has token B for same provider
- **THEN** system automatically retries request with token B

#### Scenario: Failover exhausted
- **WHEN** all tokens for provider fail
- **THEN** system returns error from last failed attempt

#### Scenario: Failover configuration
- **WHEN** configuration specifies `max_retries: 3`
- **THEN** system retries up to 3 times with different tokens

#### Scenario: Non-retryable errors
- **WHEN** request fails with 400 Bad Request
- **THEN** system does not retry (client error, not provider error)

### Requirement: Cross-Provider Fallback
The system SHALL support fallback to alternative providers when configured.

#### Scenario: Fallback chain configuration
- **WHEN** configuration specifies fallback: `gpt-4 -> claude-3-opus -> gemini-1.5-pro`
- **WHEN** OpenAI provider fails for gpt-4
- **THEN** system retries with Claude provider using claude-3-opus

#### Scenario: Fallback success
- **WHEN** fallback provider succeeds
- **WHEN** response received
- **THEN** system returns fallback provider response to client

#### Scenario: Fallback exhausted
- **WHEN** all providers in fallback chain fail
- **THEN** system returns error from last failed attempt

#### Scenario: Fallback disabled
- **WHEN** configuration disables fallback
- **THEN** system does not attempt cross-provider fallback

### Requirement: Request Retry Logic
The system SHALL support configurable retry logic for transient failures.

#### Scenario: Retryable error codes
- **WHEN** request fails with 429, 500, 502, 503, or 504
- **WHEN** retries remaining
- **THEN** system retries request

#### Scenario: Retry delay strategies
- **WHEN** configuration specifies retry delay strategy
- **THEN** system applies configured delay between retries: fixed, linear backoff, or exponential backoff

#### Scenario: Max retry delay
- **WHEN** calculated retry delay exceeds configured maximum
- **THEN** system caps delay at maximum value

#### Scenario: Retry limit enforcement
- **WHEN** retry count exceeds configured maximum
- **THEN** system stops retrying and returns final error

### Requirement: Routing Configuration
The system SHALL support flexible routing configuration via YAML.

#### Scenario: Routing rules definition
- **WHEN** configuration includes routing section
- **THEN** system loads custom model mappings, token strategies, and fallback rules

#### Scenario: Routing hot reload
- **WHEN** routing configuration changes
- **THEN** system reloads routing rules without restart

#### Scenario: Invalid routing configuration
- **WHEN** configuration contains invalid routing rules
- **THEN** system fails to start with validation error
