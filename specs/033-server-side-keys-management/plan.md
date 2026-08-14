# Implementation Plan: Server-Side Environment Keys & Tabbed Keys Management

**Branch**: `033-server-side-keys-management` | **Date**: 2026-08-10 | **Spec**: [spec.md](file:///home/tarikelmallah/Projects/FlagManagment/specs/033-server-side-keys-management/spec.md)

**Input**: Feature specification from `/specs/033-server-side-keys-management/spec.md`

## Summary

Implement a dedicated "Keys" settings tab in the Environment view to manage both Client-Side and Server-Side keys cleanly. Introduce a new `environment_server_keys` database entity to allow administrators to create, list, and revoke named server-side keys for backend integrations. A fallback for legacy keys (`api_key_hash` on `environments`) will remain for backward compatibility. 

## Technical Context

**Language/Version**: Go (Backend), React + TypeScript + Vite (Frontend)

**Primary Dependencies**: PostgreSQL, UUIDs, crypto/sha256, TailwindCSS, shadcn/ui

**Storage**: PostgreSQL

**Project Type**: Web Application / Web Service

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **API-First Contract Design**: Passes (API specs defined in `contracts/api.md`).
- **Environment Isolation**: Passes (Server keys strictly scoped to specific environment ID).
- **Governance by Default**: Passes (Multiple keys enable fine-grained access and rotation for distinct services).
- **Local Evaluation Performance**: Passes (No change to SDK performance; auth check maps token to environment quickly).
- **PII Protection & Compliance**: Passes (Server keys are securely hashed using `sha256`, the UI never displays the plain key again after creation, protecting the system against plaintext credential exposure).

## Project Structure

### Documentation (this feature)

```text
specs/033-server-side-keys-management/
├── plan.md              
├── research.md          
├── data-model.md        
├── quickstart.md        
├── contracts/           
└── tasks.md             
```

### Source Code (repository root)
```text
backend/
├── internal/
│   ├── dto/
│   ├── models/
│   ├── repository/
│   ├── services/
│   └── api/

frontend/
├── src/
│   ├── components/
│   │   ├── environments/
│   │   │   ├── EnvironmentSettingsTabs.tsx
│   │   │   ├── ClientSideKeyCard.tsx
│   │   │   ├── ServerSideKeysPanel.tsx
│   │   │   └── CreateServerKeyDialog.tsx
│   ├── pages/
│   │   └── ProjectDetail.tsx
│   └── services/
│       └── apiClient.ts
```

**Structure Decision**: 
The code follows Option 2 (Web application), split into `backend/` and `frontend/`. 
- New DB migration added in `backend/migrations`.
- New UI components added in `frontend/src/components/environments/` to encapsulate the tabbed layout and the client/server key panels.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| *None* | | |
