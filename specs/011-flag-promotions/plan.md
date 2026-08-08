# Implementation Plan: One-Click Flag Environment Promotions

**Branch**: `011-flag-promotions` | **Date**: 2026-07-25 | **Spec**: [spec.md](file:///home/tarikelmallah/Projects/FlagManagment/specs/011-flag-promotions/spec.md)

## Summary

Implement a Promotion API endpoint and UI action to copy feature flag rulesets between environments, integrating seamlessly with protected environments via Change Requests.

## Technical Context

**Language/Version**: Go 1.22, React/TypeScript
**Primary Dependencies**: PostgreSQL, Chi Router
**Storage**: PostgreSQL (`environment_flag_states`, `change_requests`)
**Testing**: Go testing, React Testing Library
**Constraints**: Protected target environments MUST trigger Change Requests instead of direct mutations.

## Constitution Check

- **API-First Contract Design**: Passes. `POST /api/v1/projects/{projectId}/flags/{flagId}/promote` endpoint.
- **Environment Isolation**: Passes. Operates across explicit source and target environment IDs.
- **Governance by Default**: Passes. Protected environments generate Change Requests; actions are logged in Audit Logs.
- **Local Evaluation Performance (NON-NEGOTIABLE)**: Passes. No impact on SDK local evaluation path.

## Project Structure

```text
specs/011-flag-promotions/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
└── checklists/
    └── requirements.md
```

### Source Code

```text
backend/
├── internal/
│   ├── api/
│   │   └── promotion.go           # Promotion HTTP handler
│   └── services/
│       └── promotion_service.go   # Logic for copying flag state or creating Change Request

frontend/
├── src/
│   ├── components/
│   │   └── PromoteFlagModal.tsx   # UI Modal to select target environment and execute promotion
│   └── services/
│       └── promotionApi.ts
```
