## Findings
- Created MVP AB testing scaffolding under internal/abtesting:
  - experiments.go with deterministic variant selection using userID and experiment name
  - splitter.go wrapper for external callers
  - metrics.go scaffolding for per-request metrics
  - config ABExperiment type added under internal/config/abconfig.go
  - API placeholder: internal/api/experiments.go
- Added unit tests for deterministic variant selection
- Added design.md and proposal.md artifacts per OpenSpec workflow

- Next steps: wire with routing layer, integrate with config loading, and implement a lightweight HTTP API for managing experiments.
