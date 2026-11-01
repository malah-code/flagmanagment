# Tasks: Flag Creation & Targeting Workflows

**Input**: Design documents from `/specs/035-flag-creation-targeting-workflows/`
**Prerequisites**: [plan.md](plan.md), [spec.md](spec.md)

---

## Phase 1: Implementation

- [x] T001 [US1] Build multi-step flag creation modal in `frontend/src/components/flags/CreateFlagDialog.tsx`
- [x] T002 [US1] Support Boolean, Multivariate, and JSON Remote Config selection with live previews
- [x] T003 [US1] Provide inline contextual rule targeting builder with operators (EQUALS, CONTAINS, REGEX)
- [x] T004 [US1] Integrate mutation call to `POST /api/v1/projects/{projectId}/flags` and invalidate TanStack cache

