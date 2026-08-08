# Implementation Plan: Contextual Targeting Engine

**Branch**: `012-contextual-targeting` | **Date**: 2026-07-25 | **Spec**: [spec.md](../spec.md)

**Input**: Feature specification from `/specs/012-contextual-targeting/spec.md`

## Summary

This feature introduces a contextual targeting engine allowing users to define rules (Equals, Contains, Regex) for feature flags. It stores these rules in the existing JSONB column and evaluates them natively in the SDK and backend using a high-performance matching algorithm with regex pattern caching.

## Technical Context

**Language/Version**: Go 1.21+, Node.js (React)

**Primary Dependencies**: Go's native `regexp`, `strings` packages, standard React for frontend

**Storage**: PostgreSQL (`environment_flag_states.targeting_rules` JSONB)

**Testing**: Go unit tests (fast evaluation benchmarks)

**Target Platform**: Linux server, modern web browsers

**Project Type**: Backend service + React frontend

**Performance Goals**: < 1ms evaluation time for SDK per rule

**Constraints**: Must pre-validate Regex upon save to prevent panics or performance crashes at runtime.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **API-First Contract Design**: Passes. The JSON payload contract for updating rules is defined in `data-model.md`.
- **Environment Isolation**: Passes. Rules are inherently isolated by being stored in `environment_flag_states`.
- **Local Evaluation Performance (NON-NEGOTIABLE)**: Passes. Native Go evaluation ensures < 1ms execution, and regex strings will be pre-compiled and validated.
- **Test-First Quality Gates**: Passes. Tests will be added for the evaluation engine prior to exposing the API logic.
- **PII Protection & Compliance**: Passes. The evaluation engine operates in-memory; context (like emails) is not logged to disk or stored persistently unless specifically tracked (out of scope here).

## Project Structure

### Documentation (this feature)

```text
specs/012-contextual-targeting/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
└── tasks.md             # Phase 2 output (to be generated next)
```

### Source Code (repository root)

```text
backend/
├── internal/
│   ├── models/
│   │   └── environment_flag_state.go # Add Rule structs
│   ├── sdk/
│   │   └── evaluator.go              # New local evaluation engine
│   └── services/
│       └── flag_state_service.go     # Add regex validation
└── tests/

frontend/
├── src/
│   ├── components/
│   │   └── flagStates/
│   │       └── TargetingRuleBuilder.tsx # New UI component
│   └── services/
└── tests/
```

**Structure Decision**: Standard web application structure (Option 2). The logic is primarily housed in the `backend/internal/sdk` package (for evaluation) and `frontend/src/components/flagStates` (for the builder UI).

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

*No violations detected.*
