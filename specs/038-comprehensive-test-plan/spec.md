# Feature Specification: Comprehensive Test Plan

**Feature Branch**: `038-comprehensive-test-plan`

**Created**: 2026-08-19

**Status**: Active

**Input**: "Write a list of the test cases to test each feature of the application using g-requirements.md, p-requirements.md, and all feature specs."

---

## Purpose

This document is the **canonical test plan** for the FlagManagment platform. It covers all features derived from `g-requirements.md`, `p-requirements.md`, and all 37 feature specs (001–037). Test cases are written in **Given/When/Then format** for use by both human testers and AI testing agents. Each case has a unique stable ID that must not be renumbered.

**Total test cases: 158** across 24 areas + 6 end-to-end scenarios.

**Layers covered**: API · SDK · UI · Infrastructure · CI/CD · IaC

**Source specs referenced**: All 37 specs in `/specs/` directory plus `docs/g-requirements.md` and `docs/p-requirements.md`.

---

## Part 1 — Authentication & User Management

### TC-AUTH-001: Login with Valid Credentials
**Layer**: UI + API | **Spec**: 027-login-error-feedback, 007-rbac-user-management

- **Given** a registered user with valid credentials exists
- **When** the user submits the login form with correct email and password
- **Then** the user is authenticated, a session token is issued, and the dashboard is displayed

### TC-AUTH-002: Login with Invalid Credentials
**Layer**: UI | **Spec**: 027-login-error-feedback

- **Given** a login form is shown
- **When** the user enters an incorrect password and submits
- **Then** a clear, specific error message is shown (not a generic 500) and login is rejected

### TC-AUTH-003: Login with Non-Existent Email
**Layer**: UI + API

- **Given** a login form is shown
- **When** the user enters an email that does not exist in the system
- **Then** a "user not found" or generic auth failure message is shown (no account enumeration)

### TC-AUTH-004: Session Expiry and Re-Authentication
**Layer**: API + UI

- **Given** a user is logged in
- **When** the session token expires
- **Then** the user is redirected to login and a clear expiry message is shown

### TC-AUTH-005: Invite New User (Admin Flow)
**Layer**: API + UI | **Spec**: 007-rbac-user-management, 034-user-management-api

- **Given** a System Administrator is logged in
- **When** they invite a new user by email and assign a role
- **Then** an invitation email is sent, the invite is recorded in the audit log, and the user appears as "Pending" in the user list

### TC-AUTH-006: Accept Invitation and Set Password
**Layer**: UI

- **Given** a user has received an invitation email
- **When** they click the invite link and set a password
- **Then** their account is activated and they are redirected to the dashboard with their assigned role

### TC-AUTH-007: Prevent Duplicate User Invitation
**Layer**: API

- **Given** a user with email `user@example.com` already exists
- **When** an admin tries to invite the same email again
- **Then** the API returns a 409 Conflict error with a descriptive message

### TC-AUTH-008: Role Assignment at Invitation
**Layer**: API | **Spec**: 034-user-management-api

- **Given** an admin invites a user with role "QA Engineer"
- **When** the user accepts the invitation
- **Then** the user has QA Engineer permissions: write in Dev/QA environments, read-only in Production

### TC-AUTH-009: Enterprise SSO Login (SAML/OIDC)
**Layer**: UI + API | **Spec**: 022-enterprise-sso

- **Given** SSO is configured with an external IdP
- **When** a user clicks "Login with SSO" and authenticates via the IdP
- **Then** a session is established, user roles are mapped from IdP attributes, and the user lands on the dashboard

### TC-AUTH-010: SSO Role Mapping
**Layer**: API | **Spec**: 022-enterprise-sso

- **Given** IdP returns group claims for "engineering-admin"
- **When** the SSO assertion is processed
- **Then** the user is assigned the mapped FlagManagment role (e.g., Release Manager) automatically

---

## Part 2 — Project Management

### TC-PROJ-001: Create New Project
**Layer**: API + UI | **Spec**: 001-core-architecture, 003-core-api

- **Given** an authenticated admin is on the Projects page
- **When** they create a new project with a unique name
- **Then** the project is created with a UUID, appears in the project list, and is recorded in the audit log

### TC-PROJ-002: Create Project with Duplicate Name
**Layer**: API

- **Given** a project named "MyApp" exists
- **When** the user attempts to create another project named "MyApp"
- **Then** the API returns a 409 error; the duplicate is not created

### TC-PROJ-003: List All Projects
**Layer**: API + UI

- **Given** multiple projects exist
- **When** an authenticated user navigates to the Projects list
- **Then** all projects the user has access to are listed with name, creation date, and environment count

### TC-PROJ-004: Delete Project
**Layer**: API | **Spec**: 001-core-architecture

- **Given** a project exists with environments and flags
- **When** an admin deletes the project
- **Then** the project, all its environments, flags, and flag states are removed; an audit log entry is created

### TC-PROJ-005: Unlimited Projects (No Artificial Cap)
**Layer**: API | **Ref**: g-requirements §4.1

- **Given** a self-hosted deployment
- **When** more than 1 project is created
- **Then** all projects are accepted with no licensing enforcement blocking creation

---

## Part 3 — Environment Management

### TC-ENV-001: Create Environment in Project
**Layer**: API + UI | **Spec**: 002-data-model-state, 003-core-api

- **Given** a project exists
- **When** a user creates a new environment (e.g., "QA") within the project
- **Then** the environment is created with a unique cryptographically secure SDK API key, and appears in the environment list

### TC-ENV-002: Unlimited Environments Per Project (No Cap)
**Layer**: API | **Ref**: g-requirements §4.1

- **Given** a project already has environments: Dev, QA, Staging, Production
- **When** the user creates a 5th, 6th, or Nth environment
- **Then** all are created successfully with no error or enforcement limit

### TC-ENV-003: Environment SDK Key Is Unique and Secure
**Layer**: API | **Spec**: 032-environment-sdk-key-ux, 033-server-side-keys-management

- **Given** two environments are created in the same project
- **When** their SDK keys are compared
- **Then** the keys are cryptographically different; using Environment A's key cannot access Environment B's flags

### TC-ENV-004: Mark Environment as Protected
**Layer**: API + UI | **Spec**: 008-change-requests-workflow

- **Given** an environment named "Production" exists
- **When** an admin sets `isProtected = true`
- **Then** any subsequent mutation to flags in Production creates a Change Request rather than applying directly

### TC-ENV-005: Environment Isolation — Cross-Access Prevention
**Layer**: API + SDK | **Spec**: 005-sdk-evaluation-api

- **Given** SDK_KEY_DEV and SDK_KEY_PROD are different
- **When** an SDK configured with SDK_KEY_DEV requests flag state
- **Then** only Dev environment flags are returned; Production flag states are inaccessible

### TC-ENV-006: Delete Environment
**Layer**: API

- **Given** an environment exists with flags
- **When** an admin deletes the environment
- **Then** the environment and all its flag states are removed; the audit log captures the action

### TC-ENV-007: Clone Environment (Ephemeral Environments)
**Layer**: API | **Spec**: 018-environment-cloning

- **Given** a "QA" environment with flags and targeting rules exists
- **When** an API call clones the QA environment to "PR-456-Test"
- **Then** a new environment is created with the same flag configurations, a new unique SDK key, and it appears in the environment list

### TC-ENV-008: Delete Ephemeral/Cloned Environment via API
**Layer**: API | **Spec**: 018-environment-cloning

- **Given** an ephemeral environment "PR-456-Test" exists
- **When** the API call to delete it is made
- **Then** the environment is deleted with all its flag states; the API returns 204

### TC-ENV-009: Context Switching in UI
**Layer**: UI | **Spec**: 030-env-context-switching

- **Given** a project has multiple environments
- **When** the user selects a different environment from the environment switcher
- **Then** all flag lists, states, and targeting rules displayed are specific to the newly selected environment

### TC-ENV-010: SDK Key Display and Copy in UI
**Layer**: UI | **Spec**: 032-environment-sdk-key-ux

- **Given** an environment exists
- **When** the user navigates to the environment SDK Key section
- **Then** the SDK key is displayed with a one-click copy button and an obfuscated masked view by default

---

## Part 4 — Feature Flags (Boolean & Multivariate)

### TC-FLAG-001: Create Boolean Feature Flag
**Layer**: API + UI | **Spec**: 003-core-api, 035-flag-creation-targeting-workflows

- **Given** a project and environment exist
- **When** a user creates a new boolean flag with key `new-checkout-flow`
- **Then** the flag is created, stored in the database, and defaults to `disabled` in all environments

### TC-FLAG-002: Toggle Boolean Flag On
**Layer**: API + UI | **Spec**: 029-flag-toggle-feedback

- **Given** a boolean flag exists in disabled state
- **When** the user toggles the flag to enabled
- **Then** the flag state updates to `true`, the change is reflected in the UI, an SDK delta is pushed, and the action is audit logged

### TC-FLAG-003: Toggle Boolean Flag Off
**Layer**: API + UI

- **Given** a boolean flag is enabled
- **When** the user toggles it to disabled
- **Then** the flag returns `false`, the SDK delta is pushed, and the toggle is logged in the audit trail

### TC-FLAG-004: Toggle Feedback — Visual Confirmation
**Layer**: UI | **Spec**: 029-flag-toggle-feedback

- **Given** a flag exists in the dashboard
- **When** the user toggles its state
- **Then** a visual confirmation (e.g., spinner, success toast) is shown before and after the API response with no ambiguous state

### TC-FLAG-005: Create Multivariate Flag with Variants
**Layer**: API + UI | **Spec**: 013-multivariate-flags, 035-flag-creation-targeting-workflows

- **Given** a project and environment exist
- **When** a user creates a multivariate flag with variants: "control" (50%), "variant-A" (30%), "variant-B" (20%)
- **Then** the flag is created with three named variants whose percentages sum to 100%

### TC-FLAG-006: Multivariate Variant Percentages Must Sum to 100%
**Layer**: API + UI | **Spec**: 013-multivariate-flags

- **Given** a multivariate flag creation form
- **When** the user saves variants with percentages summing to 90%
- **Then** the API rejects the request with a validation error (400) and the UI displays the error

### TC-FLAG-007: Deterministic Bucketing for Multivariate Flags
**Layer**: SDK | **Spec**: 013-multivariate-flags

- **Given** a multivariate flag with variants A (50%) and B (50%) is active
- **When** the SDK evaluates the flag for the same user identity 1000 times
- **Then** the same variant is returned every time (deterministic hash-based bucketing)

### TC-FLAG-008: Stable Assignment Across Sessions
**Layer**: SDK | **Spec**: 013-multivariate-flags

- **Given** user ID "user-abc" was bucketed into "variant-A" in session 1
- **When** user-abc evaluates the same flag in a new session
- **Then** "variant-A" is returned (stable assignment persists across sessions)

### TC-FLAG-009: Remote Configuration — String Payload
**Layer**: API + SDK | **Spec**: 026-remote-config-ui

- **Given** a flag with a remote config string payload `{"theme": "dark"}` exists
- **When** the SDK evaluates the flag
- **Then** the SDK returns the string payload without code redeployment

### TC-FLAG-010: Remote Configuration — JSON Payload
**Layer**: API + SDK | **Spec**: 026-remote-config-ui

- **Given** a flag with JSON remote config `{"maxRetries": 3, "timeout": 5000}` exists
- **When** the SDK evaluates the flag
- **Then** the strongly typed JSON object is returned intact

### TC-FLAG-011: Remote Configuration — Update Payload Without Redeploy
**Layer**: API + SDK | **Spec**: 026-remote-config-ui

- **Given** a remote config flag has value `{"limit": 10}` and the SDK is connected
- **When** an admin updates the payload to `{"limit": 20}` in the dashboard
- **Then** the SDK receives a delta update and the new value is returned within seconds, without any application restart

### TC-FLAG-012: Update Flag Name and Description
**Layer**: API + UI

- **Given** a flag exists
- **When** the user edits the flag's name or description
- **Then** the changes are saved, the flag key remains unchanged, and the audit log captures the modification

### TC-FLAG-013: Delete Feature Flag
**Layer**: API + UI

- **Given** a feature flag exists
- **When** an admin deletes it
- **Then** the flag is removed from all environments, SDK evaluations return the default fallback value, and the deletion is audit logged

### TC-FLAG-014: List Flags with Pagination
**Layer**: API + UI

- **Given** 50+ feature flags exist in a project
- **When** the user views the flag list
- **Then** flags are paginated or virtualized; navigation between pages works correctly

### TC-FLAG-015: Search and Filter Flags
**Layer**: UI | **Spec**: 004-frontend-dashboard

- **Given** multiple flags exist with various names
- **When** the user types a search term
- **Then** only matching flags are shown in real-time

### TC-FLAG-016: Flag Creation Empty State / Onboarding
**Layer**: UI | **Spec**: 028-empty-state-onboarding

- **Given** a project has no flags
- **When** the user opens the Flags section
- **Then** an empty state with a clear CTA ("Create your first flag") is shown

---

## Part 5 — Contextual Targeting Engine

### TC-TARGET-001: Create Targeting Rule — Equals Operator
**Layer**: API + UI | **Spec**: 012-contextual-targeting, 035-flag-creation-targeting-workflows

- **Given** a boolean flag exists
- **When** a targeting rule is added: "IF `country` equals `US` THEN enable flag"
- **Then** the rule is saved; identities with `country=US` evaluate to `true`, others to `false`

### TC-TARGET-002: Create Targeting Rule — Contains Operator
**Layer**: API + SDK | **Spec**: 012-contextual-targeting

- **Given** a flag has rule: "IF `email` contains `@acme.com` THEN enable"
- **When** an identity with `email=jane@acme.com` evaluates the flag
- **Then** the flag returns `true`

### TC-TARGET-003: Create Targeting Rule — Regex Match
**Layer**: API + SDK | **Spec**: 012-contextual-targeting

- **Given** a flag has rule: "IF `userId` matches regex `^admin-.*`"
- **When** an identity with userId `admin-007` evaluates
- **Then** the flag returns `true`; an identity with userId `user-007` returns `false`

### TC-TARGET-004: Create Targeting Rule — Array Inclusion
**Layer**: API + SDK | **Spec**: 012-contextual-targeting

- **Given** a flag has rule: "IF `planTier` in `[premium, enterprise]`"
- **When** identities with planTier `free`, `premium`, and `enterprise` evaluate
- **Then** `premium` and `enterprise` receive `true`; `free` receives `false`

### TC-TARGET-005: Create Targeting Rule — Not Equals Operator
**Layer**: API + SDK | **Spec**: 012-contextual-targeting

- **Given** a flag has rule: "IF `region` not equals `EU`"
- **When** identity with `region=US` and `region=EU` evaluate
- **Then** US returns `true`, EU returns `false`

### TC-TARGET-006: Multiple Targeting Rules — AND Logic
**Layer**: API + SDK

- **Given** a flag has two rules: `country=US` AND `plan=premium`
- **When** evaluations occur for combinations
- **Then** only identities satisfying both conditions receive `true`

### TC-TARGET-007: Targeting Rule with Fallback Default
**Layer**: SDK

- **Given** a flag has a rule matching only premium users
- **When** a free-tier user evaluates the flag
- **Then** the flag falls back to the default value (disabled)

### TC-TARGET-008: Visual Rule Builder — Add and Preview Rule
**Layer**: UI | **Spec**: 035-flag-creation-targeting-workflows

- **Given** the flag editing UI is open
- **When** the user uses the visual rule builder to add a rule
- **Then** the rule is shown in a human-readable preview and can be saved

### TC-TARGET-009: Delete Targeting Rule
**Layer**: API + UI

- **Given** a flag has a targeting rule
- **When** the user removes the rule
- **Then** the rule is deleted; subsequent evaluations use the default state; change is audit logged

### TC-TARGET-010: Targeting Rules Stored as Structured JSON
**Layer**: API

- **Given** targeting rules are created via the UI
- **When** the API retrieves the flag state
- **Then** the rules are returned in structured JSONB format (e.g., `{"operator":"equals","attribute":"country","value":"US"}`)

---

## Part 6 — Percentage Rollouts

### TC-ROLLOUT-001: Create Percentage Rollout (Gradual)
**Layer**: API + UI | **Spec**: 013-multivariate-flags

- **Given** a boolean flag exists
- **When** the user sets a 10% gradual rollout
- **Then** approximately 10% of identities (by hash) evaluate to `true`; the rest evaluate to `false`

### TC-ROLLOUT-002: Gradual Rollout Distribution — Statistical Accuracy
**Layer**: SDK

- **Given** a flag is set to 50% rollout
- **When** 10,000 unique user IDs evaluate the flag
- **Then** approximately 50% (±5%) receive `true`, with distribution based on deterministic identity hash

### TC-ROLLOUT-003: Change Rollout Percentage in Real-Time
**Layer**: API + SDK

- **Given** a flag is at 10% rollout and SDKs are connected
- **When** an admin changes rollout to 50%
- **Then** the new percentage is pushed via delta update; SDKs begin distributing at the new rate within seconds

### TC-ROLLOUT-004: Set Rollout to 100% (Full Release)
**Layer**: API

- **Given** a flag at 50% rollout
- **When** the rollout is set to 100%
- **Then** all identities evaluate to `true`

### TC-ROLLOUT-005: Set Rollout to 0% (Kill Switch)
**Layer**: API

- **Given** a flag at 100% rollout
- **When** the rollout is set to 0%
- **Then** all identities evaluate to `false` regardless of other targeting rules

---

## Part 7 — Sequential Flag Dependencies

### TC-DEP-001: Create Dependent Flag (Child Inherits Parent State)
**Layer**: API + SDK | **Spec**: 016-sequential-dependencies

- **Given** a parent flag "feature-X" is disabled
- **When** a child flag "sub-feature-X" is configured to depend on "feature-X"
- **Then** evaluating "sub-feature-X" returns the safe fallback value (disabled) regardless of its own state

### TC-DEP-002: Enable Parent — Child Evaluates Normally
**Layer**: SDK | **Spec**: 016-sequential-dependencies

- **Given** a parent flag is enabled and a child depends on it
- **When** the child flag is evaluated
- **Then** the child evaluates using its own state/rules (not forced to fallback)

### TC-DEP-003: Disable Parent — All Children Fall Back
**Layer**: SDK | **Spec**: 016-sequential-dependencies

- **Given** multiple child flags depend on one parent flag
- **When** the parent is disabled
- **Then** all children immediately return their safe fallback values via delta update

### TC-DEP-004: Circular Dependency Prevention at Creation
**Layer**: API | **Spec**: 016-sequential-dependencies

- **Given** Flag A depends on Flag B
- **When** the user tries to make Flag B depend on Flag A
- **Then** the API rejects the request with a 400 error: "Circular dependency detected"

### TC-DEP-005: Multi-Level Dependency Chain
**Layer**: SDK | **Spec**: 016-sequential-dependencies

- **Given** Flag A → Flag B → Flag C (chain)
- **When** Flag A is disabled
- **Then** both Flag B and Flag C fall back to their safe defaults

---

## Part 8 — Scheduled Flag Changes

### TC-SCHED-001: Schedule Flag Enable at Future Time
**Layer**: API + UI | **Spec**: 014-scheduled-flags

- **Given** a flag is currently disabled
- **When** a schedule is created to enable the flag at a specific future timestamp
- **Then** before the scheduled time the flag is disabled; at the scheduled time the flag automatically enables and an audit log entry is created

### TC-SCHED-002: Schedule Flag Disable (Sunset)
**Layer**: API | **Spec**: 014-scheduled-flags

- **Given** a flag is enabled
- **When** a schedule is created to disable it at a future timestamp
- **Then** the flag automatically disables at the scheduled time

### TC-SCHED-003: Cancel Scheduled Change
**Layer**: API + UI | **Spec**: 014-scheduled-flags

- **Given** a scheduled flag change exists
- **When** the user cancels the schedule
- **Then** the schedule is removed; the flag retains its current state

### TC-SCHED-004: Scheduled Change in Protected Environment Requires Approval
**Layer**: API | **Spec**: 014-scheduled-flags, 008-change-requests-workflow

- **Given** Production is a protected environment
- **When** a user schedules a flag change in Production
- **Then** the schedule creates a pending Change Request that must be approved before it can take effect

---

## Part 9 — Flag Promotion Pipeline

### TC-PROMO-001: Promote Flag from QA to Staging
**Layer**: API + UI | **Spec**: 011-flag-promotions

- **Given** a flag has specific targeting rules in the QA environment
- **When** the user triggers "Promote to Staging"
- **Then** the flag state, targeting rules, rollout percentage, and remote config are copied to Staging; QA is unchanged; action is audit logged

### TC-PROMO-002: Promote to Protected Environment Requires Change Request
**Layer**: API + UI | **Spec**: 011-flag-promotions, 008-change-requests-workflow

- **Given** Production is a protected environment
- **When** the user promotes a flag from Staging to Production
- **Then** a Change Request is created with a diff showing proposed vs current Production state; the change is not applied until approved

### TC-PROMO-003: Promotion Diff View Shows Correct Changes
**Layer**: UI | **Spec**: 011-flag-promotions, 008-change-requests-workflow

- **Given** a flag in QA has different targeting rules than in Production
- **When** a promotion is initiated to Production
- **Then** the diff view clearly shows current rules (removed) and proposed rules (added) in a git-diff visual format

### TC-PROMO-004: Full Flag Payload Promotion (All Fields)
**Layer**: API | **Spec**: 011-flag-promotions

- **Given** a flag has rollout %, targeting rules, and remote config in QA
- **When** promoted to Staging
- **Then** all three data sets (rollout, rules, remote config) are replicated exactly to Staging

### TC-PROMO-005: Selective Field Promotion
**Layer**: API | **Spec**: 011-flag-promotions

- **Given** a flag has both targeting rules and remote config in QA
- **When** the user promotes only the targeting rules to Staging
- **Then** only the targeting rules are copied; the remote config in Staging is unchanged

---

## Part 10 — Change Requests & Approval Workflow

### TC-CR-001: Mutation in Protected Environment Creates Change Request
**Layer**: API + UI | **Spec**: 008-change-requests-workflow

- **Given** Production is a protected environment
- **When** a Release Manager-level user toggles a flag
- **Then** a Change Request is created with status "Pending"; the flag state does not change immediately

### TC-CR-002: Change Request Diff Displays Current vs Proposed State
**Layer**: UI | **Spec**: 008-change-requests-workflow

- **Given** a Change Request exists for a toggle action
- **When** a Release Manager opens the Change Request
- **Then** a git-style diff showing previous state (current) vs proposed state is displayed clearly

### TC-CR-003: Approve Change Request Applies Changes Atomically
**Layer**: API | **Spec**: 008-change-requests-workflow

- **Given** a Change Request is in "Pending" status
- **When** a Release Manager approves it
- **Then** the flag change is applied atomically; the CR status becomes "Approved"; the approval is audit logged with approver ID and timestamp

### TC-CR-004: Reject Change Request Discards Changes
**Layer**: API | **Spec**: 008-change-requests-workflow

- **Given** a Change Request is in "Pending" status
- **When** a Release Manager rejects it
- **Then** no changes are applied; the CR status becomes "Rejected"; the rejection is audit logged

### TC-CR-005: Only Release Manager Can Approve/Reject
**Layer**: API | **Spec**: 008-change-requests-workflow, 007-rbac-user-management

- **Given** a Change Request exists
- **When** a QA Engineer (non-Release-Manager) tries to approve it via API
- **Then** the API returns 403 Forbidden

### TC-CR-006: Change Request Requires Justification (if configured)
**Layer**: UI | **Spec**: 008-change-requests-workflow

- **Given** an environment policy requires approval justification
- **When** a Release Manager approves without providing justification
- **Then** the UI blocks submission and shows a validation message

### TC-CR-007: List All Pending Change Requests
**Layer**: API + UI | **Spec**: 008-change-requests-workflow

- **Given** multiple change requests exist
- **When** a Release Manager views the Change Requests list
- **Then** all pending CRs are shown with requester, environment, flag name, and timestamps

---

## Part 11 — RBAC & Permissions

### TC-RBAC-001: System Administrator Has Full Access
**Layer**: API | **Spec**: 007-rbac-user-management

- **Given** a user has the System Administrator role
- **When** they perform any operation (create project, delete flag, approve CR)
- **Then** all operations succeed without permission errors

### TC-RBAC-002: QA Engineer Can Write in Dev/QA Environments
**Layer**: API | **Spec**: 007-rbac-user-management

- **Given** a user has the QA Engineer role
- **When** they toggle a flag in the Dev or QA environment
- **Then** the toggle succeeds

### TC-RBAC-003: QA Engineer Cannot Write in Production
**Layer**: API | **Spec**: 007-rbac-user-management

- **Given** a user has the QA Engineer role
- **When** they attempt to toggle a flag in Production via API
- **Then** the API returns 403 Forbidden

### TC-RBAC-004: Read-Only Auditor Cannot Modify Flags
**Layer**: API | **Spec**: 007-rbac-user-management

- **Given** a user has the Read-Only Auditor role
- **When** they attempt to create, update, or delete a flag via API
- **Then** all mutation endpoints return 403 Forbidden

### TC-RBAC-005: Project-Level Role Scoping
**Layer**: API | **Spec**: 007-rbac-user-management

- **Given** a user is a Project Owner in Project A but has no role in Project B
- **When** they access flags in Project B
- **Then** access is denied (403); they can access flags in Project A normally

### TC-RBAC-006: Environment-Level Permission Override
**Layer**: API | **Spec**: 007-rbac-user-management

- **Given** a user has write access at project level but read-only in the Production environment specifically
- **When** they write to flags in other environments
- **Then** writes succeed; writing to Production returns 403

### TC-RBAC-007: Role Change Takes Effect Immediately
**Layer**: API + UI | **Spec**: 007-rbac-user-management

- **Given** a user is a QA Engineer
- **When** an admin changes their role to Release Manager
- **Then** the user immediately gains Change Request approval capabilities without re-login

---

## Part 12 — Immutable Audit Logs

### TC-AUDIT-001: Flag Toggle Is Audit Logged
**Layer**: API | **Spec**: 020-audit-logs-siem

- **Given** a user toggles a flag
- **When** the audit log is queried
- **Then** an entry exists containing: timestamp, actor user ID, action type, previous state (JSON), new state (JSON), target environment/flag IDs, and actor IP address

### TC-AUDIT-002: User Invitation Is Audit Logged
**Layer**: API | **Spec**: 020-audit-logs-siem

- **Given** an admin invites a new user
- **When** the audit log is reviewed
- **Then** the invitation event appears with the admin's ID, timestamp, and invited email

### TC-AUDIT-003: Change Request Lifecycle Events Are Logged
**Layer**: API | **Spec**: 020-audit-logs-siem

- **Given** a Change Request goes through Pending → Approved
- **When** the audit log is queried for the related flag
- **Then** both the "Created" and "Approved" events are recorded with respective actor IDs

### TC-AUDIT-004: Audit Log Is Immutable (Append-Only)
**Layer**: API | **Spec**: 020-audit-logs-siem

- **Given** an audit log entry exists
- **When** a direct API call attempts to modify or delete an audit log entry
- **Then** the API returns 405 Method Not Allowed or 403 Forbidden; the entry is unchanged

### TC-AUDIT-005: Audit Log — Export as CSV
**Layer**: API + UI | **Spec**: 020-audit-logs-siem

- **Given** audit logs contain 100+ entries
- **When** the user clicks "Export CSV"
- **Then** a CSV file is downloaded with all fields (timestamp, actor, action, previous state, new state, IP)

### TC-AUDIT-006: Audit Log — Filter by Date Range
**Layer**: UI | **Spec**: 020-audit-logs-siem

- **Given** audit logs span multiple days
- **When** the user filters by a specific date range
- **Then** only entries within that range are shown

### TC-AUDIT-007: Audit Log — SIEM Webhook Streaming
**Layer**: API | **Spec**: 020-audit-logs-siem

- **Given** a SIEM webhook endpoint is configured
- **When** a flag change occurs
- **Then** within seconds, the audit event is POSTed to the configured webhook URL in the defined JSON format

### TC-AUDIT-008: Audit Log Does Not Capture Plaintext API Keys
**Layer**: API | **Spec**: 021-pii-hashing-privacy, p-requirements §15.3

- **Given** flag evaluation or administrative actions occur
- **When** the audit logs are reviewed
- **Then** no plaintext SDK API keys or raw PII targeting values appear in any log entry

---

## Part 13 — SDK: Server-Side Local Evaluation

### TC-SDK-001: SDK Bootstraps with Full Ruleset Snapshot
**Layer**: SDK | **Spec**: 005-sdk-evaluation-api

- **Given** a Node.js SDK is initialized with a valid environment SDK key
- **When** the SDK connects to the FlagManagment server on startup
- **Then** the SDK downloads the complete flag ruleset for that environment into local memory within 5 seconds

### TC-SDK-002: SDK Evaluates Flag in Sub-Millisecond (Local)
**Layer**: SDK | **Spec**: 005-sdk-evaluation-api, p-requirements §7.1

- **Given** the SDK has a full snapshot in memory
- **When** 10,000 flag evaluations are performed sequentially
- **Then** each evaluation completes in under 1ms (no outbound network call per evaluation)

### TC-SDK-003: SDK Receives Delta Update After Flag Change
**Layer**: SDK | **Spec**: 005-sdk-evaluation-api

- **Given** an SDK is connected via streaming gRPC
- **When** an admin changes a flag state in the dashboard
- **Then** the SDK receives a delta update and starts returning the new flag value within 1 second

### TC-SDK-004: SDK Continues Evaluating on Connection Loss (Resilience)
**Layer**: SDK | **Spec**: 005-sdk-evaluation-api, p-requirements §7.3

- **Given** an SDK is connected and has a full snapshot
- **When** the network connection to the FlagManagment server drops
- **Then** the SDK continues evaluating all flags using the last known good snapshot; no errors are thrown to the application

### TC-SDK-005: SDK Reconnects with Exponential Backoff
**Layer**: SDK | **Spec**: 005-sdk-evaluation-api

- **Given** an SDK has lost its streaming connection
- **When** the server becomes available again
- **Then** the SDK reconnects automatically using exponential backoff and re-downloads the latest snapshot

### TC-SDK-006: SDK Returns Safe Default on Unknown Flag Key
**Layer**: SDK

- **Given** an SDK is connected
- **When** a flag key that does not exist in the snapshot is evaluated
- **Then** the SDK returns the user-specified default value (not an error), and logs a warning

### TC-SDK-007: Node.js SDK — Boolean Evaluation
**Layer**: SDK | **Spec**: 006-nodejs-sdk

- **Given** a boolean flag `feature-payments` is enabled
- **When** `client.getBooleanValue("feature-payments", false, context)` is called
- **Then** `true` is returned

### TC-SDK-008: Node.js SDK — String Evaluation (Remote Config)
**Layer**: SDK | **Spec**: 006-nodejs-sdk

- **Given** a flag has string remote config value `"v2"`
- **When** `client.getStringValue("api-version", "v1", context)` is called
- **Then** `"v2"` is returned

### TC-SDK-009: Node.js SDK — Multivariate Variant Evaluation
**Layer**: SDK | **Spec**: 006-nodejs-sdk, 013-multivariate-flags

- **Given** a multivariate flag has variant "beta" allocated to 30% of users
- **When** 1000 different user contexts evaluate the flag
- **Then** approximately 30% return "beta"

### TC-SDK-010: OpenFeature API Compatibility
**Layer**: SDK | **Spec**: 024-openfeature-compliance

- **Given** the FlagManagment Node.js SDK is initialized as an OpenFeature provider
- **When** standard OpenFeature `client.getBooleanDetails(...)` is called
- **Then** the response conforms to the OpenFeature `ResolutionDetails` shape (value, reason, variant, errorCode)

---

## Part 14 — SDK Event Forwarding (A/B Analytics)

### TC-EVENT-001: SDK Fires Evaluation Event on Flag Resolve
**Layer**: SDK | **Spec**: 017-sdk-event-forwarding

- **Given** an event forwarding hook is registered on the SDK
- **When** a flag is evaluated and a user is bucketed into a variant
- **Then** the evaluation event (flag key, variant assigned, user identity) is passed to the registered hook

### TC-EVENT-002: Event Forward to PostHog (Integration)
**Layer**: SDK | **Spec**: 017-sdk-event-forwarding

- **Given** a PostHog hook is configured on the SDK
- **When** a multivariate flag is evaluated for a user
- **Then** an event appears in PostHog with the flag key, variant, and user ID

### TC-EVENT-003: No Event Fired for Default Fallback
**Layer**: SDK | **Spec**: 017-sdk-event-forwarding

- **Given** an SDK evaluates a flag and returns the default value (flag not found or server offline)
- **When** the evaluation hook executes
- **Then** the event is marked with reason `DEFAULT` and optionally suppressed by hook configuration

---

## Part 15 — Stale Flag Detection

### TC-STALE-001: Flag Marked Stale After Inactivity
**Layer**: API + UI | **Spec**: 015-stale-flag-management

- **Given** a flag has been at 100% rollout with no state change for 31+ days
- **When** the dashboard loads the flag list
- **Then** the flag is visually highlighted with a "Stale" badge/warning indicator

### TC-STALE-002: `last_evaluated_at` Timestamp Updated on Evaluation
**Layer**: API + SDK | **Spec**: 015-stale-flag-management

- **Given** a flag has a `last_evaluated_at` field in the database
- **When** the SDK evaluates the flag
- **Then** the `last_evaluated_at` timestamp is updated to the current time

### TC-STALE-003: Stale Flag Alert Threshold is Configurable
**Layer**: API | **Spec**: 015-stale-flag-management

- **Given** the system has a configurable stale threshold (default 30 days)
- **When** the threshold is changed to 7 days
- **Then** flags inactive for 7+ days at 100% rollout are now highlighted as stale

### TC-STALE-004: Dismissing Stale Flag Warning
**Layer**: UI | **Spec**: 015-stale-flag-management

- **Given** a flag is flagged as stale
- **When** the user dismisses or archives the stale warning
- **Then** the warning is cleared; an audit log entry records the dismissal

---

## Part 16 — Telemetry & Automated Kill-Switches

### TC-TELEM-001: Configure Telemetry Trigger Rule
**Layer**: API + UI | **Spec**: 009-telemetry-kill-switches

- **Given** a flag is active in Production
- **When** the user creates a rule: "IF Datadog alert `high-error-rate` fires, SET flag to 0% rollout"
- **Then** the rule is saved and associated with the flag and environment

### TC-TELEM-002: Telemetry Webhook Received — Trigger Fires
**Layer**: API | **Spec**: 009-telemetry-kill-switches

- **Given** a kill-switch rule is configured for a flag
- **When** the FlagManagment API receives a webhook POST from Datadog matching the rule condition
- **Then** the flag rollout is automatically set to 0%; an audit log entry records the automated action and its trigger source

### TC-TELEM-003: Automated Kill-Switch Logged in Audit Log
**Layer**: API | **Spec**: 009-telemetry-kill-switches, 020-audit-logs-siem

- **Given** a kill-switch was automatically triggered
- **When** the audit log is reviewed
- **Then** the entry shows actor as "system/automation", the trigger condition, and the action taken (rollout set to 0%)

### TC-TELEM-004: Manual Override After Automated Kill-Switch
**Layer**: API + UI | **Spec**: 009-telemetry-kill-switches

- **Given** a kill-switch has set a flag to 0%
- **When** an engineer manually overrides the rollout to 25% in the dashboard
- **Then** the override takes effect immediately; the manual change is audit logged

---

## Part 17 — Admin Dashboard & UX

### TC-UI-001: Projects List Page Renders Correctly
**Layer**: UI | **Spec**: 004-frontend-dashboard

- **Given** multiple projects exist
- **When** the admin logs in
- **Then** the Projects list page shows all accessible projects with names and environment counts

### TC-UI-002: Create Project via UI
**Layer**: UI

- **Given** the admin is on the Projects page
- **When** they click "New Project" and fill in a name
- **Then** the project is created, a success notification is shown, and the project appears in the list

### TC-UI-003: Competitor UX Parity — Unleash-Style Flag Management
**Layer**: UI | **Spec**: 031-flagcompatitor-ux-parity

- **Given** the dashboard flags page is open
- **When** the user interacts with flag toggles, search, and pagination
- **Then** the UX is comparable or superior to Unleash's open-source UI (visual clarity, responsiveness, speed)

### TC-UI-004: Environment Switcher is Prominently Visible
**Layer**: UI | **Spec**: 030-env-context-switching

- **Given** a project with multiple environments is open
- **When** the user views any page (flags, targeting, etc.)
- **Then** the active environment is clearly shown in the header/nav; switching is a single click

### TC-UI-005: Empty State on First-Time Project Setup
**Layer**: UI | **Spec**: 028-empty-state-onboarding

- **Given** a new project has no environments or flags
- **When** the user opens the project
- **Then** a guided onboarding flow/empty state is shown with CTAs to create their first environment and flag

### TC-UI-006: Audit Log UI — Queryable List
**Layer**: UI | **Spec**: 020-audit-logs-siem

- **Given** audit logs exist
- **When** the user navigates to the Audit Log section
- **Then** logs are shown in reverse-chronological order with actor, action, environment, and timestamp; filtering by actor and action type is available

### TC-UI-007: Change Requests List in UI
**Layer**: UI | **Spec**: 008-change-requests-workflow

- **Given** pending Change Requests exist
- **When** the user navigates to Change Requests
- **Then** all pending CRs are listed with environment name, flag name, requester, and submission timestamp

---

## Part 18 — PII Hashing & Data Privacy

### TC-PII-001: User Identity Hashed Before Storage
**Layer**: API | **Spec**: 021-pii-hashing-privacy, p-requirements §15.3

- **Given** an evaluation event includes a user email as identity
- **When** the event is stored or logged
- **Then** the email is stored as a salted hash, not plaintext

### TC-PII-002: API Keys Are Hashed in Database
**Layer**: API | **Spec**: 021-pii-hashing-privacy

- **Given** an environment SDK key is created
- **When** the database `environments` table is inspected
- **Then** only the hashed form of the API key (`api_key_hash`) is present; no plaintext key is stored

### TC-PII-003: Audit Logs Contain No Sensitive Targeting Metadata
**Layer**: API | **Spec**: 021-pii-hashing-privacy, 020-audit-logs-siem

- **Given** a targeting rule uses email addresses as attribute values
- **When** the corresponding audit log entry is retrieved
- **Then** the `previous_state` and `new_state` fields reference attribute names/operators only; individual user PII is not present

---

## Part 19 — Slack & Webhook Notifications

### TC-SLACK-001: Configure Slack Webhook for Flag Changes
**Layer**: API + UI | **Spec**: 010-slack-notifications

- **Given** a Slack webhook URL is configured
- **When** a flag is toggled in a watched environment
- **Then** a Slack message is sent with the flag name, action, environment, and actor name

### TC-SLACK-002: Slack Notification for Change Request Created
**Layer**: API | **Spec**: 010-slack-notifications

- **Given** a Slack webhook is configured
- **When** a Change Request is created in Production
- **Then** a Slack notification is sent to the configured channel with a link to the Change Request

### TC-SLACK-003: Slack Notification for Change Request Approved/Rejected
**Layer**: API | **Spec**: 010-slack-notifications

- **Given** a Slack webhook is configured
- **When** a Change Request is approved or rejected
- **Then** a Slack notification includes the outcome, approver name, and the flag affected

### TC-SLACK-004: Slack Notifications Are Opt-In Per Environment
**Layer**: API + UI | **Spec**: 010-slack-notifications

- **Given** a webhook is configured only for Production
- **When** a flag is toggled in the Dev environment
- **Then** no Slack notification is sent for the Dev change

---

## Part 20 — Edge Proxy / Relay Node

### TC-EDGE-001: Edge Proxy Connects to FlagManagment Backend
**Layer**: Infrastructure | **Spec**: 019-edge-proxy-relay

- **Given** an Edge Proxy is deployed in a private subnet
- **When** it is initialized with the environment SDK key
- **Then** it establishes a persistent gRPC connection to the FlagManagment backend and downloads the snapshot

### TC-EDGE-002: Internal SDK Connects to Edge Proxy (No Direct Internet)
**Layer**: SDK + Infrastructure | **Spec**: 019-edge-proxy-relay

- **Given** internal microservice SDKs are configured to point to the Edge Proxy URL
- **When** the SDKs evaluate flags
- **Then** evaluations succeed without any direct outbound internet access from the microservice; all traffic flows through the proxy

### TC-EDGE-003: Edge Proxy Propagates Delta Updates
**Layer**: Infrastructure | **Spec**: 019-edge-proxy-relay

- **Given** an Edge Proxy is running with SDKs connected
- **When** an admin changes a flag
- **Then** the Edge Proxy receives the delta from the backend and forwards it to all connected internal SDKs

---

## Part 21 — Terraform Provider

### TC-TERRAFORM-001: Create Project via Terraform
**Layer**: IaC | **Spec**: 025-terraform-provider

- **Given** a Terraform configuration declares a `flagmanagment_project` resource
- **When** `terraform apply` is run
- **Then** the project is created in FlagManagment via API; the state file reflects the created resource

### TC-TERRAFORM-002: Create Environment via Terraform
**Layer**: IaC | **Spec**: 025-terraform-provider

- **Given** a Terraform config declares `flagmanagment_environment` resources under a project
- **When** `terraform apply` is run
- **Then** all declared environments are created with correct names and protection settings

### TC-TERRAFORM-003: Destroy Project via Terraform
**Layer**: IaC | **Spec**: 025-terraform-provider

- **Given** a Terraform-managed project exists
- **When** `terraform destroy` is run
- **Then** the project and all its environments are removed; the Terraform state is updated accordingly

### TC-TERRAFORM-004: Terraform Drift Detection
**Layer**: IaC | **Spec**: 025-terraform-provider

- **Given** a Terraform-managed environment has been modified outside of Terraform (via UI)
- **When** `terraform plan` is run
- **Then** the plan shows the drift and proposes corrective changes

---

## Part 22 — CI/CD Pipeline & DevOps

### TC-CICD-001: CI Pipeline Passes on Valid Go Code
**Layer**: CI/CD | **Spec**: 036-cicd-security-release-workflows

- **Given** code is pushed to a feature branch
- **When** the GitHub Actions CI workflow runs
- **Then** Go lint, Go unit tests, and frontend tests all pass; the workflow reports green

### TC-CICD-002: CI Fails on Lint Errors
**Layer**: CI/CD | **Spec**: 036-cicd-security-release-workflows

- **Given** code with a golangci-lint violation is pushed
- **When** the CI lint job runs
- **Then** the lint job fails with a descriptive error; the PR is blocked from merging

### TC-CICD-003: Docker Images Built Successfully
**Layer**: CI/CD | **Spec**: 036-cicd-security-release-workflows

- **Given** a commit is pushed to the main branch
- **When** the Docker build job runs
- **Then** backend and frontend Docker images are built successfully with no errors

### TC-CICD-004: Trivy Vulnerability Scan Runs on Images
**Layer**: CI/CD | **Spec**: 036-cicd-security-release-workflows

- **Given** Docker images are built
- **When** the Trivy scan job runs
- **Then** the vulnerability report is generated; CRITICAL vulnerabilities fail the build (or warn, per policy)

### TC-CICD-005: Ephemeral Test Environment Created in CI
**Layer**: CI/CD + API | **Spec**: 018-environment-cloning, 036-cicd-security-release-workflows

- **Given** a CI/CD workflow runs integration tests
- **When** the workflow calls the FlagManagment API to clone the QA environment for `PR-123`
- **Then** the ephemeral environment is created; tests run against it; the environment is deleted after tests complete

---

## Part 23 — Non-Functional & Performance

### TC-PERF-001: Single SDK Flag Evaluation Under 1ms
**Layer**: SDK | **Ref**: p-requirements §7.1, §15.2

- **Given** the SDK has a full snapshot in memory
- **When** a single flag evaluation is performed
- **Then** wall-clock time for the evaluation is under 1 millisecond on commodity hardware

### TC-PERF-002: Backend Handles 1000 Concurrent SDK Connections
**Layer**: Backend | **Ref**: p-requirements §7.1

- **Given** a single FlagManagment backend instance
- **When** 1000 SDK clients connect simultaneously via gRPC streaming
- **Then** all connections are maintained; memory usage is stable; flag delta updates are delivered to all clients within 2 seconds

### TC-PERF-003: API Response Time Under 200ms (P95)
**Layer**: API | **Ref**: p-requirements §7.1

- **Given** the backend is under normal load
- **When** REST API calls for flag CRUD operations are made
- **Then** P95 response time is under 200ms

### TC-PERF-004: Backend Unit Test Coverage >= 80%
**Layer**: Backend | **Ref**: p-requirements §15.1

- **Given** the backend test suite is run
- **When** coverage is measured
- **Then** total unit test coverage is >= 80%

### TC-PERF-005: Frontend Test Coverage >= 70%
**Layer**: Frontend | **Ref**: p-requirements §15.1

- **Given** the frontend test suite is run
- **When** coverage is measured
- **Then** total unit/component test coverage is >= 70%

### TC-SEC-001: All External Traffic Requires HTTPS/TLS
**Layer**: Infrastructure | **Ref**: p-requirements §7.4

- **Given** the platform is deployed with TLS configured
- **When** an HTTP (non-TLS) request is made to any external endpoint
- **Then** the request is redirected to HTTPS or rejected

### TC-SEC-002: SDK Keys Not Exposed in Client-Side Logs
**Layer**: SDK + API | **Ref**: p-requirements §7.4

- **Given** a server-side SDK key is used for evaluation
- **When** SDK logs are inspected at INFO/DEBUG level
- **Then** the full SDK key is masked or absent from log output

---

## Part 24 — End-to-End Scenarios (Cross-Feature Flows)

### TC-E2E-001: Full Feature Rollout Pipeline (Dev to Production)
**Layer**: API + SDK + UI

1. Create a project with Dev, QA, Staging, Production environments (Production protected)
2. Create a boolean flag `new-payments-v2` in Dev, enable it
3. Verify SDK in Dev returns `true`
4. Promote flag from Dev to QA; verify SDK in QA returns `true`, Dev unchanged
5. Promote QA to Staging; verify Staging SDK returns `true`
6. Promote Staging to Production; verify a Change Request is created
7. Release Manager approves the Change Request
8. Verify Production SDK returns `true`; audit log shows full lifecycle

### TC-E2E-002: A/B Test with Analytics Event Forwarding
**Layer**: SDK + Analytics

1. Create a multivariate flag with variants Control (60%) and Variant-A (40%)
2. Configure SDK event forwarding to PostHog
3. Run 1000 evaluations with unique user IDs
4. Verify approximately 40% bucket into Variant-A
5. Verify PostHog receives evaluation events with variant and user ID
6. Verify deterministic bucketing: same user always gets same variant

### TC-E2E-003: Automated Rollback via Kill-Switch
**Layer**: API + Telemetry + Audit

1. Enable a flag in Production at 50% rollout
2. Configure a kill-switch rule: "IF error-rate-alert fires, set rollout to 0%"
3. Simulate a Datadog webhook firing the alert
4. Verify the flag rollout is automatically set to 0%
5. Verify the audit log records the automated action with trigger source
6. Manually re-enable the flag to 25% rollout and verify it takes effect

### TC-E2E-004: Change Request Full Cycle in Protected Environment
**Layer**: API + UI

1. Create Production as a protected environment
2. Assign QA Engineer and Release Manager roles to two users
3. QA Engineer toggles a flag in Production (via UI/API)
4. Verify a Change Request is created (Pending); flag is unchanged
5. Release Manager views the diff (current vs proposed)
6. Release Manager approves the Change Request
7. Verify the flag state changes; audit log records approval with approver ID
8. Repeat with rejection scenario; verify no change is applied

### TC-E2E-005: Stale Flag Lifecycle Management
**Layer**: API + UI

1. Create and fully enable a flag at 100% rollout
2. Simulate 31 days of inactivity (no state change, no evaluation)
3. Open the dashboard; verify the flag is highlighted as "Stale"
4. Dismiss the stale warning
5. Verify the dismissal is recorded in the audit log

### TC-E2E-006: Ephemeral Environment for PR Integration Test
**Layer**: API + CI/CD

1. CI/CD triggers: clone QA environment to `PR-789-test` via API
2. Verify the cloned environment has all flags from QA with a new SDK key
3. Run integration tests against `PR-789-test`
4. Delete `PR-789-test` via API
5. Verify the environment no longer exists; QA environment is unchanged

---

## Requirements

### Functional Requirements

- **FR-TEST-001**: Test cases MUST cover all flag types: Boolean, Multivariate, and Remote Configuration
- **FR-TEST-002**: Test cases MUST cover all RBAC roles at global, project, and environment levels
- **FR-TEST-003**: Test cases MUST include both happy-path and failure/negative scenarios
- **FR-TEST-004**: End-to-end tests MUST demonstrate a complete Dev to Production promotion cycle
- **FR-TEST-005**: SDK tests MUST validate sub-millisecond local evaluation, delta updates, and resilience
- **FR-TEST-006**: Audit log tests MUST confirm append-only immutability and no PII exposure
- **FR-TEST-007**: Each test case MUST be independently executable in isolation
- **FR-TEST-008**: Test cases MUST be executable by both human testers (manual) and AI agents (automated)

### Non-Functional Requirements

- **NFR-TEST-001**: The test plan MUST be stored in a location accessible to CI/CD automation
- **NFR-TEST-002**: Test IDs MUST be unique and stable (not renumbered between iterations)
- **NFR-TEST-003**: Test cases MUST reference their source spec and requirements document section

---

## Success Criteria

1. **Complete Coverage**: Every functional feature in g-requirements.md and p-requirements.md has at least one positive and one negative test case
2. **SDK Latency**: SDK local evaluation tests confirm sub-1ms performance under load
3. **RBAC Enforcement**: All 7 RBAC roles can be verified with no unauthorized access passing any test
4. **End-to-End Flows**: All 6 E2E scenario tests execute successfully with real system state verification
5. **Audit Completeness**: Every data mutation has a corresponding verifiable audit log entry
6. **Human Usability**: A QA engineer can understand and execute any test case without additional documentation
7. **AI Parsability**: Each test case follows a strict Given/When/Then format parseable by test automation agents

---

## Assumptions

- The SDK tests target the Node.js SDK (spec 006) as the reference implementation; patterns apply to all SDKs (spec 023)
- "Admin" in test cases refers to a System Administrator role unless stated otherwise
- Performance tests are run on a standard cloud VM (2 vCPU, 4GB RAM)
- Integration with external tools (PostHog, Datadog, Slack) uses mock/test endpoints during automated testing
- Terraform tests require a running FlagManagment instance and valid API credentials
- Stale flag "31 days inactivity" is simulated by directly manipulating the `last_evaluated_at` timestamp in the database

---

## Dependencies

| Dependency | Required For |
|---|---|
| Backend API running locally or in CI | All API and SDK tests |
| Frontend app running (`npm run dev`) | All UI tests |
| PostgreSQL and Redis running | All API and SDK tests |
| Node.js SDK library | SDK test cases |
| Mock Datadog/Slack webhook endpoint | Telemetry and Slack tests |
| Terraform CLI + provider binary | Terraform IaC tests |
| Docker + docker-compose | CI/CD and edge proxy tests |
