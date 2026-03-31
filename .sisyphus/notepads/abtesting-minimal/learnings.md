## Learnings
- MVP AB testing design can be implemented with a compact in-process deterministic splitter using a stable hash of userID and experiment name.
- We can keep config loading minimal, using a dedicated ABExperiment type to wire with future config loaders.
- Unit tests for deterministic mapping are essential to guard against regressions in the hashing or weighting logic.
