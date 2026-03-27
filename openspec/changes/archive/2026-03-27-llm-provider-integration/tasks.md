## 1. HTTP Client Layer

- [x] 1.1 Create `internal/httpclient/` package structure
- [x] 1.2 Implement HTTP client with configurable timeout and connection pooling
- [x] 1.3 Add provider-specific client constructors (OpenAI, Anthropic, Google)
- [x] 1.4 Implement request builder with authentication header injection
- [x] 1.5 Add response reader with error detection and status code handling
- [x] 1.6 Create connection pool configuration in YAML config
- [ ] 1.7 Add HTTP client unit tests

## 2. Provider Configuration

- [x] 2.1 Extend config structure with `providers` section
- [x] 2.2 Add base URL configuration for each provider
- [x] 2.3 Add timeout configuration per provider
- [x] 2.4 Add enabled/disabled flag per provider
- [x] 2.5 Update example configuration file with provider settings
- [x] 2.6 Add configuration validation for provider settings

## 3. OpenAI Provider Integration

- [x] 3.1 Implement OpenAI HTTP client for chat completions
- [x] 3.2 Implement OpenAI HTTP client for embeddings
- [x] 3.3 Add request forwarding from handler to OpenAI API
- [x] 3.4 Add response passthrough (OpenAI format is compatible)
- [x] 3.5 Implement error handling and status code mapping
- [ ] 3.6 Add integration tests with mock OpenAI server

## 4. Anthropic Provider Integration

- [x] 4.1 Implement Anthropic HTTP client for messages API
- [x] 4.2 Add request transformation from OpenAI to Anthropic format
- [x] 4.3 Handle system message extraction for Anthropic API
- [x] 4.4 Add response transformation from Anthropic to OpenAI format
- [x] 4.5 Implement error handling and status code mapping
- [ ] 4.6 Add integration tests with mock Anthropic server

## 5. Google Gemini Provider Integration

- [x] 5.1 Implement Google Gemini HTTP client
- [x] 5.2 Add request transformation from OpenAI to Gemini format
- [x] 5.3 Add response transformation from Gemini to OpenAI format
- [x] 5.4 Implement authentication (API key in query param or header)
- [x] 5.5 Implement error handling and status code mapping
- [ ] 5.6 Add integration tests with mock Google server

## 6. Request Pipeline

- [x] 6.1 Create `internal/pipeline/` package structure
- [x] 6.2 Implement pipeline stage interface
- [x] 6.3 Create validation stage (request format checking)
- [x] 6.4 Create routing stage (model → provider resolution)
- [x] 6.5 Create token resolution stage (user + provider → token)
- [x] 6.6 Create forwarding stage (HTTP request to provider)
- [x] 6.7 Create response transformation stage
- [x] 6.8 Wire pipeline into chat completions handler
- [x] 6.9 Wire pipeline into embeddings handler
- [ ] 6.10 Add pipeline unit tests

## 7. Streaming Proxy

- [x] 7.1 Implement streaming detection in request pipeline
- [x] 7.2 Create streaming response writer with SSE headers
- [x] 7.3 Implement bidirectional streaming with io.Copy
- [x] 7.4 Add context cancellation propagation
- [x] 7.5 Implement OpenAI SSE stream passthrough
- [x] 7.6 Implement Anthropic SSE stream transformation
- [x] 7.7 Implement Google SSE stream transformation
- [x] 7.8 Add streaming error handling
- [ ] 7.9 Add streaming metrics collection
- [ ] 7.10 Add streaming integration tests

## 8. Error Handling and Resilience

- [x] 8.1 Implement provider error translation to OpenAI format
- [x] 8.2 Add timeout handling for provider requests
- [x] 8.3 Add connection error handling
- [ ] 8.4 Implement retry logic with exponential backoff
- [ ] 8.5 Add circuit breaker for failing providers (future)
- [x] 8.6 Add error logging with request context

## 9. Token Integration

- [x] 9.1 Wire token resolver into request pipeline
- [x] 9.2 Implement token selection based on user and provider
- [x] 9.3 Add token decryption for provider authentication
- [x] 9.4 Implement token memory clearing after use
- [x] 9.5 Add multi-token load balancing (round-robin)
- [ ] 9.6 Add token failover on authentication error

## 10. Testing and Validation

- [ ] 10.1 Create mock HTTP servers for each provider
- [ ] 10.2 Add end-to-end tests for chat completions
- [ ] 10.3 Add end-to-end tests for embeddings
- [ ] 10.4 Add end-to-end tests for streaming
- [ ] 10.5 Add error scenario tests (timeout, auth failure, rate limit)
- [ ] 10.6 Add performance benchmarks for proxy latency
- [ ] 10.7 Test with real provider APIs (manual validation)

## 11. Documentation

- [x] 11.1 Update README with provider setup instructions
- [x] 11.2 Add provider configuration documentation
- [x] 11.3 Document error codes and troubleshooting
- [x] 11.4 Add examples for each provider
