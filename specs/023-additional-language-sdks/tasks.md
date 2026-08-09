# Tasks: Additional Language SDKs (Go, Java, Python, .NET, React, iOS, Android)

**Input**: Design documents from `/specs/023-additional-language-sdks/`

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure for the multi-language monorepo.

- [x] T001 Create multi-language SDK directory structure (`sdk/go`, `sdk/java`, `sdk/python`, `sdk/dotnet`, `sdk/react`, `sdk/ios`, `sdk/android`)
- [x] T002 [P] Initialize Go SDK module in `sdk/go/go.mod`
- [x] T003 [P] Initialize Java SDK project in `sdk/java/pom.xml` or `build.gradle`
- [x] T004 [P] Initialize Python SDK project in `sdk/python/pyproject.toml`
- [x] T005 [P] Initialize .NET SDK project in `sdk/dotnet/FlagManagment.Sdk.csproj`
- [x] T006 [P] Initialize React SDK project in `sdk/react/package.json`
- [x] T007 [P] Initialize iOS SDK package in `sdk/ios/Package.swift`
- [x] T008 [P] Initialize Android SDK project in `sdk/android/build.gradle`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented.

- [x] T009 Implement generic MurmurHash3 hashing utility in backend (if missing) and verify consistency tests in `backend/internal/utils/hash.go`
- [x] T010 Implement backend SSE streaming endpoint `/api/v1/sdk/stream` in `backend/internal/api/sdk_stream.go`

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Server-Side Language SDKs (Go, Java, Python, .NET) (Priority: P1) 🎯 MVP

**Goal**: Deliver high-performance, native SDKs that evaluate feature flags locally in-memory with sub-millisecond latency and sync via streaming.

**Independent Test**: Can be independently verified by initializing any of the server SDKs with an environment API key, evaluating flags locally under concurrency, and verifying sub-millisecond evaluation times and offline resilience.

### Implementation for User Story 1

- [x] T011 [P] [US1] Go SDK: Create OpenFeature provider interface in `sdk/go/provider.go`
- [x] T012 [P] [US1] Go SDK: Implement SSE streaming client in `sdk/go/client.go`
- [x] T013 [P] [US1] Go SDK: Implement in-memory evaluation engine with MurmurHash3 in `sdk/go/evaluator.go`
- [x] T014 [P] [US1] Java SDK: Create OpenFeature provider interface in `sdk/java/src/main/java/com/flagmanagment/sdk/Provider.java`
- [x] T015 [P] [US1] Java SDK: Implement SSE streaming client in `sdk/java/src/main/java/com/flagmanagment/sdk/Client.java`
- [x] T016 [P] [US1] Java SDK: Implement in-memory evaluation engine with MurmurHash3 in `sdk/java/src/main/java/com/flagmanagment/sdk/Evaluator.java`
- [x] T017 [P] [US1] Python SDK: Create OpenFeature provider interface in `sdk/python/flagmanagment/provider.py`
- [x] T018 [P] [US1] Python SDK: Implement SSE streaming client in `sdk/python/flagmanagment/client.py`
- [x] T019 [P] [US1] Python SDK: Implement in-memory evaluation engine with MurmurHash3 in `sdk/python/flagmanagment/evaluator.py`
- [x] T020 [P] [US1] .NET SDK: Create OpenFeature provider interface in `sdk/dotnet/Provider.cs`
- [x] T021 [P] [US1] .NET SDK: Implement SSE streaming client in `sdk/dotnet/Client.cs`
- [x] T022 [P] [US1] .NET SDK: Implement in-memory evaluation engine with MurmurHash3 in `sdk/dotnet/Evaluator.cs`
- [x] T023 [P] [US1] Implement cross-SDK MurmurHash3 deterministic bucketing test suite for Go, Java, Python, and .NET.

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently

---

## Phase 4: User Story 2 - React Client-Side Web SDK (Priority: P1)

**Goal**: Deliver a declarative React SDK (Hooks and Provider) that manages feature flag states, context updates, and UI re-renders.

**Independent Test**: Wrap a React app with `<FlagProvider>`, use `useFlag('my-flag')` in child components, update flag values, and confirm immediate UI re-rendering.

### Implementation for User Story 2

- [x] T024 [P] [US2] React SDK: Implement `<FlagProvider>` wrapper in `sdk/react/src/FlagProvider.tsx`
- [x] T025 [P] [US2] React SDK: Implement `useFlag` and `useFlags` hooks using `useSyncExternalStore` in `sdk/react/src/hooks.ts`
- [x] T026 [P] [US2] React SDK: Implement SSE streaming client for React context updates in `sdk/react/src/client.ts`
- [x] T027 [US2] React SDK: Add OpenFeature web provider compliance in `sdk/react/src/provider.ts`

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently

---

## Phase 5: User Story 3 - Native Mobile SDKs (iOS & Android) (Priority: P2)

**Goal**: Deliver a mobile-optimized SDK that handles offline caching, battery/network-conscious rule fetching, and secure local storage.

**Independent Test**: Run iOS and Android sample apps, toggle flight mode while evaluating flags, update contexts on background/foreground transitions, and verify persistent flag state.

### Implementation for User Story 3

- [x] T028 [P] [US3] iOS SDK: Create OpenFeature provider interface in `sdk/ios/Sources/FlagManagment/Provider.swift`
- [x] T029 [P] [US3] iOS SDK: Implement SSE streaming and lifecycle-aware syncing in `sdk/ios/Sources/FlagManagment/Client.swift`
- [x] T030 [P] [US3] iOS SDK: Implement secure Keychain storage and offline evaluation cache in `sdk/ios/Sources/FlagManagment/Storage.swift`
- [x] T031 [P] [US3] Android SDK: Create OpenFeature provider interface in `sdk/android/src/main/kotlin/com/flagmanagment/sdk/Provider.kt`
- [x] T032 [P] [US3] Android SDK: Implement SSE streaming and lifecycle-aware syncing in `sdk/android/src/main/kotlin/com/flagmanagment/sdk/Client.kt`
- [x] T033 [P] [US3] Android SDK: Implement EncryptedSharedPreferences storage and offline cache in `sdk/android/src/main/kotlin/com/flagmanagment/sdk/Storage.kt`

**Checkpoint**: All user stories should now be independently functional

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [x] T034 [P] Document SDK usage, configuration options, and integration examples in `docs/sdk-reference.md`
- [x] T035 [P] Set up automated CI/CD publishing pipelines for all SDKs (npm, pip, maven, nuget, swift, gradle)
- [x] T036 Run `quickstart.md` validation tests against all implemented SDKs

---

## Dependencies & Execution Order

### Phase Dependencies
- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
  - User stories can then proceed in parallel
- **Polish (Final Phase)**: Depends on all desired user stories being complete

### User Story Dependencies
- **User Story 1 (P1)**: Can start after Foundational (Phase 2)
- **User Story 2 (P1)**: Can start after Foundational (Phase 2)
- **User Story 3 (P2)**: Can start after Foundational (Phase 2)

### Parallel Opportunities
- All language SDK setups in Phase 1 can be done in parallel.
- Go, Java, Python, and .NET SDK implementations in US1 can be built simultaneously by different developers.
- US1, US2, and US3 can be executed entirely in parallel since they involve separate codebases (`sdk/go`, `sdk/react`, `sdk/ios`, etc.) after the backend SSE endpoint (T010) is deployed.
