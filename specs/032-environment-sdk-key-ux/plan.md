# Implementation Plan: Environment SDK Key UX & Integration Guide

**Branch**: `032-environment-sdk-key-ux` | **Date**: 2026-08-10 | **Spec**: [spec.md](file:///home/tarikelmallah/Projects/FlagManagment/specs/032-environment-sdk-key-ux/spec.md)

**Input**: Feature specification from `/specs/032-environment-sdk-key-ux/spec.md`

## Summary

Enhance Environment Settings by adding an always-visible, copyable Client SDK Key for each environment, clearly badging Public Client Keys vs. Private Admin Keys, and creating an interactive Integration Guide modal with ready-to-use React, Node.js, Python, and Go SDK code snippets.

## Technical Context

**Language/Version**: React + TypeScript (Frontend), Go (Backend DTO/API update to expose client key if necessary)

**Primary Dependencies**: React, TailwindCSS, Lucide Icons, react-hot-toast

**Storage**: PostgreSQL (Backend)

**Testing**: Vitest / React Testing Library

**Target Platform**: Web Browser

**Project Type**: Web Application Dashboard

**Performance Goals**: Instant client key copy, instantaneous tab switching in code guide modal.

**Constraints**: High contrast for code blocks (dark theme code blocks), responsive modal layout.

**Scale/Scope**: `EnvironmentsList.tsx`, `EnvironmentSidebar.tsx`, and new `SDKIntegrationModal.tsx`.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **API-First Contract Design**: Public client keys must be returned safely in Environment DTOs for authenticated members.
- **Environment Isolation**: Each environment's key is isolated and maps to its unique ID.
- **Technology Stack Constraints**: React + TypeScript + Tailwind.

All gates pass.

## Project Structure

### Documentation (this feature)

```text
specs/032-environment-sdk-key-ux/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
└── quickstart.md        # Phase 1 output
```

### Source Code (repository root)

```text
frontend/
├── src/
│   ├── components/
│   │   └── environments/
│   │       ├── SDKIntegrationModal.tsx   # New integration guide modal
│   │       └── EnvironmentsList.tsx      # Updated environment card layout
```

**Structure Decision**: Web application (frontend component addition)

## Complexity Tracking

None.
