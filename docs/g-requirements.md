# Comprehensive Product Requirements Document (PRD)

**Project Name:** FlagManagment (Working Title)
**Document Purpose:** To serve as the standalone, single source of truth for engineering, architecture, and product teams to design and build a next-generation feature flag and remote configuration platform.
**Target Audience:** Software Architects, Developers, Open-Source Contributors, and Product Analysts.

---

## 1. Product Vision & Executive Summary

FlagManagment is a cloud-native, open-source feature flag management and remote configuration platform. It is engineered to bridge the gap between basic open-source toggles and prohibitively expensive enterprise suites. Existing platforms intentionally cripple their open-source offerings by gating critical multi-environment workflows, granular governance, and compliance features behind high-tier paywalls.

FlagManagment rejects this model. It provides unlimited scaling, multi-environment promotion pipelines, and enterprise-grade role-based access control (RBAC) completely out of the box for self-hosted deployments. The platform is designed to be fully decoupled from application deployments, allowing IT, QA, and Product teams to manage feature rollouts dynamically and safely.

---

## 2. Licensing & Business Model

The core problem with current open-source feature flag tools is the vulnerability of the maintainers to major cloud providers, balanced against corporate reluctance to adopt restrictive copyleft licenses.

* **License Type:** Business Source License (BSL) or a comparable "Fair-Source" model.
* **Corporate Freedom:** Organizations can self-host, modify, and utilize the entire platform internally for their own business operations with zero licensing costs and zero artificial feature caps.
* **Integrations & Ecosystem**
  - Webhooks for flag changes (creation, modification, deletion, toggle).
  - OpenFeature compliant SDKs for seamless integration.
  - Supported Language SDKs:
    - **Server-side**: Go, Java, Python, .NET (Node.js/TypeScript native)
    - **Client-side**: React (Web), iOS (Swift), Android (Kotlin)
* **Commercial Protection:** The license explicitly and legally prohibits third-party managed service providers (e.g., major cloud vendors) from wrapping the codebase and reselling it as a competing public Software-as-a-Service (SaaS).
* **Monetization Strategy:** The founding organization will offer a fully managed SaaS cloud version. Revenue will be generated through managed hosting, high-throughput API SLAs, and premium integrations, rather than intentionally limiting the fundamental security and governance tools needed by modern engineering teams.

---

## 3. Market Differentiation & Competitive Gap Analysis

FlagManagment is built specifically to solve the deliberate bottlenecks engineered into current market leaders.

### The Unleash Problem

* Unleash caps its open-source version at exactly 1 project and 2 environments.


* This artificially forces teams operating standard pipelines (Local, Dev, QA, Staging, Production) into premium tiers.


* FlagManagment will support unlimited projects and unlimited environments natively.

### The Flagsmith Problem

* Flagsmith restricts granular RBAC, environment-level permissions, and approval workflows to their "Scale-Up" and "Enterprise" tiers.


* Teams requiring compliance or separated duties are forced to pay premium prices.


* FlagManagment makes granular RBAC, change requests, and immutable audit logs core, non-paywalled features.

---

## 4. Proposed Technology Stack & Architecture

To guarantee the system is cloud-native, highly performant, and attractive to the global open-source community, the architecture must utilize modern, widely supported tooling.

| Component | Selected Technology | Technical Rationale |
| --- | --- | --- |
| **Backend Engine** | Go (Golang) | Industry standard for cloud-native infrastructure. Offers exceptional memory management, rapid execution, and high concurrency via goroutines for serving millions of local evaluation SDK requests. |
| **Frontend Dashboard** | React + TypeScript + Vite | Provides a massive developer ecosystem for future open-source contributors. TypeScript ensures strict safety for complex targeting rule builders, and Vite delivers rapid build times. |
| **Primary Datastore** | PostgreSQL | Required for complex relational schemas mapping RBAC hierarchies, project/environment isolation, and immutable audit logs. JSONB column support allows for schema-less storage of multivariate flag targeting rules. |
| **Caching & Pub/Sub** | Redis | Acts as a high-throughput, low-latency caching layer to shield PostgreSQL from heavy read loads and powers the real-time event broadcasting to server-side SDKs. |
| **API Protocols** | gRPC & REST | Internal communications and Server-Side SDK streaming will utilize gRPC/Protobuf for maximum throughput. Standard REST APIs will be exposed for external tooling and UI interactions. |

---

## 5. Core Feature Requirements

### Flag Types and Evaluation Rules

* **Boolean Toggles:** Standard on/off switches for code execution paths.
* **Multivariate Flags:** A/B/n testing variants with percentage-based bucketing that is deterministic based on identity hashes.
* **Remote Configuration:** Flags must return strongly typed payload values (Strings, Numbers, Booleans, or complex JSON objects) to allow runtime configuration changes without redeploying code.


* **Sequential Dependencies:** A flag must be configurable to depend on the state of a parent flag. If the parent flag is disabled, the dependent flag automatically defaults to its safe fallback state, preventing architectural conflicts.


* **Contextual Targeting:** Evaluation must support operators like equals, contains, regex matching, and array inclusions against custom user attributes (e.g., tenant ID, region, email domain).

### Multi-Environment Promotion Pipelines

* **Strict Isolation:** Every environment must generate its own unique, cryptographically secure SDK authentication tokens.
* **State Separation:** Flag states, targeting rules, and rollout percentages must exist independently within each environment.
* **One-Click Promotions:** The UI and API must provide a "Promote" action to copy a flag's exact rule configuration seamlessly from a lower environment (e.g., QA) to a higher environment (e.g., Staging) without manual data entry.

---

## 6. Governance, RBAC, and Compliance

This is the primary differentiator of FlagManagment. These features must be deeply integrated into the data model from day one.

### Granular Role-Based Access Control (RBAC)

* Permissions must be calculable at the global, project, and individual environment levels.
* Roles must allow separation of duties. A QA engineer can have full read/write access to toggle flags in testing environments but strict read-only access in production.

### Change Requests and Mandatory Approvals

* Environments can be flagged as "Protected" (e.g., Production).
* Any mutation to a protected environment cannot apply immediately. The action must automatically generate a "Change Request."
* The Change Request must display a git-style visual diff comparing the current rule matrix against the proposed rule matrix.
* The change remains pending until a user with a designated "Release Manager" role reviews and executes the approval.

### Immutable Audit Logging

* The system must record an append-only ledger of all administrative actions, flag creations, rule modifications, and user invitations.
* Logs must strictly capture the timestamp, exact user ID, target environment, previous JSON state, new JSON state, and the actor's IP address.
* Logs must be exportable via CSV and capable of being streamed in real-time via webhooks to external SIEM tools like Splunk or Datadog.

---

## 7. SDKs and the Local Evaluation Engine

To ensure FlagManagment is viable for massive enterprise applications, it cannot introduce network latency into application logic.

### Broad Language Support

* The platform must launch with SDKs for primary backend languages (Go, Java, Python, Node.js, .NET) and frontend frameworks (React, Next.js, iOS, Android).
* All SDKs must conform strictly to the CNCF OpenFeature API standard to allow easy vendor migration.



### Server-Side Local Evaluation Engine Protocol

* **Bootstrapping:** Upon application startup, the Server-Side SDK uses its environment token to open a persistent streaming RPC connection to the FlagManagment server. It downloads the complete ruleset snapshot into local memory.
* **In-Memory Evaluation:** When the application queries a flag variant, the SDK evaluates the rules locally in memory, guaranteeing sub-millisecond response times without executing any outbound network calls.
* **Delta Updates:** When an administrator alters a flag in the FlagManagment dashboard, the server pushes an instantaneous, lightweight delta patch down the open gRPC stream to all connected SDKs, updating their local cache in real-time.
* **Resiliency:** If the connection to the FlagManagment server drops, the SDK must continue evaluating flags accurately using its last known good memory state, preventing application outages.

---

## 8. Observability and Automated Rollback

FlagManagment must integrate natively with modern telemetry systems to act as an automated safety net for deployments.

### Telemetry Ingestion

* The platform must expose endpoints designed to consume webhooks and alerts from external APM platforms (e.g., Datadog, New Relic, Prometheus).

### Automated Action Triggers

* Engineers must be able to bind specific active telemetry thresholds (e.g., "HTTP 500 errors spike above 2%") to a specific feature flag rollout.
* If the alert fires, FlagManagment must automatically execute a predetermined safety action, such as executing an immediate "Kill-Switch" rollback to 0% distribution to stop a bad deployment before human intervention is required.



---

## 9. Baseline Database Schema Requirements

The PostgreSQL schema must be designed for strict relational integrity.

### Core Tables

* `projects`: UUID, Name, Timestamps.
* `environments`: UUID, Project ID (FK), Name, API Key Hash, Is Protected Boolean.
* `feature_flags`: UUID, Project ID (FK), Key (String), Type (Boolean/Multivariate), Parent Flag ID (FK).
* `environment_flag_states`: Junction table connecting Environments to Feature Flags. Contains Boolean State, Targeting Rules (JSONB), and Remote Config (JSONB).

### Governance Tables

* `change_requests`: UUID, Environment ID (FK), Requester ID, Status (Pending/Approved/Rejected), Proposed Delta (JSONB).
* `change_request_approvals`: Tracks the exact user who authorized the change request.
* `audit_logs`: UUID, Target IDs, Actor ID, Action Type, Previous State, New State, Actor IP.