# OpenSpec Proposal: abtesting-minimal

- Change: abtesting-minimal
- Schema: spec-driven
- Objective: Introduce a minimal AB testing framework (config-driven experiments with two variants) integrated into the project with a small API surface, metrics scaffolding, and sticky user assignment.

## Context
- The PRD section for P2 requires an AB testing capability that can route users to variant A/B with deterministic stickiness per user and collect lightweight metrics per request.
- The new module should live under internal/abtesting and integrate with internal/config for experiment configuration and with routing hooks for selection at request time.

## Problem Statement
- Without an AB framework, feature experiments are hard to manage, track, and verify deterministically across traffic.
- We need a small, dependency-free layer that can be wired into routing and config, and provide an API surface to list/start/stop experiments when feasible in scope.

## Goals & Success Criteria
- A minimal AB testing module at internal/abtesting with:
  - Deterministic sticky assignment per user for given experiment
  - Traffic split across two or more variants according to config
  - Basic per-request metrics scaffolding (latency, cost, success)
  - A small API surface to list/start/stop experiments (optional exposure)
- Integration hooks wired into existing config and routing points (non-invasive).
- Tests for sticky assignment correctness and deterministic routing per user.

## Scope (In-Scope)
- Implement internal/abtesting with:
  - experiments.go: core Experiment, Variant definitions, and config-driven loading
  - splitter.go: deterministic variant selection function based on user id and experiment name
  - metrics.go: scaffolding for per-request metrics capture
  - config.go (new abridged config): Unrealized, but provide minimal interfaces to load AB experiments from config files
- Basic API surface: internal/api/experiments.go (optional) to expose list/start/stop endpoints
- Wire minimal tests for sticky function and traffic split

## Non-Goals
- Full-featured analytics backend or external dependencies
- Deep integration with auth, token providers, or migrations beyond skeleton wiring
- Production-grade A/B metrics pipeline

## Acceptance Criteria (Definition of Done)
- Code builds locally with go build ./...
- Unit tests pass (sticky assignment deterministic for given user and config)
- Sample config yields deterministic routing per user for an experiment with two variants
- No changes to core authentication, token adapters, or provider code beyond wiring hooks

## Risks & Mitigations
- Risk: Over-fitting to two variants only. Mitigation: design the API to support multiple variants later without breaking changes.
- Risk: Missing dependencies in integration points. Mitigation: implement small adapter layer and keep wiring minimal.

## Open Questions
- How will configuration be supplied in CI/test environments? We'll implement a simple YAML-config loader as a starting point.
- Do we want to expose a live API for managers to start/stop experiments in this MVP? We can add a lightweight http handler later if needed.

## Next Steps (Post-Approval)
- Create design.md detailing architecture and wiring points
- Implement the abtesting package and config loader
- Add unit tests for sticky routing and traffic split
- Wire into routing at a minimal touchpoint (e.g., a middleware hook)

---
This proposal is a starting point. If approved, I will generate the remaining artifacts (design.md, specs, tasks) and implement the MVP.
