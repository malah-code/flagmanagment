# FlagManagment — Test Plan Index

> **Canonical test plan location**: [`/specs/038-comprehensive-test-plan/spec.md`](../specs/038-comprehensive-test-plan/spec.md)

This document serves as the **entry point** to the comprehensive test plan for the FlagManagment platform. The full test plan contains **158 test cases** covering all 37 feature specs, `g-requirements.md`, and `p-requirements.md`.

---

## How to Use This Test Plan

### For Human Testers

1. Open [`/specs/038-comprehensive-test-plan/spec.md`](../specs/038-comprehensive-test-plan/spec.md)
2. Navigate to the relevant area (Parts 1–24)
3. Execute each test case sequentially using its Given/When/Then steps
4. Mark test results (Pass/Fail) in your test tracking tool with the TC ID (e.g., `TC-FLAG-001`)

### For AI Testing Agents

1. Load the full spec: `/specs/038-comprehensive-test-plan/spec.md`
2. Parse each test case using the structured format:
   - **ID**: `TC-{AREA}-{NNN}` (e.g., `TC-AUTH-001`)
   - **Layer**: The testing layer(s) to invoke (API, SDK, UI, etc.)
   - **Spec**: Source feature spec(s) to cross-reference
   - **Given/When/Then**: Precondition, action, and assertion

### Test Execution Prerequisites

```bash
# Start the full stack
docker-compose up -d

# Or run locally
cd backend && go run ./cmd/server  # Backend API + gRPC
cd frontend && npm run dev         # Frontend dashboard
```

---

## Test Areas Summary

| Part | Area | Count | Layers |
|---|---|---|---|
| 1 | Authentication & User Management | 10 | UI, API |
| 2 | Project Management | 5 | API, UI |
| 3 | Environment Management | 10 | API, UI, SDK |
| 4 | Feature Flags (Boolean & Multivariate) | 16 | API, UI, SDK |
| 5 | Contextual Targeting Engine | 10 | API, UI, SDK |
| 6 | Percentage Rollouts | 5 | API, SDK |
| 7 | Sequential Flag Dependencies | 5 | API, SDK |
| 8 | Scheduled Flag Changes | 4 | API, UI |
| 9 | Flag Promotion Pipeline | 5 | API, UI |
| 10 | Change Requests & Approval Workflow | 7 | API, UI |
| 11 | RBAC & Permissions | 7 | API |
| 12 | Immutable Audit Logs | 8 | API, UI |
| 13 | SDK: Server-Side Local Evaluation | 10 | SDK |
| 14 | SDK Event Forwarding (A/B Analytics) | 3 | SDK |
| 15 | Stale Flag Detection | 4 | API, UI |
| 16 | Telemetry & Automated Kill-Switches | 4 | API, UI |
| 17 | Admin Dashboard & UX | 7 | UI |
| 18 | PII Hashing & Data Privacy | 3 | API |
| 19 | Slack & Webhook Notifications | 4 | API, UI |
| 20 | Edge Proxy / Relay Node | 3 | Infrastructure |
| 21 | Terraform Provider | 4 | IaC |
| 22 | CI/CD Pipeline & DevOps | 5 | CI/CD |
| 23 | Non-Functional & Performance | 7 | API, SDK, Infrastructure |
| 24 | End-to-End Scenarios | 6 | All |
| **Total** | | **158** | |

---

## Test ID Naming Convention

```
TC-{AREA}-{NNN}

AREA codes:
  AUTH     Authentication & User Management
  PROJ     Project Management
  ENV      Environment Management
  FLAG     Feature Flags
  TARGET   Contextual Targeting
  ROLLOUT  Percentage Rollouts
  DEP      Sequential Dependencies
  SCHED    Scheduled Flag Changes
  PROMO    Flag Promotion Pipeline
  CR       Change Requests
  RBAC     RBAC & Permissions
  AUDIT    Audit Logs
  SDK      SDK Evaluation
  EVENT    SDK Event Forwarding
  STALE    Stale Flag Detection
  TELEM    Telemetry & Kill-Switches
  UI       Admin Dashboard UX
  PII      Data Privacy
  SLACK    Slack/Webhook Notifications
  EDGE     Edge Proxy
  TERRAFORM Terraform Provider
  CICD     CI/CD Pipeline
  PERF     Performance
  SEC      Security
  E2E      End-to-End Scenarios
```

---

## Source Documents

- [`docs/g-requirements.md`](./g-requirements.md) — General system requirements
- [`docs/p-requirements.md`](./p-requirements.md) — Product & technical requirements
- [`specs/`](../specs/) — 37 feature specifications (001–037)
- [`specs/038-comprehensive-test-plan/spec.md`](../specs/038-comprehensive-test-plan/spec.md) — **This test plan** (full test cases)
- [`specs/038-comprehensive-test-plan/checklists/requirements.md`](../specs/038-comprehensive-test-plan/checklists/requirements.md) — Quality checklist

---

*Last updated: 2026-08-19 | Generated from all 37 feature specs + g-requirements.md + p-requirements.md*
