# FlagManagment Constitution — Detailed Reference Guide

> This document provides an in-depth explanation of each constitutional principle
> for the development team. The canonical constitution lives at
> `.specify/memory/constitution.md`. This file is a companion reference.

---

## Overview

The FlagManagment constitution defines **8 core principles** that govern all feature
development. Every feature specification is validated against these principles
during the `/speckit-plan` constitution check. Violations must be justified or
the plan is blocked.

The principles are derived from two source documents:
- [g-requirements.md](file:///home/tarikelmallah/Projects/FlagManagment/docs/g-requirements.md) — Comprehensive Product Requirements Document
- [p-requirements.md](file:///home/tarikelmallah/Projects/FlagManagment/docs/p-requirements.md) — Product and Technical Requirements Specification

---

## Principle I: API-First Contract Design

### Statement

Every external and internal interface MUST be formally defined before any
business logic implementation begins.

### PRD Sources

| Document | Section | Requirement |
|----------|---------|-------------|
| p-requirements.md | §14.1 Contract-Driven Development | "The API must be strictly defined before business logic implementation" |
| p-requirements.md | §14.1 | "A mocked version of the API must be provided immediately to unblock parallel frontend development" |
| g-requirements.md | §4 (API Protocols row) | "gRPC & REST" with defined boundary responsibilities |

### Rationale

FlagManagment has three distinct consumer surfaces: the React dashboard (REST), the
server-side SDKs (gRPC streaming), and external tooling/Terraform (REST). If
these contracts are not formally defined upfront, integration becomes the
bottleneck. A mocked API immediately unblocks the frontend team from day one.

### Compliance Criteria

- [ ] OpenAPI 3.0 specification exists before any REST endpoint is implemented
- [ ] Protobuf `.proto` files exist before any gRPC handler is implemented
- [ ] A mock server (e.g., Prism, gRPC mock) is available for frontend development
- [ ] API changes go through contract versioning before implementation

### Example Violations

- ❌ Implementing a REST handler without an OpenAPI schema entry
- ❌ Adding a new gRPC message type directly in code without updating `.proto` files
- ❌ Frontend team blocked because backend APIs aren't defined yet

---

## Principle II: Environment Isolation

### Statement

Every environment MUST have cryptographically unique SDK tokens, independent
flag states, and total access isolation.

### PRD Sources

| Document | Section | Requirement |
|----------|---------|-------------|
| g-requirements.md | §5 Multi-Environment Promotion | "Every environment must generate its own unique, cryptographically secure SDK authentication tokens" |
| p-requirements.md | §6.2.1 Environment Isolation | "It must be impossible for an SDK configured for Environment A to access or mutate flags in Environment B" |
| p-requirements.md | §7.4 Security | "Environment tokens must be stored and transmitted securely" |
| g-requirements.md | §5 | "Flag states, targeting rules, and rollout percentages must exist independently within each environment" |

### Rationale

Environment isolation is the foundation of safe multi-environment promotion
pipelines. Without it, a QA engineer could accidentally mutate production flags,
or a staging SDK could evaluate production rules. This is a data integrity and
security concern that cannot be retrofitted.

### Compliance Criteria

- [ ] Each environment generates a unique API key using a CSPRNG
- [ ] API key is hashed (not stored in plaintext) in the `environments` table
- [ ] SDK authentication middleware rejects tokens that don't match the requested environment
- [ ] Flag state queries are always scoped to a single environment ID
- [ ] Promotion operations copy data between environments, never share references

### Example Violations

- ❌ Using the same API key across Dev and Staging environments
- ❌ A database query that returns flag states without an `environment_id` WHERE clause
- ❌ Storing environment tokens in plaintext in the database

---

## Principle III: Governance by Default

### Statement

Granular RBAC, change requests, and immutable audit logging are core capabilities.
They MUST NOT be gated behind paid tiers. No artificial caps on projects or environments.

### PRD Sources

| Document | Section | Requirement |
|----------|---------|-------------|
| g-requirements.md | §1 Product Vision | "unlimited scaling, multi-environment promotion pipelines, and enterprise-grade RBAC completely out of the box" |
| g-requirements.md | §3 Market Differentiation | Explicit gap analysis vs. Unleash (1 project / 2 env caps) and Flagsmith (RBAC paywalled) |
| p-requirements.md | §4.1-4.2 Competitive Constraints | "FlagManagment must not impose such artificial limits" |
| p-requirements.md | §3.2 Commercial SaaS | "self-hosted deployments can run with the same core functionality as the SaaS" |

### Rationale

This is the existential reason FlagManagment exists. If governance features are
limited or paywalled, FlagManagment becomes just another Unleash/Flagsmith clone.
This principle protects the project's competitive identity.

### Compliance Criteria

- [ ] No hardcoded caps on project count, environment count, or flag count
- [ ] RBAC enforcement is active in all deployment modes (local Docker, self-hosted, SaaS)
- [ ] Change request workflows function identically across all deployment modes
- [ ] Audit logging cannot be disabled or feature-gated
- [ ] License checks (if any) do not restrict governance features

### Example Violations

- ❌ An `if tier == "enterprise"` check before enabling change requests
- ❌ Limiting audit log retention to 30 days in the free tier
- ❌ Capping projects at 5 without a technical (performance) justification

---

## Principle IV: Local Evaluation Performance (NON-NEGOTIABLE)

### Statement

SDK flag evaluation MUST complete in under 1 millisecond with zero network calls.

### PRD Sources

| Document | Section | Requirement |
|----------|---------|-------------|
| g-requirements.md | §7 SDKs and Local Evaluation | "guaranteeing sub-millisecond response times without executing any outbound network calls" |
| p-requirements.md | §6.4.3 Server-Side Local Evaluation | "Evaluation must be designed to be sub-millisecond under normal conditions" |
| p-requirements.md | §15.2 Performance Baselines | "must demonstrably execute in under 1 ms during simulated load testing before the feature can be marked as complete" |
| p-requirements.md | §7.1 Performance | "Single flag evaluation should typically complete in under 1 ms on commodity hardware" |

### Rationale

Feature flags sit in the hot path of every application request. A 10ms network
call to evaluate a flag at 10,000 RPS adds 100 seconds of cumulative latency per
second. Local evaluation is the only architecture that scales to enterprise
traffic volumes.

### Compliance Criteria

- [ ] Benchmark test proves <1ms p99 evaluation time under simulated load
- [ ] Evaluation function makes zero network calls (pure in-memory operation)
- [ ] SDK bootstraps by downloading a complete snapshot, not querying per-flag
- [ ] Delta updates arrive via push (gRPC stream), not poll
- [ ] SDK continues evaluating correctly when backend connection is lost
- [ ] Load test results are included in the feature's quickstart.md

### Example Violations

- ❌ An SDK that makes an HTTP call to evaluate a flag
- ❌ A caching layer with a 5-second TTL that causes periodic network fetches
- ❌ Marking an SDK feature "complete" without sub-millisecond load test proof

---

## Principle V: Test-First Quality Gates

### Statement

Mandatory coverage minimums (80% backend, 70% frontend) and comprehensive test
domains across all platform capabilities.

### PRD Sources

| Document | Section | Requirement |
|----------|---------|-------------|
| p-requirements.md | §15.1 Code Quality & Testing Minimums | "80% unit test coverage for the backend engine, 70% for the frontend UI" |
| p-requirements.md | §15.1 | "Code must pass standard static analysis and linting pipelines with zero critical warnings" |
| p-requirements.md | §9.2 Testing and Validation | Full list of mandatory test domains |
| p-requirements.md | §9.2 | End-to-end test scenarios: feature rollout, change request cycle, automated rollback |

### Rationale

Feature flag platforms are infrastructure software — bugs here cascade into
every application that depends on them. Coverage thresholds and domain-specific
test requirements ensure reliability for enterprise customers who deploy
FlagManagment in production-critical paths.

### Compliance Criteria

- [ ] Go backend maintains ≥80% test coverage (measured by `go test -cover`)
- [ ] React frontend maintains ≥70% test coverage (measured by Jest/Vitest)
- [ ] CI pipeline fails on coverage regression below thresholds
- [ ] Static analysis (golangci-lint, ESLint) passes with zero critical warnings
- [ ] End-to-end tests exist for: Dev→QA→Staging→Prod rollout, change request cycle, automated rollback
- [ ] SDK evaluation tests verify cross-language determinism (same input → same output)

### Example Violations

- ❌ Merging a PR that drops backend coverage to 75%
- ❌ Skipping integration tests for the change request workflow
- ❌ No end-to-end test for the promotion pipeline

---

## Principle VI: OpenFeature Interoperability

### Statement

SDKs conform to the CNCF OpenFeature API standard. Deviations are documented.

### PRD Sources

| Document | Section | Requirement |
|----------|---------|-------------|
| g-requirements.md | §7 SDKs | "All SDKs must conform strictly to the CNCF OpenFeature API standard" |
| p-requirements.md | §6.4.2 OpenFeature Compatibility | "SDK APIs must conform to the CNCF OpenFeature standard where feasible" |
| p-requirements.md | §6.4.2 | "Where OpenFeature semantics conflict with FlagManagment-specific features, document any deviations clearly" |

### Rationale

OpenFeature is the industry standard for feature flag APIs. Conformance ensures
users can adopt FlagManagment without rewriting their evaluation code, and can migrate
away if needed — which paradoxically increases trust and adoption.

### Compliance Criteria

- [ ] Each SDK implements the OpenFeature Provider interface
- [ ] Standard evaluation methods (getBooleanValue, getStringValue, etc.) are present
- [ ] FlagManagment-specific extensions (sequential dependencies, analytics hooks) are clearly marked
- [ ] Deviation documentation exists in SDK README and feature spec
- [ ] Interoperability test: swap an OpenFeature provider to verify API compatibility

### Example Violations

- ❌ An SDK that uses a completely proprietary API surface with no OpenFeature provider
- ❌ Breaking the OpenFeature `evaluate` method signature for convenience
- ❌ Undocumented deviations from the OpenFeature specification

---

## Principle VII: PII Protection & Compliance

### Statement

All PII used for identity bucketing is salted and hashed before storage. Audit
logs are sanitized. SOC 2-style compliance readiness is mandatory.

### PRD Sources

| Document | Section | Requirement |
|----------|---------|-------------|
| p-requirements.md | §15.3 Data Privacy & SOC2 Readiness | "must natively salt and hash any PII utilized for identity bucketing prior to database storage" |
| p-requirements.md | §15.3 | "Audit logs must be actively sanitized to ensure no plaintext API keys or sensitive targeting metadata are inadvertently captured" |
| p-requirements.md | §7.5 Compliance and Auditability | "Design audit logs and RBAC to support common compliance requirements (e.g., SOC 2-style change management)" |
| p-requirements.md | §7.5 | "Provide configuration options for log retention policies" |

### Rationale

Feature flag platforms receive user identity attributes (email, tenant ID, user
ID) for targeting. Storing this data in plaintext creates regulatory exposure
under GDPR, CCPA, and similar frameworks. Audit log sanitization prevents
accidental credential leakage.

### Compliance Criteria

- [ ] Identity attributes are salted + hashed (e.g., bcrypt, Argon2) before database INSERT
- [ ] Audit log serialization strips API keys, tokens, and PII from JSON state snapshots
- [ ] Log retention is configurable (default: 90 days, adjustable per compliance needs)
- [ ] No plaintext email addresses, user IDs, or targeting attributes in `audit_logs` table
- [ ] Security review checklist item on every feature spec

### Example Violations

- ❌ Storing `user_email` in plaintext in the `evaluation_context` column
- ❌ Audit log entry that includes `"api_key": "sk-live-abc123..."` in the JSON state
- ❌ Hardcoded 365-day log retention with no configuration option

---

## Principle VIII: Cloud-Native Portability

### Statement

Docker/K8s deployment, multi-arch builds (x86 + ARM + RPi4), environment-driven
configuration, single-command local bootstrap.

### PRD Sources

| Document | Section | Requirement |
|----------|---------|-------------|
| p-requirements.md | §5.2 Deployment and Operations | "Docker images for all services", "Kubernetes manifests or Helm charts" |
| p-requirements.md | §11.2 Multi-Architecture Support | "must be optimized to build and run seamlessly across varying hardware profiles" |
| p-requirements.md | §11.2 | "range from standard x86 and ARM-based Mac environments down to...Raspberry Pi 4" |
| p-requirements.md | §11.1 Local Workspace | "Developers should be able to bootstrap the entire stack locally with a single command" |
| p-requirements.md | §5.2 Configuration | "Environment-driven configuration for DB connections, cache, logging" |

### Rationale

Multi-architecture support drives both open-source contributor adoption (ARM
Macs are now dominant in the developer ecosystem) and enterprise edge-testing
scenarios. A single `docker compose up` command eliminates onboarding friction
and makes the project attractive to casual contributors.

### Compliance Criteria

- [ ] Dockerfiles use multi-stage builds with `--platform` flags for x86 and ARM
- [ ] `docker compose up` bootstraps the entire stack (backend, frontend, PostgreSQL, Redis)
- [ ] All configuration is via environment variables (no hardcoded connection strings)
- [ ] Kubernetes manifests or Helm charts are provided in the repository
- [ ] ARM build is tested on Apple Silicon and Raspberry Pi 4 in CI
- [ ] IDE workspace configs (.vscode/, .windsurf/) are committed to the repository

### Example Violations

- ❌ A Dockerfile that only builds for `linux/amd64`
- ❌ Database connection string hardcoded in a Go source file
- ❌ Local setup requires 15 manual steps instead of `docker compose up`

---

## Quick Reference Table

| # | Principle | Constitution Section | Key Metric |
|---|-----------|---------------------|------------|
| I | API-First Contract Design | Core Principles | OpenAPI + Protobuf before code |
| II | Environment Isolation | Core Principles | Unique crypto tokens per env |
| III | Governance by Default | Core Principles | Zero paywalled features |
| IV | Local Evaluation Performance | Core Principles | <1ms p99 evaluation |
| V | Test-First Quality Gates | Core Principles | 80% backend / 70% frontend |
| VI | OpenFeature Interoperability | Core Principles | CNCF standard conformance |
| VII | PII Protection & Compliance | Core Principles | Salt+hash all PII |
| VIII | Cloud-Native Portability | Core Principles | Multi-arch Docker + single-command setup |
