## ADDED Requirements

### Requirement: HTTP Client Initialization
The system SHALL initialize HTTP clients for each configured LLM provider with connection pooling and timeouts.

#### Scenario: Client creation with default settings
- **WHEN** system starts with provider configuration
- **THEN** system creates HTTP client for each provider with default timeout (30s) and connection pool (100)

#### Scenario: Client creation with custom settings
- **WHEN** provider configuration specifies custom timeout and pool size
- **THEN** system creates HTTP client with specified settings

#### Scenario: Client reuse
- **WHEN** multiple requests go to same provider
- **THEN** system reuses existing HTTP client connection pool

### Requirement: Request Construction
The system SHALL construct properly formatted HTTP requests for each LLM provider.

#### Scenario: OpenAI request construction
- **WHEN** system forwards request to OpenAI
- **THEN** system creates POST request to `/v1/chat/completions` with JSON body and `Authorization: Bearer {token}` header

#### Scenario: Anthropic request construction
- **WHEN** system forwards request to Anthropic
- **THEN** system creates POST request to `/v1/messages` with JSON body, `x-api-key: {token}` header, and `anthropic-version: 2023-06-01` header

#### Scenario: Google request construction
- **WHEN** system forwards request to Google
- **THEN** system creates POST request to `/models/{model}:generateContent` with JSON body and API key parameter

### Requirement: Response Handling
The system SHALL handle provider responses and propagate status codes.

#### Scenario: Successful response
- **WHEN** provider returns 200 status
- **THEN** system returns 200 to client with provider response body

#### Scenario: Provider error
- **WHEN** provider returns 4xx or 5xx status
- **THEN** system translates error to OpenAI format and returns appropriate status to client

#### Scenario: Timeout
- **WHEN** provider request exceeds configured timeout
- **THEN** system returns 504 Gateway Timeout to client

#### Scenario: Connection error
- **WHEN** system cannot connect to provider
- **THEN** system returns 502 Bad Gateway to client

### Requirement: Connection Pooling
The system SHALL reuse HTTP connections for efficiency.

#### Scenario: Connection reuse
- **WHEN** sequential requests go to same provider
- **THEN** system reuses TCP connection from pool

#### Scenario: Pool exhaustion
- **WHEN** all connections in pool are busy
- **THEN** system waits for available connection or creates new one up to limit

#### Scenario: Connection cleanup
- **WHEN** connection idle for configured duration
- **THEN** system closes idle connection to free resources
