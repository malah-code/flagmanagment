<!-- Sync Impact Report
  Version change: 0.0.0 (template) → 1.0.0
  Modified principles: N/A (initial population from template)
  Added sections:
    - I. API-First Contract Design
    - II. Environment Isolation
    - III. Governance by Default
    - IV. Local Evaluation Performance (NON-NEGOTIABLE)
    - V. Test-First Quality Gates
    - VI. OpenFeature Interoperability
    - VII. PII Protection & Compliance
    - VIII. Cloud-Native Portability
    - Technology Stack Constraints
    - Development Workflow & Quality Gates
  Removed sections: None
  Templates requiring updates:
    - .specify/templates/plan-template.md ✅ (no changes needed — gates are dynamic)
    - .specify/templates/spec-template.md ✅ (no changes needed)
    - .specify/templates/tasks-template.md ✅ (no changes needed)
  Follow-up TODOs: None
-->
# FlagManagment Constitution

## Core Principles

### I. API-First Contract Design

Every external and internal interface MUST be formally defined before any
business logic implementation begins. REST APIs require an OpenAPI 3.0
specification. gRPC interfaces require Protobuf definitions. A mocked
version of each API MUST be provided immediately to unblock parallel
frontend and SDK development.

**Rationale:** Contract-driven development eliminates integration
surprises and enables true parallel workstreams between backend,
frontend, and SDK teams.

### II. Environment Isolation

Every environment MUST generate its own unique, cryptographically secure
SDK authentication token. Flag states, targeting rules, and rollout
percentages MUST exist independently within each environment. It MUST be
impossible for an SDK configured for Environment A to access or mutate
flags in Environment B through normal usage. Environment tokens MUST be
stored and transmitted securely; server-side tokens MUST NOT be exposed
in client-side code.

**Rationale:** Strict isolation prevents accidental cross-environment
contamination and is a prerequisite for safe multi-environment promotion
pipelines (Dev → QA → Staging → Production).

### III. Governance by Default

Granular RBAC, change request workflows, and immutable audit logging
MUST be core capabilities available in all deployments — self-hosted and
SaaS alike. These features MUST NOT be gated behind paid tiers or
artificial feature caps. The platform MUST support unlimited projects and
unlimited environments with no licensing-driven limits.

**Rationale:** This is FlagManagment's primary market differentiator. Gating
governance features behind paywalls is the exact problem FlagManagment exists
to solve.

### IV. Local Evaluation Performance (NON-NEGOTIABLE)

Server-side SDK flag evaluation MUST execute in under 1 millisecond on
commodity hardware. All evaluations MUST occur in-memory within the SDK
with zero outbound network calls during evaluation. SDKs MUST bootstrap
by downloading a complete ruleset snapshot, receive lightweight delta
updates over persistent streaming connections, and continue evaluating
accurately using the last known good state if the connection drops.

**Rationale:** Network latency in the hot path is unacceptable for
enterprise-scale applications serving millions of requests. This
constraint MUST be verified via load testing before any SDK feature is
marked complete.

### V. Test-First Quality Gates

Automated test coverage MUST meet or exceed: 80% unit test coverage for
the Go backend engine, 70% coverage for the React frontend. Code MUST
pass standard static analysis and linting pipelines with zero critical
warnings before merge. Automated tests MUST cover: flag evaluation logic
(including multivariate and dependencies), promotion pipelines, RBAC
enforcement, protected environment behavior, audit logging, change
request workflows, SDK local evaluation, delta synchronization, and
observability triggers.

**Rationale:** High coverage thresholds and comprehensive test domains
ensure platform reliability for enterprise customers who depend on
feature flags in production-critical paths.

### VI. OpenFeature Interoperability

All SDKs MUST conform to the CNCF OpenFeature API standard where
feasible to allow easy vendor migration and interoperability. Where
OpenFeature semantics conflict with FlagManagment-specific features (e.g.,
sequential dependencies, analytics forwarding), deviations MUST be
documented clearly in both SDK documentation and the feature's spec.

**Rationale:** OpenFeature compliance reduces adoption friction and gives
users confidence that they are not locked into a proprietary API surface.

### VII. PII Protection & Compliance

The system MUST natively salt and hash any Personal Identifiable
Information used for identity bucketing prior to database storage. Audit
logs MUST be actively sanitized to ensure no plaintext API keys or
sensitive targeting metadata are inadvertently captured. The system MUST
be designed to support common compliance requirements (SOC 2-style change
management). Configurable log retention policies MUST be provided.

**Rationale:** Feature flag platforms inherently process user identity
data for targeting. Failing to protect PII exposes the platform and its
users to regulatory and legal risk.

### VIII. Cloud-Native Portability

All services MUST be containerized as Docker images with reference
Kubernetes manifests or Helm charts. Local containerized orchestration
(Docker Compose) MUST build and run seamlessly on x86, ARM-based Mac
environments, and Raspberry Pi 4. Configuration MUST be environment-driven
(environment variables for DB connections, cache, logging). Developers
MUST be able to bootstrap the entire stack locally with a single command.

**Rationale:** Multi-architecture support and frictionless local setup
directly drive open-source contributor adoption and edge-testing
flexibility for enterprise users.

## Technology Stack Constraints

The following technology choices are mandatory unless a formal amendment
to this constitution is approved:

- **Backend Engine:** Go (Golang) — goroutine concurrency for high-throughput
  SDK serving
- **Frontend Dashboard:** React + TypeScript + Vite — standardized component
  library (Shadcn/UI or MUI) for professional aesthetics
- **Primary Datastore:** PostgreSQL 16+ — ACID compliance, relational schemas,
  JSONB for targeting rules and remote config payloads
- **Caching & Pub/Sub:** Redis 7+ — in-memory flag caching, real-time change
  broadcasting to SDK connections
- **API Protocols:** REST/JSON over HTTPS (external), gRPC/Protobuf
  (internal and SDK streaming)
- **Deployment:** Docker multi-arch images, Kubernetes manifests, Helm charts
- **Database Migrations:** golang-migrate or equivalent versioned migration tool
- **Hashing Algorithm:** MurmurHash3 for deterministic cross-language SDK bucketing

Alternative technologies MAY be proposed only if they meet or exceed the
functional and non-functional requirements defined in the product
requirements documents and maintain cloud-native, open-source-friendly
characteristics.

## Development Workflow & Quality Gates

### Spec-Driven Development

All feature work MUST follow the Spec Kit lifecycle:

1. `/speckit-specify` — Define the feature (WHAT and WHY, not HOW)
2. `/speckit-plan` — Generate implementation plan with constitution check
3. `/speckit-tasks` — Break into dependency-ordered, actionable tasks
4. `/speckit-implement` — Execute tasks

No implementation code may be written until the corresponding feature
spec passes the quality checklist and the plan's constitution check
passes with no unjustified violations.

### Phased Delivery

Features MUST be delivered in the phased order defined in the execution
roadmap (PLAN.md). Each phase's features MUST have their dependencies
satisfied before implementation begins. The three delivery phases are:

1. **Core Engine MVP** — Boolean/Multivariate flags, unlimited environments,
   PostgreSQL schema, gRPC/REST APIs, React dashboard, initial SDKs
2. **Governance & RBAC** — Granular RBAC, change requests, visual diffing,
   mandatory approvals, immutable audit log
3. **Enterprise & Observability** — Telemetry triggers, kill-switch engine,
   full OpenFeature compliance, remaining SDKs, Terraform provider

### Code Review & Approval

- Feature specs require review before proceeding to planning
- All code changes require at least one peer review before merge
- Architectural deviations MUST be retroactively documented in the feature
  spec and validated against this constitution

## Governance

This constitution is the supreme governance document for the FlagManagment
project. It supersedes all other development practices, conventions, and
ad-hoc agreements. Compliance with these principles is mandatory for all
contributors.

### Source Control & Repository
- **Official Repository**: [https://github.com/malah-code/flagmanagment](https://github.com/malah-code/flagmanagment)
- All feature branches, specifications, and pull requests MUST target this central repository.

### Amendment Procedure

1. Propose the amendment in writing with rationale
2. Obtain review from at least one Lead Architect and one Lead Engineer
3. Update this document with the change and bump the version per semantic
   versioning:
   - **MAJOR:** Backward-incompatible governance/principle removals or
     redefinitions
   - **MINOR:** New principle/section added or materially expanded guidance
   - **PATCH:** Clarifications, wording, typo fixes, non-semantic refinements
4. Update dependent templates if the amendment affects spec, plan, or task
   generation rules

### Compliance Review

All pull requests and code reviews MUST verify compliance with these
principles. Any complexity that contradicts a principle MUST be justified
in the feature's plan.md Complexity Tracking table. Unjustified violations
are grounds for blocking merge.

**Version**: 1.0.0 | **Ratified**: 2026-07-18 | **Last Amended**: 2026-07-18
