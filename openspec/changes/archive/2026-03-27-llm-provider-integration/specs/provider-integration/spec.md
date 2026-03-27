## ADDED Requirements

### Requirement: OpenAI Provider Integration
The system SHALL forward requests to OpenAI API and return responses in OpenAI-compatible format.

#### Scenario: Chat completions request
- **WHEN** client sends POST to `/v1/chat/completions` with model `gpt-4`
- **WHEN** user has configured OpenAI token
- **THEN** system forwards request to OpenAI API and returns response to client

#### Scenario: Streaming chat completions
- **WHEN** client sends POST to `/v1/chat/completions` with `stream: true`
- **THEN** system streams response from OpenAI API to client in real-time

#### Scenario: Embeddings request
- **WHEN** client sends POST to `/v1/embeddings` with model `text-embedding-ada-002`
- **THEN** system forwards request to OpenAI API and returns embedding vectors

#### Scenario: OpenAI authentication injection
- **WHEN** system forwards request to OpenAI
- **THEN** system injects `Authorization: Bearer {access_token}` header from user's token

### Requirement: Anthropic Provider Integration
The system SHALL forward requests to Anthropic API with protocol translation.

#### Scenario: Chat completions request
- **WHEN** client sends POST to `/v1/chat/completions` with model `claude-3-opus`
- **WHEN** user has configured Anthropic token
- **THEN** system translates request to Anthropic format, forwards to API, translates response back

#### Scenario: System message handling
- **WHEN** client sends messages with `role: system`
- **THEN** system extracts system message and includes in Anthropic `system` parameter

#### Scenario: Streaming chat completions
- **WHEN** client sends POST with `stream: true` for Claude model
- **THEN** system streams response from Anthropic API, translating chunks to OpenAI format

#### Scenario: Anthropic authentication injection
- **WHEN** system forwards request to Anthropic
- **THEN** system injects `x-api-key: {access_token}` and `anthropic-version: 2023-06-01` headers

### Requirement: Google Gemini Provider Integration
The system SHALL forward requests to Google Gemini API with protocol translation.

#### Scenario: Chat completions request
- **WHEN** client sends POST to `/v1/chat/completions` with model `gemini-pro`
- **WHEN** user has configured Google token
- **THEN** system translates request to Gemini format, forwards to API, translates response back

#### Scenario: Streaming chat completions
- **WHEN** client sends POST with `stream: true` for Gemini model
- **THEN** system streams response from Gemini API, translating chunks to OpenAI format

#### Scenario: Google authentication injection
- **WHEN** system forwards request to Google
- **THEN** system injects API key as query parameter or Authorization header per Gemini spec

### Requirement: Provider Configuration
The system SHALL support configurable provider endpoints and settings.

#### Scenario: Base URL configuration
- **WHEN** configuration specifies `providers.openai.base_url`
- **THEN** system uses configured URL instead of default

#### Scenario: Timeout configuration
- **WHEN** configuration specifies provider timeout
- **THEN** system enforces timeout on requests to that provider

#### Scenario: Provider disabled
- **WHEN** configuration marks provider as disabled
- **THEN** system does not route requests to disabled provider

### Requirement: Error Translation
The system SHALL translate provider errors to OpenAI-compatible format.

#### Scenario: Authentication error
- **WHEN** provider returns 401 Unauthorized
- **THEN** system returns OpenAI-format error with type `authentication_error`

#### Scenario: Rate limit error
- **WHEN** provider returns 429 Too Many Requests
- **THEN** system returns OpenAI-format error with type `rate_limit_error`

#### Scenario: Model not found
- **WHEN** provider returns 404 for unknown model
- **THEN** system returns OpenAI-format error with type `invalid_request_error`

#### Scenario: Server error
- **WHEN** provider returns 5xx error
- **THEN** system returns OpenAI-format error with type `server_error`
