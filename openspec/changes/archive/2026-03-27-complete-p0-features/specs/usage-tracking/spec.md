## ADDED Requirements

### Requirement: Request Counting
The system SHALL count requests per user.

#### Scenario: Request counted
- **WHEN** user makes request
- **THEN** system increments user's request counter

#### Scenario: Request counted by model
- **WHEN** user makes request for specific model
- **THEN** system increments counter for user and model

### Requirement: Token Counting
The system SHALL count token consumption per user.

#### Scenario: Tokens counted
- **WHEN** request completes with token usage
- **THEN** system adds prompt_tokens, completion_tokens, total_tokens to user's counters

#### Scenario: Tokens counted by provider
- **WHEN** request completes
- **THEN** system tracks tokens per provider

### Requirement: Cost Estimation
The system SHALL estimate cost per request.

#### Scenario: Cost calculated
- **WHEN** request completes
- **THEN** system calculates estimated cost based on model pricing and token usage

#### Scenario: Cost accumulated
- **WHEN** multiple requests complete
- **THEN** system accumulates cost per user

### Requirement: Usage Persistence
The system SHALL persist usage statistics.

#### Scenario: Usage saved to database
- **WHEN** request completes
- **THEN** system saves usage data to usage_stats table

#### Scenario: Usage queryable
- **WHEN** admin queries usage
- **THEN** system returns usage statistics for specified time period
