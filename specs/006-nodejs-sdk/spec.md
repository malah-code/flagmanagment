# Feature Specification: Node.js / TypeScript SDK

**Feature Branch**: `[006-nodejs-sdk]`
**Created**: 2026-07-22
**Status**: Draft
**Input**: User description: "Node.js / TypeScript SDK (matches frontend language)"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Initialize the SDK Client (Priority: P1)

As a Node.js developer, I want to initialize the FlagManagment SDK using an environment token, so that I can establish a secure connection to the backend and retrieve the initial ruleset snapshot.

**Why this priority**: Without initialization and the initial ruleset snapshot, no evaluation can happen.

**Independent Test**: Can be tested by instantiating the client and verifying it pulls the initial ruleset from the REST/gRPC endpoint.

**Acceptance Scenarios**:

1. **Given** a valid environment token, **When** `client.init()` is called, **Then** it fetches the initial snapshot and resolves a ready promise.
2. **Given** an invalid token, **When** `client.init()` is called, **Then** the promise rejects with an authentication error.

---

### User Story 2 - Local Flag Evaluation (Priority: P1)

As a developer, I want to evaluate feature flags locally in memory without network latency, so that my application remains highly performant (under 1ms).

**Why this priority**: Required by the Constitution for enterprise scalability.

**Independent Test**: Test evaluation using a pre-populated in-memory store.

**Acceptance Scenarios**:

1. **Given** a boolean flag enabled with no targeting rules, **When** evaluating the flag, **Then** it instantly returns true.
2. **Given** a multivariate flag with percentage rollouts, **When** evaluating with different user context identities, **Then** it consistently returns the correct variation via MurmurHash3 bucketing.

---

### User Story 3 - Real-time Delta Updates (Priority: P2)

As a developer, I want the SDK to automatically receive rule updates over a streaming connection, so that flag changes take effect instantly without needing to poll.

**Why this priority**: Delta updates provide a seamless, real-time experience while minimizing backend load.

**Independent Test**: Verify that sending an update to the streaming endpoint updates the local cache automatically.

**Acceptance Scenarios**:

1. **Given** a running SDK client, **When** a ruleset delta is broadcasted by the server, **Then** the SDK updates its in-memory ruleset and evaluates future requests using the new state.

---

### User Story 4 - OpenFeature Provider Compliance (Priority: P2)

As an architect using OpenFeature, I want to use the FlagManagment SDK as an OpenFeature Provider, so that I avoid vendor lock-in.

**Why this priority**: Mandated by the Constitution (VI. OpenFeature Interoperability).

**Independent Test**: Verify the SDK can be plugged into `@openfeature/server-sdk` and evaluate flags successfully.

**Acceptance Scenarios**:

1. **Given** the OpenFeature API, **When** the FlagManagment provider is registered, **Then** OpenFeature `getClient().getBooleanValue()` works seamlessly.

### Edge Cases

- What happens when the initial network request times out? (Should fall back to a local default value and continuously retry)
- How does the system handle an interruption in the streaming connection? (Should transparently reconnect with exponential backoff)

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The SDK MUST be written in TypeScript and compiled to CommonJS and ESM.
- **FR-002**: The SDK MUST fetch the complete ruleset snapshot over HTTP/gRPC on initialization.
- **FR-003**: The SDK MUST evaluate flags strictly in-memory without making outbound network requests during `evaluateFlag`.
- **FR-004**: The SDK MUST subscribe to the server's streaming endpoint to receive delta updates.
- **FR-005**: The SDK MUST implement deterministic identity hashing using MurmurHash3 for rollout percentages.
- **FR-006**: The SDK MUST include an OpenFeature-compliant Provider wrapper.

### Key Entities

- **FlagClient**: The main entry point for developers.
- **RuleStore**: The in-memory data structure holding the latest flag definitions.
- **Evaluator**: The pure-function logic engine that evaluates rules against context.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Flag evaluation completes in under 1 millisecond.
- **SC-002**: SDK reconnects automatically within 5 seconds if the streaming connection drops.
- **SC-003**: SDK achieves 100% compliance with the OpenFeature Server SDK provider test suite.

## Assumptions

- Node.js version 20+ is the target environment.
- Evaluation engine logic (MurmurHash3, targeting rules) precisely mirrors the server-side logic we already built in `internal/sdk/evaluator.go`.
