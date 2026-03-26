## ADDED Requirements

### Requirement: Chat Completions Endpoint
The system SHALL provide a POST endpoint at `/v1/chat/completions` that accepts requests in OpenAI-compatible format and returns responses in OpenAI-compatible format.

#### Scenario: Non-streaming request
- **WHEN** client sends POST to `/v1/chat/completions` with valid JSON body containing `model` and `messages` fields
- **THEN** system returns JSON response in OpenAI chat completion format with `id`, `object`, `created`, `model`, `choices`, and `usage` fields

#### Scenario: Streaming request
- **WHEN** client sends POST to `/v1/chat/completions` with `stream: true` in request body
- **THEN** system returns Server-Sent Events (SSE) stream with `data:` prefixed JSON chunks in OpenAI streaming format, ending with `data: [DONE]`

#### Scenario: Invalid request body
- **WHEN** client sends POST to `/v1/chat/completions` with missing required fields (`model` or `messages`)
- **THEN** system returns 400 Bad Request with OpenAI-compatible error response

#### Scenario: Unsupported model
- **WHEN** client sends request with model name that has no configured provider mapping
- **THEN** system returns 404 Not Found with error message indicating model not found

### Requirement: Embeddings Endpoint
The system SHALL provide a POST endpoint at `/v1/embeddings` that accepts requests in OpenAI-compatible format and returns responses in OpenAI-compatible format.

#### Scenario: Single input embedding
- **WHEN** client sends POST to `/v1/embeddings` with `input` as string and valid `model`
- **THEN** system returns JSON response with embedding vector in OpenAI embeddings format

#### Scenario: Batch input embedding
- **WHEN** client sends POST to `/v1/embeddings` with `input` as array of strings
- **THEN** system returns JSON response with array of embedding vectors corresponding to each input

#### Scenario: Invalid embedding request
- **WHEN** client sends POST to `/v1/embeddings` with missing `input` or `model` field
- **THEN** system returns 400 Bad Request with OpenAI-compatible error response

### Requirement: Request Validation
The system SHALL validate incoming requests against OpenAI API specification without modifying valid parameters.

#### Scenario: Valid parameters pass through
- **WHEN** client sends request with valid OpenAI parameters (model, messages, temperature, top_p, max_tokens, stream, etc.)
- **THEN** system passes all parameters unchanged to upstream provider

#### Scenario: Unknown parameters rejected
- **WHEN** client sends request with parameters not in OpenAI specification
- **THEN** system returns 400 Bad Request with error indicating unknown parameter (strict mode) or passes through (permissive mode based on configuration)

#### Scenario: Parameter type validation
- **WHEN** client sends request with incorrect parameter types (e.g., string for temperature)
- **THEN** system returns 400 Bad Request with error indicating type mismatch

### Requirement: Error Response Standardization
The system SHALL return all errors in OpenAI-compatible error format with appropriate HTTP status codes.

#### Scenario: Proxy error format
- **WHEN** system encounters internal error
- **THEN** system returns JSON with structure: `{"error": {"message": "...", "type": "...", "code": "..."}}` and appropriate HTTP status

#### Scenario: Upstream error forwarding
- **WHEN** upstream provider returns error
- **THEN** system translates error to OpenAI format and forwards with original HTTP status code

#### Scenario: Authentication error
- **WHEN** request fails authentication
- **THEN** system returns 401 Unauthorized with error message indicating authentication failure

#### Scenario: Rate limit error
- **WHEN** request exceeds rate limit
- **THEN** system returns 429 Too Many Requests with error message indicating rate limit exceeded

### Requirement: Streaming Response Format
The system SHALL support Server-Sent Events (SSE) streaming for chat completions with zero buffering and immediate chunk forwarding.

#### Scenario: SSE chunk format
- **WHEN** upstream provider sends streaming response chunk
- **THEN** system forwards chunk immediately with format: `data: {json}\n\n` where json matches OpenAI streaming response format

#### Scenario: Stream completion
- **WHEN** upstream provider completes streaming response
- **THEN** system sends final `data: [DONE]\n\n` to signal stream completion

#### Scenario: Client disconnection during streaming
- **WHEN** client disconnects during streaming response
- **THEN** system immediately closes connection to upstream provider and terminates request

#### Scenario: Stream error during transmission
- **WHEN** upstream provider encounters error during streaming
- **THEN** system forwards error in OpenAI-compatible format and terminates stream
