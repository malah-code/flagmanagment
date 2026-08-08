# Implementation Plan: [FEATURE]

**Branch**: `[###-feature-name]` | **Date**: [DATE] | **Spec**: [link]

**Input**: Feature specification from `/specs/[###-feature-name]/spec.md`

**Note**: This template is filled in by the `/speckit-plan` command; its definition describes the execution workflow.

## Summary

Implement Immutable Audit Logs and SIEM Webhooks to provide an append-only ledger of all administrative actions. The system will sanitize PII and sensitive keys from JSON payloads before insertion, expose APIs for querying and CSV export (via streaming), and dispatch webhook notifications asynchronously via in-memory goroutine queues with backoff retries.

## Technical Context

<!--
  ACTION REQUIRED: Replace the content in this section with the technical details
  for the project. The structure here is presented in advisory capacity to guide
  the iteration process.
-->

**Language/Version**: Go 1.22+

**Primary Dependencies**: `encoding/csv` (standard library), standard Go HTTP router

**Storage**: PostgreSQL (new `audit_logs` and `webhook_integrations` tables)

**Testing**: `go test` with integration tests for webhook delivery

**Target Platform**: Linux server, Docker containers

**Project Type**: Backend Service Module (`backend/internal/audit` and `backend/internal/webhook`)

**Performance Goals**: Asynchronous webhook dispatch without blocking the API response; CSV streaming to prevent OOM on large exports.

**Constraints**: Strict append-only constraints; generic PII sanitization for arbitrary JSONB payloads.

**Scale/Scope**: Audit log table will grow unbounded; pagination must be efficient using indexed timestamp/actor filtering.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **VII. PII Protection & Compliance**: PASS - We are explicitly designing a `Scrub()` phase before DB insertion to sanitize API keys, passwords, and tokens.
- **III. Governance by Default**: PASS - This feature directly implements the immutable audit logging requirement for all deployments.
- **IV. Local Evaluation Performance**: PASS - We explicitly bypass audit logging for read operations and SDK evaluations to prevent latency and DB bloat.
- **VIII. Cloud-Native Portability**: PASS - The webhook dispatcher uses an in-memory goroutine approach instead of adding a new external dependency like RabbitMQ, maintaining the lightweight Docker Compose portability.

## Project Structure

### Documentation (this feature)

```text
specs/[###-feature]/
├── plan.md              # This file (/speckit-plan command output)
├── research.md          # Phase 0 output (/speckit-plan command)
├── data-model.md        # Phase 1 output (/speckit-plan command)
├── quickstart.md        # Phase 1 output (/speckit-plan command)
├── contracts/           # Phase 1 output (/speckit-plan command)
└── tasks.md             # Phase 2 output (/speckit-tasks command - NOT created by /speckit-plan)
```

### Source Code (repository root)
<!--
  ACTION REQUIRED: Replace the placeholder tree below with the concrete layout
  for this feature. Delete unused options and expand the chosen structure with
  real paths (e.g., apps/admin, packages/something). The delivered plan must
  not include Option labels.
-->

```text
backend/
├── internal/
│   ├── audit/
│   │   ├── handler.go
│   │   ├── repository.go
│   │   └── sanitizer.go
│   └── webhook/
│       ├── dispatcher.go
│       ├── handler.go
│       └── repository.go
├── migrations/
│   └── [migration_files]
└── tests/
```

**Structure Decision**: Integrated directly into the existing Go `backend/` monorepo. We introduce two new packages: `internal/audit` for storage/API/sanitization, and `internal/webhook` for the SIEM dispatching engine.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| [e.g., 4th project] | [current need] | [why 3 projects insufficient] |
| [e.g., Repository pattern] | [specific problem] | [why direct DB access insufficient] |
