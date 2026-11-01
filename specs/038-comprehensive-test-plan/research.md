# Research: Automated Test Execution Framework & Harness

**Branch**: `038-comprehensive-test-plan` | **Date**: 2026-10-31

This research document consolidates architectural and tooling decisions to implement the automated execution and tracking of the 158 test cases defined in [spec.md](spec.md).

---

## Decision 1: Unified Test Orchestration & Reporting

- **Decision**: Implement a root-level test harness script [`scripts/run-all-tests.sh`](../../scripts/run-all-tests.sh) supporting scoped execution (`--backend`, `--frontend`, `--e2e`, `--sdk`, `--perf`, or `--all`) that generates structured test reports (JUnit XML / JSON) and consolidates exit codes.
- **Rationale**: Developers and CI pipelines need a single deterministic entry point to run any tier or the complete suite. Fragmented commands across languages increase maintenance burden and lead to overlooked test failures.
- **Alternatives Considered**:
  - *Make targets only (`make test-all`)*: Retained as convenience wrappers, but shell scripts provide superior argument parsing, environment validation, and selective execution.
  - *Monorepo tools (Turborepo/Nx)*: Adds unnecessary heavy Node.js dependencies across Go and Rust/C SDK sub-projects.

---

## Decision 2: Backend Integration Test Isolation

- **Decision**: Use a dedicated PostgreSQL test database (`flagmanagment_test`) and Redis test DB (`15`). Each test package or suite initializes with clean migrations (`golang-migrate`), and individual test functions utilize unique project/environment UUIDs and cleanups.
- **Rationale**: In-memory SQLite or mock drivers fail to validate PostgreSQL-specific JSONB queries, MurmurHash3 distribution, and Redis Pub/Sub notification semantics. Real ephemeral containers guarantee true behavioral fidelity.
- **Alternatives Considered**:
  - *Testcontainers-Go*: Good for isolated spins, but slower in local workflows where developers already have `docker-compose.yml` running.
  - *SQLite in-memory*: Incompatible with PostgreSQL JSONB operators (`@>`, `->>`) used heavily in flag targeting logic.

---

## Decision 3: Mocking External Service Dependencies (OIDC & Webhooks)

- **Decision**: Implement hermetic test doubles using Go's `net/http/httptest`:
  1. **Mock OIDC Provider (`testutil/oidc_mock.go`)**: Serves OpenID configuration discovery (`/.well-known/openid-configuration`), JWKS key set, and returns valid signed JWTs with customizable claims to test SSO role mappings (TC-AUTH-009, TC-AUTH-010).
  2. **Mock Webhook Receiver (`testutil/webhook_mock.go`)**: Captures HTTP POST payloads for Slack notifications and SIEM audit exports (TC-SLACK-001–004, TC-AUDIT-004) without internet access.
- **Rationale**: Eliminates network flakiness, credential leaks, and third-party rate limits in CI.

---

## Decision 4: Sub-Millisecond Benchmark Verification (Constitution Principle IV)

- **Decision**: Integrate automated latency assertions into Go and Node.js SDK test suites using standard Go benchmarks (`testing.B`) and high-resolution timers (`hrtime.bigint()`).
  - The benchmark suite evaluates a flag with 10 targeting rules 100,000 times consecutively in-memory.
  - Test fails if p95 latency exceeds 1.0 millisecond.
- **Rationale**: Constitution Principle IV explicitly mandates: *"Server-side SDK flag evaluation MUST execute in under 1 millisecond on commodity hardware... This constraint MUST be verified via load testing before any SDK feature is marked complete."*

---

## Decision 5: Test Case Mapping & Automated Tracker Synchronization

- **Decision**: Adopt a standardized code annotation convention across all test files:
  ```go
  // @TestCase: TC-AUTH-001
  func TestLoginWithValidCredentials(t *testing.T) { ... }
  ```
  ```typescript
  test('TC-AUTH-001: Login with valid credentials', async ({ page }) => { ... });
  ```
  A Python synchronizer script (`scripts/sync-test-tracker.py`) parses these annotations, cross-references test results against the 158 canonical test IDs in `spec.md`, and updates [`docs/TEST_EXECUTION_TRACKER.md`](../../docs/TEST_EXECUTION_TRACKER.md) with pass/fail badges.
- **Rationale**: Keeps documentation 100% synchronized with actual code and provides instant visibility into which of the 158 test cases are covered and passing.

