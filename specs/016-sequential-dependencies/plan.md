# Implementation Plan: Sequential Dependencies

**Branch**: `016-sequential-dependencies` | **Date**: 2026-08-07 | **Spec**: [spec.md](file:///home/tarikelmallah/Projects/FlagManagment/specs/016-sequential-dependencies/spec.md)

**Input**: Feature specification from `/specs/016-sequential-dependencies/spec.md`

## Summary

This feature implements Sequential Dependencies, allowing a feature flag to depend on the state of a parent flag. The technical approach involves adding a self-referencing `parent_flag_id` to the `feature_flags` table, implementing directed graph cycle detection at the API layer to prevent infinite loops, and updating the server-side SDK evaluation engine to lazily evaluate parent dependencies and short-circuit to safe fallback states when necessary.

## Technical Context

**Language/Version**: Go 1.21 (Backend) / TypeScript (Frontend)

**Primary Dependencies**: PostgreSQL, standard `database/sql` driver, React/Vite.

**Storage**: PostgreSQL (`feature_flags` table).

**Testing**: `go test` for backend, including extensive API-level tests for cycle detection.

**Target Platform**: Linux server, modern web browsers.

**Project Type**: Web service + Frontend dashboard.

**Performance Goals**: < 1ms flag evaluation locally inside the SDK engine.

**Constraints**: Must strictly prevent circular dependencies (A -> B -> A). Dependency chain evaluation must maintain the sub-millisecond evaluation latency constraint.

**Scale/Scope**: Impacts core feature flag creation/update API and SDK evaluation engine.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- [X] API-First Contract Design: N/A (minor DB addition, no new external API contracts needed beyond modifying the Flag creation payload).
- [X] Local Evaluation Performance: Checked. Using lazy evaluation inside SDKs keeps it well under 1ms.
- [X] Test-First Quality Gates: Will add test coverage for DFS cycle detection algorithm.

## Project Structure

### Documentation (this feature)

```text
specs/016-sequential-dependencies/
├── plan.md              
├── research.md          
├── data-model.md        
├── quickstart.md        
└── tasks.md             
```

### Source Code (repository root)

```text
backend/
├── migrations/
│   └── 000014_sequential_dependencies.up.sql
├── internal/
│   ├── models/
│   │   └── feature_flag.go
│   ├── repository/
│   │   └── flag_repo.go
│   ├── services/
│   │   ├── flag_service.go
│   │   └── cycle_detector.go
│   └── sdk/
│       └── evaluator.go

frontend/
├── src/
│   ├── components/
│   │   └── flags/
│   │       ├── CreateFlagDialog.tsx
│   │       └── FlagDetails.tsx
│   └── types/
│       └── index.ts
```

**Structure Decision**: Utilizing the existing backend/frontend structure. Adding dedicated cycle detection logic in `backend/internal/services` to isolate the graph traversal logic.

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| In-memory cycle detection on backend | To prevent DB recursion and maintain cross-DB compatibility | Database recursive CTEs are database-specific and hard to test. |
