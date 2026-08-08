# Implementation Plan: PII Hashing & Data Privacy

**Branch**: `021-pii-hashing-privacy` | **Date**: 2026-08-08 | **Spec**: [spec.md](file:///home/tarikelmallah/Projects/FlagManagment/specs/021-pii-hashing-privacy/spec.md)

**Input**: Feature specification from `/specs/021-pii-hashing-privacy/spec.md`

## Summary

This feature implements strict PII (Personally Identifiable Information) protection by ensuring that any identity attributes used for targeting or analytics are cryptographically hashed and salted before being stored in the database. As decided in the specification, we will use a hybrid approach: **MurmurHash3** will remain the mechanism for ultra-fast deterministic SDK bucketing (as required by the Constitution), while **SHA-256** combined with an environment-level salt will be used to hash sensitive data (e.g. emails) for persistent storage in audit logs and evaluation analytics. 

## Technical Context

**Language/Version**: Go 1.22
**Primary Dependencies**: `crypto/sha256` (std lib), `github.com/spaolacci/murmur3`
**Storage**: PostgreSQL (GORM)
**Testing**: Go standard testing suite
**Target Platform**: Linux server containerized
**Project Type**: Backend web service & SDK Evaluator
**Performance Goals**: < 1ms flag evaluation execution time.
**Constraints**: Must strictly isolate environment data using unique salts. Must not log plaintext PII.
**Scale/Scope**: Supports unlimited environments.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **II. Environment Isolation**: Passed. The introduction of environment-specific salts explicitly hardens environment isolation.
- **IV. Local Evaluation Performance**: Passed. MurmurHash3 is preserved for the evaluation hot path.
- **VII. PII Protection & Compliance**: Passed. This feature directly implements the PII Protection & Compliance constitutional mandate.
- **Technology Stack Constraints**: Passed. Uses PostgreSQL, Go, and MurmurHash3.

## Project Structure

### Documentation (this feature)

```text
specs/021-pii-hashing-privacy/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
└── tasks.md             # Phase 2 output (future)
```

### Source Code (repository root)

```text
backend/
├── cmd/
├── internal/
│   ├── api/
│   ├── models/        # Update Environment struct to include Salt
│   ├── repository/
│   ├── sdk/           # Update evaluator.go to use SHA-256 + Salt
│   └── services/      # Add background retention cleanup job
├── migrations/        # Add DB migration for 'salt' column in environments
└── tests/
```

**Structure Decision**: The implementation will strictly follow the existing Go backend monolithic structure.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

No constitutional violations exist.

## Proposed Changes

> Note: This section is typically expanded in the task breakdown, but here is the technical blueprint:

1. **Database Migration**: Create `000008_add_environment_salt.up.sql` to add a `salt VARCHAR(64)` column to the `environments` table.
2. **Model Update**: Update `internal/models/environment.go` to include `Salt string` and generate the salt in `BeforeCreate`.
3. **Evaluator Refactor**: Update `internal/sdk/evaluator.go`. The `HashPII` function must accept a `salt` parameter and use `sha256.Sum256([]byte(salt + strVal))`.
4. **API Integration**: Update endpoints in `internal/api/evaluate.go` to retrieve the environment salt and pass it to `HashPII` before passing the context to the analytics service.
5. **Data Retention**: Create a simple periodic background task (or Go ticker in the server startup) to delete evaluation analytics older than 30 days.
