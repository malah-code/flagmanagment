# FlagManagment Product and Technical Requirements Specification

## 1. Document Overview

### 1.1 Purpose
This document defines the product vision, functional requirements, non-functional requirements, suggested technical architecture, and delivery expectations for the FlagManagment feature flag and remote configuration platform.
It is intended to be a standalone specification that a third-party development company can use to design, implement, and deliver the platform end-to-end.

### 1.2 Audience
- Software architects and technical leads
- Backend and frontend developers
- DevOps/SRE engineers
- QA and test automation engineers
- Product managers and business analysts

### 1.3 Guiding Principles
- Functional and non-functional behavior described in this document is mandatory unless explicitly labeled as “Recommended” or “Optional”.
- Technology stack choices in this document are **recommended defaults**; the vendor may propose alternatives that meet or exceed the functional and non-functional requirements and align with cloud-native, open-source-friendly tooling.
- The platform must be designed so that organizations can self-host and operate it independently, without relying on the vendor’s SaaS offering.

## 2. Product Vision

### 2.1 Vision Statement
FlagManagment is a cloud-native, open-source feature flag and remote configuration platform designed to bridge the gap between basic open-source toggle libraries and prohibitively expensive enterprise suites.
Current market leaders intentionally constrain their open-source offerings by gating multi-environment workflows, granular governance, and compliance controls behind enterprise tiers; FlagManagment rejects this model.
FlagManagment will provide unlimited projects, unlimited environments, enterprise-grade RBAC, change requests, and immutable audit logs out of the box for self-hosted deployments, while still enabling a managed SaaS offering as a separate commercial product.

### 2.2 Core Objectives
- Enable engineering, QA, and product teams to safely decouple deployment from release via feature flags and remote configuration.
- Provide rich multi-environment promotion pipelines (Dev → QA → Staging → Production) with strict isolation and governance.
- Deliver enterprise-grade governance, RBAC, and auditability as core capabilities, not paywalled add-ons.
- Support high-throughput, low-latency local flag evaluation via server-side SDKs with streaming updates.
- Integrate with modern observability stacks to support automated kill-switches and rollbacks.

## 3. Licensing and Business Model

### 3.1 License Model
- Use a **Business Source License (BSL)** or comparable “fair-source” license.
- Permit organizations to self-host, modify, and utilize the entire platform internally for their own business operations at zero licensing cost.
- Explicitly prohibit third-party managed service providers from wrapping the codebase and reselling it as a competing public SaaS.

### 3.2 Commercial SaaS Offering (Informational)
- The founding organization will offer a fully managed, multi-tenant SaaS hosting option.
- SaaS monetization should be based on managed hosting, high-throughput API SLAs, and premium integrations, **not** on withholding core security or governance features.
- The codebase must be structured so that self-hosted deployments can run with the same core functionality as the SaaS, differing only in operational and commercial aspects.

## 4. Market Differentiation and Competitive Constraints

This section exists to clarify design intent; it is not an implementation of competitor features but a set of constraints FlagManagment must satisfy.

### 4.1 Unleash-Specific Constraints
- The open-source version of Unleash caps usage at **one project and two environments**; FlagManagment **must not** impose such artificial limits.
- FlagManagment must support:
  - Unlimited projects
  - Unlimited environments per project
  - No hard-coded caps based on licensing; limits, if any, should be technical (e.g., performance) rather than licensing-driven

### 4.2 Flagsmith-Specific Constraints
- Flagsmith restricts granular RBAC, environment-level permissions, and approval workflows to higher-priced tiers.
- FlagManagment must provide:
  - Granular RBAC at global, project, and environment levels
  - Environment-level permissions and separation of duties
  - Change requests and approval workflows for protected environments
  - Immutable audit logs with export and streaming capabilities

These features must be available in all deployments (self-hosted and SaaS) as core capabilities.

## 5. High-Level Architecture and Technology Stack

### 5.1 Architectural Goals
- Cloud-native design suitable for containerized deployment (Docker/Kubernetes).
- Horizontal scalability to handle millions of feature evaluation requests per second via local evaluation in SDKs.
- Clear separation between:
  - Core flag engine and data store
  - Admin/dashboard UI
  - Public API layer
  - SDKs and client libraries

### 5.2 Recommended Technology Stack

The vendor may propose alternatives, but any changes must preserve:
- Cloud-native deployment compatibility
- Strong community ecosystem and long-term maintainability
- Support for relational integrity and JSON-based targeting rules

**Backend Engine (Recommended)**  
- Language: **Go (Golang)**
- Rationale: Concurrency (goroutines), strong performance for backend services, wide adoption in cloud-native ecosystems

**Frontend Dashboard (Recommended)**  
- Framework: **React**
- Language: **TypeScript**
- Tooling: **Vite** or a similarly modern bundler
- Rationale: Strong ecosystem, TypeScript safety for complex UI logic, fast builds

**Primary Datastore (Required Capabilities)**  
- Recommended: **PostgreSQL**
- Requirements:
  - Full ACID compliance
  - Support for relational schemas (projects, environments, flags, RBAC tables)
  - JSON/JSONB columns for complex targeting rules and remote configuration payloads

**Caching and Pub/Sub (Required Capabilities)**  
- Recommended: **Redis**
- Requirements:
  - In-memory caching for flag evaluations and configuration snapshots
  - Pub/Sub or streaming features to broadcast changes to SDK connections

**API Protocols (Required Capabilities)**  
- External APIs: REST/JSON over HTTPS
- Internal and SDK streaming: gRPC/Protobuf or equivalent bidirectional streaming protocol

**Deployment and Operations**  
- Containerization: Docker images for all services
- Orchestration: Kubernetes manifests or Helm charts for reference deployments
- Configuration: Environment-driven configuration for DB connections, cache, logging, etc.

## 6. Functional Requirements

### 6.1 Flag Types and Evaluation Rules

#### 6.1.1 Boolean Toggles
- Support standard on/off flags to gate code paths.
- Evaluation returns a boolean value for each identity and environment.

#### 6.1.2 Multivariate Flags
- Support A/B/n flags with named variants.
- Allocate variants via deterministic percentage-based bucketing derived from identity hashes (e.g., user ID, session ID).
- Ensure stable assignment per identity within a given environment.

#### 6.1.3 Remote Configuration
- Flags must be able to return strongly typed payload values along with the flag state:
  - String
  - Number
  - Boolean
  - JSON object
- Payloads should allow runtime configuration changes without code redeploys.

#### 6.1.4 Sequential Dependencies
- Allow flags to depend on parent flags.
- If a parent flag is disabled, all dependent flags must evaluate to a safe fallback state.
- The engine must detect and prevent circular dependencies at creation time.

#### 6.1.5 Contextual Targeting
- Evaluation must support targeting operators against identity/context attributes:
  - equals
  - not equals
  - contains / substring
  - regex matching
  - array inclusion/exclusion
- Attributes may include tenant ID, region, email domain, plan tier, etc.
- Targeting rules must be stored in structured form (e.g., JSON) to enable visual rule builders in the UI.

### 6.2 Multi-Environment Promotion Pipelines

#### 6.2.1 Environment Isolation
- Each environment must have a unique, cryptographically secure SDK authentication token.
- Flag states, targeting rules, and rollout percentages must be independent by environment.
- It must be impossible for an SDK configured for Environment A to access or mutate flags in Environment B via normal usage.

#### 6.2.2 Environment Hierarchy and Promotion
- Support typical hierarchies: Local, Dev, QA, Staging, Production, plus additional custom environments.
- Provide a “Promote” operation in UI and API:
  - Copies a flag’s configuration from a source environment (e.g., QA) to a target environment (e.g., Staging)
  - Includes rollout percentages, targeting rules, remote config payloads, and dependencies
  - Allows selective promotion (e.g., only specific rules or fields) when needed

#### 6.2.3 Promotion Safety
- Promotions to protected environments (e.g., Production) must be subject to change request and approval workflows (see Governance).
- Promotion operations must be fully logged in the audit system.

### 6.3 Governance, RBAC, and Compliance

#### 6.3.1 RBAC Model
- Roles and permissions must be calculable at:
  - Global level (instance-wide roles)
  - Project level
  - Environment level
- At minimum, support the following role types (names can vary):
  - System Administrator
  - Project Owner
  - Release Manager
  - QA Engineer
  - Read-Only Auditor
- Example permission behavior:
  - QA Engineer: full read/write access in Dev/QA, read-only in Production
  - Release Manager: can approve change requests in protected environments

#### 6.3.2 Protected Environments and Change Requests
- Environments can be marked as **Protected** (e.g., Production).
- Any mutation to flags or rules in a protected environment must:
  - Create a Change Request object rather than applying the change immediately
  - Present a diff view (git-style) showing current rule configuration vs proposed configuration
- Change Requests must support statuses:
  - Pending
  - Approved
  - Rejected
- Only users with appropriate roles (e.g., Release Manager) can approve or reject Change Requests.
- Upon approval, changes are applied atomically; upon rejection, no changes are applied.

#### 6.3.3 Immutable Audit Logging
- Maintain an append-only ledger of all administrative actions:
  - Flag creation, modification, deletion
  - Environment creation and protection state changes
  - Role assignments and user invitations
  - Change Request lifecycle transitions
- Each audit log entry must include:
  - Timestamp
  - Actor user ID
  - Target project/environment/flag IDs
  - Previous JSON state
  - New JSON state
  - Actor IP address (where possible)
- Audit logs must be:
  - Queryable via UI
  - Exportable in CSV or similar format
  - Streamable via webhooks to external SIEM tools (e.g., Splunk, Datadog)

### 6.4 SDKs and Local Evaluation Engine

#### 6.4.1 SDK Language Coverage
- Initial release must include SDKs for:
  - Go
  - Java
  - Python
  - Node.js (JavaScript/TypeScript)
  - .NET (C#)
  - React and Next.js (frontend)
  - iOS (Swift)
  - Android (Kotlin/Java)
- SDKs must be designed for extension to additional languages.

#### 6.4.2 OpenFeature Compatibility
- SDK APIs must conform to the CNCF OpenFeature standard where feasible, to ease vendor migration and interoperability.
- Where OpenFeature semantics conflict with FlagManagment-specific features, document any deviations clearly.

#### 6.4.3 Server-Side Local Evaluation Protocol
- Bootstrapping:
  - On application startup, the server-side SDK opens a persistent secure streaming connection (e.g., gRPC) to the FlagManagment backend using the environment token.
  - The SDK downloads a full ruleset snapshot (all flags for that environment) into local memory.
- In-Memory Evaluation:
  - All flag evaluations occur in-memory within the SDK.
  - Evaluation must be designed to be sub-millisecond under normal conditions.
- Delta Updates:
  - When an administrator modifies flags, the backend pushes lightweight delta messages (changesets) over the stream.
  - SDKs apply deltas to their in-memory cache in real-time.
- Resiliency:
  - If the streaming connection is lost, SDKs must continue evaluating flags based on the last known good snapshot.
  - SDKs should implement exponential backoff and reconnection strategies.

### 6.5 Observability and Automated Rollback

#### 6.5.1 Telemetry Ingestion
- Provide endpoints and/or webhook handlers that can ingest alerts and metrics from external APM/monitoring systems (e.g., Datadog, New Relic, Prometheus).
- Support mapping external signals (e.g., error rate spikes) to specific flags or environments.

#### 6.5.2 Automated Action Triggers
- Allow engineers to bind telemetry thresholds to flag behaviors:
  - Example: “If HTTP 500 errors for service X exceed 2% for 5 minutes, execute kill-switch on flag Y in Production”.
- Implement a configurable rule engine for triggers:
  - Conditions based on metrics and time windows
  - Actions such as:
    - Set rollout percentage to 0%
    - Toggle flag off
    - Revert to previous known good configuration
- All automated actions must be recorded in the audit log.

### 6.6 Admin UI and Developer Experience
- Provide a web-based dashboard for:
  - Browsing projects/environments
  - Creating and editing flags and targeting rules via a visual rule builder
  - Viewing promotion pipelines and environment states
  - Managing change requests and approvals
  - Inspecting audit logs and export options
- UX should prioritize clarity around:
  - Which environment is currently active
  - The impact scope of changes
  - The status of change requests and automated triggers

## 7. Non-Functional Requirements

### 7.1 Performance
- Local evaluation in SDKs:
  - Single flag evaluation should typically complete in under 1 ms on commodity hardware.
- Backend throughput:
  - Must support thousands of concurrent SDK connections per instance.
  - Must support high rates of configuration changes without noticeable lag in delta propagation.

### 7.2 Scalability
- Horizontal scalability via stateless application nodes and shared data stores/caches.
- Ability to run separate clusters or tenants for large organizations.

### 7.3 Availability and Reliability
- Design for high availability with support for:
  - Multi-instance backend deployments
  - Failover strategies at the cache and database layers
- SDK resiliency:
  - Continued local evaluation during backend outages

### 7.4 Security
- Enforce secure transport (HTTPS/TLS) for all external traffic.
- Environment tokens must be stored and transmitted securely; avoid exposing them in client-side code for server-side environments.
- Provide hooks for authentication and authorization integration (e.g., SAML/OIDC, OAuth2) for the admin UI.

### 7.5 Compliance and Auditability
- Design audit logs and RBAC to support common compliance requirements (e.g., SOC 2-style change management).
- Provide configuration options for log retention policies.

### 7.6 Operability and Maintainability
- Provide structured logging and metrics from all services for monitoring.
- Offer admin tools or scripts for:
  - Database migrations
  - Archiving or cleaning up stale flags
  - Backup and restore

## 8. Baseline Data Model

The vendor may refine the schema, but the following entities and relationships must be present.

### 8.1 Core Tables
- **projects**
  - Fields: UUID, Name, CreatedAt, UpdatedAt
- **environments**
  - Fields: UUID, ProjectID (FK), Name, APIKeyHash, IsProtected (Boolean), CreatedAt, UpdatedAt
- **feature_flags**
  - Fields: UUID, ProjectID (FK), Key (String), Type (Enum: Boolean, Multivariate), ParentFlagID (FK, nullable), CreatedAt, UpdatedAt
- **environment_flag_states**
  - Junction table linking Environments and Feature Flags
  - Fields: EnvironmentID (FK), FeatureFlagID (FK), BooleanState, TargetingRules (JSONB), RemoteConfig (JSONB), CreatedAt, UpdatedAt

### 8.2 Governance Tables
- **change_requests**
  - Fields: UUID, EnvironmentID (FK), RequesterID (FK), Status (Enum: Pending, Approved, Rejected), ProposedDelta (JSONB), CreatedAt, UpdatedAt
- **change_request_approvals**
  - Fields: UUID, ChangeRequestID (FK), ApproverID (FK), ApprovedAt
- **audit_logs**
  - Fields: UUID, ActorID (FK), TargetProjectID (FK, nullable), TargetEnvironmentID (FK, nullable), TargetFlagID (FK, nullable), ActionType (String/Enum), PreviousState (JSONB), NewState (JSONB), ActorIP, CreatedAt

### 8.3 RBAC Tables (Example)
- **roles**
  - Fields: UUID, Name, Scope (Global/Project/Environment), CreatedAt, UpdatedAt
- **user_roles**
  - Fields: UUID, UserID, RoleID, ProjectID (FK, nullable), EnvironmentID (FK, nullable)

## 9. Delivery Expectations and Acceptance Criteria

### 9.1 Code Quality and Documentation
- Deliver well-structured, idiomatic code in chosen languages.
- Provide:
  - API documentation (OpenAPI/Swagger for REST; protobuf definitions for gRPC)
  - SDK documentation and example integrations
  - Deployment documentation (Docker/Kubernetes, configuration, scaling guidance)

### 9.2 Testing and Validation
- Automated tests must cover:
  - Flag evaluation logic (including multivariate and dependencies)
  - Promotion pipelines across environments
  - RBAC enforcement and protected environment behavior
  - Audit logging and change request workflows
  - SDK local evaluation behavior and delta synchronization
  - Observability triggers and automated actions
- Include end-to-end test scenarios demonstrating:
  - A typical feature rollout across Dev → QA → Staging → Production
  - A change request cycle in a protected environment
  - Automated rollback triggered by telemetry

### 9.3 Handover and Support
- Provide initial support period for:
  - Production deployment assistance (for self-hosted use)
  - Bug fixes for critical issues discovered immediately after delivery
- Handover must include:
  - Source code repositories
  - CI/CD pipelines (or scripts) used for build/test/deploy
  - All configuration required to run the platform in a standard cloud environment

---
# Appendix A
Here are the missing details you should weave into the specification:

### 1. Infrastructure-as-Code (IaC) and Automation

While section 5.2 covers containerization and Kubernetes, enterprise IT teams rarely configure environments or RBAC manually via a UI.

* **Add to Section 6.6 (Admin UI and Developer Experience):** Require the development of an official **Terraform Provider**. Teams utilizing monorepos for their cloud infrastructure will need to provision projects, environments, and base feature flags programmatically via Terraform to ensure disaster recovery and strict configuration management.
* **Add to Section 6.2 (Multi-Environment):** Specify support for API-driven environment cloning so CI/CD automation tools (like n8n or GitHub Actions) can spin up ephemeral environments for automated integration testing and tear them down afterward.

### 2. The Open-Source Contributor Experience

To ensure the project gains traction in the open-source community, the local development environment needs to be entirely frictionless.

* **Add to Section 9.1 (Code Quality and Documentation):** Mandate standard workspace configurations for modern environments like Windsurf and VS Code.
* **Specify Architecture Support:** The local containerized orchestration (e.g., Docker Compose) must be explicitly optimized to build and run seamlessly across varying hardware, from ARM-based Mac environments down to lightweight deployments capable of running on a local Raspberry Pi 4 for edge-testing scenarios.

### 3. Technical Debt and Stale Flag Management

We noted in our market research that a major pain point is the accumulation of dead flags, but there is no functional requirement to address *how* the system handles this.

* **Add a new subsection under 6.1 (Flag Types):** "Stale Flag Detection." The engine must track a `last_evaluated_at` timestamp for every flag. The dashboard should highlight flags that have been fully rolled out to 100% and haven't changed state in 30+ days, explicitly flagging them for codebase cleanup.

### 4. SDK Event Forwarding for A/B Analytics

Section 6.1.2 covers *how* multivariate flags are distributed, but not how the results are measured.

* **Add to Section 6.4 (SDKs):** Server-side and client-side SDKs must support a hook or interceptor pattern to broadcast evaluation events. If a user is bucketed into Variant B, the SDK needs a standardized way to forward that bucketing event to external product analytics tools (like Amplitude or PostHog) so product teams can calculate the actual impact of the A/B test.

### 5. Edge Proxy / Relay Node (Optional but Recommended)

For high-security corporate environments, databases and core backends are often locked down entirely.

* **Add to Section 5.1 (Architectural Goals):** Require the design of a stateless "Relay Node" or "Edge Proxy." This is a lightweight, containerized microservice that sits inside a private corporate subnet. It maintains the gRPC connection to the FlagManagment backend and serves as the local evaluation hub for internal microservices, preventing hundreds of internal SDKs from needing outbound internet access.

---


# Appendix B

Here is the formatted Markdown text. we added this to the abive existing PRD document to ensure the development squad has strict operational and execution boundaries.

## 10. Infrastructure, Automation & DevOps

### 10.1 Infrastructure-as-Code (IaC)
- The platform must be provisionable via a dedicated Terraform Provider.
- Organizations utilizing monorepos must be able to declare projects, environments, and base feature flags programmatically to ensure disaster recovery and strict configuration management.

### 10.2 CI/CD & Ephemeral Environments
- The backend API must support programmatic environment cloning.
- Automated orchestration workflows (e.g., n8n, GitHub Actions) must be able to spin up ephemeral test environments (e.g., `PR-123-Test`) for automated integration testing and cleanly tear them down upon completion.

### 10.3 Edge Proxy / Relay Node (Recommended)
- The architecture should support a stateless "Relay Node" or "Edge Proxy."
- This containerized microservice will sit inside a private corporate subnet, maintain the gRPC connection to the FlagManagment backend, and serve as the local evaluation hub for internal microservices, eliminating the need for all internal SDKs to have outbound internet access.

## 11. Open-Source Developer & Contributor Experience

### 11.1 Local Workspace Configuration
- The repository must include standardized, out-of-the-box workspace configurations for modern IDEs, explicitly supporting Windsurf and VS Code.
- Developers should be able to bootstrap the entire stack locally with a single command.

### 11.2 Multi-Architecture Support
- Local containerized orchestration (e.g., Docker Compose) must be optimized to build and run seamlessly across varying hardware profiles.
- Support must range from standard x86 and ARM-based Mac environments down to lightweight deployments capable of running on a local Raspberry Pi 4 for edge-testing and local automation scenarios.

## 12. Lifecycle & Analytics Integrations

### 12.1 Stale Flag Management
- The engine must track a `last_evaluated_at` timestamp and state-change history for every flag.
- The dashboard must proactively identify and highlight "stale" flags (e.g., flags rolled out to 100% with no state changes in 30+ days) to encourage codebase cleanup and reduce technical debt.

### 12.2 SDK Event Forwarding for Analytics
- Server-side and client-side SDKs must implement a standardized hook or interceptor pattern to broadcast evaluation events.
- When an identity is bucketed into a specific variant, the SDK must seamlessly forward that event to external product analytics tools (e.g., PostHog, Amplitude) to measure the operational impact of A/B tests.

## 13. Agile Execution & Phased Delivery Scope

### 13.1 Phase 1: Core Engine MVP
- Boolean and Multivariate flags, unlimited environments, core PostgreSQL schema, internal gRPC/REST APIs, and the primary React dashboard.
- Delivery of the local evaluation SDKs for at least two primary backend languages.

### 13.2 Phase 2: Governance & RBAC
- Implementation of granular RBAC matrices.
- Delivery of the Change Request entity, visual diffing, mandatory approval workflows for protected environments, and the immutable audit log.

### 13.3 Phase 3: Enterprise & Observability
- Telemetry ingestion endpoints and automated kill-switch triggers.
- Full OpenFeature API compliance and expansion of remaining language SDKs.

## 14. API-First Contract & UI Constraints

### 14.1 Contract-Driven Development
- The API must be strictly defined before business logic implementation. The development squad must deliver an OpenAPI 3.0 specification for REST and Protobuf definitions for gRPC.
- A mocked version of the API must be provided immediately to unblock parallel frontend development.

### 14.2 Design System Constraints
- The frontend dashboard must utilize a standardized, widely adopted open-source component library (e.g., Tailwind CSS with Shadcn/UI or MUI) to guarantee a clean, professional aesthetic without excessive custom CSS overhead.

## 15. Quality Assurance & Definition of Done (DoD)

### 15.1 Code Quality & Testing Minimums
- Mandatory automated test coverage thresholds: 80% unit test coverage for the backend engine, 70% for the frontend UI.
- Code must pass standard static analysis and linting pipelines with zero critical warnings before merge.

### 15.2 Performance Baselines
- The server-side SDK local evaluation function must demonstrably execute in under 1 ms during simulated load testing before the feature can be marked as complete.

### 15.3 Data Privacy & SOC2 Readiness
- The system must natively salt and hash any Personal Identifiable Information (PII) utilized for identity bucketing prior to database storage.
- Audit logs must be actively sanitized to ensure no plaintext API keys or sensitive targeting metadata are inadvertently captured.