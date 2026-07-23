# Research: Node.js SDK

## Topic 1: OpenFeature Compliance in Node.js

**Decision**: Build the `FlagManagmentProvider` implementing `@openfeature/server-sdk`'s `Provider` interface.

**Rationale**: The OpenFeature specification defines a standard `Provider` interface that SDK developers must implement to plug their proprietary evaluation logic into the OpenFeature API. By exposing this provider natively, users can do `OpenFeature.setProvider(new FlagManagmentProvider(client))` and instantly use standard OpenFeature evaluation calls.

**Alternatives considered**: 
- *Building our own proprietary wrapper*: Rejected as it violates Constitution VI (OpenFeature Interoperability).

## Topic 2: Deterministic Hashing for Rollouts

**Decision**: Use `murmurhash3js-revisited` or implement a lightweight pure-JS MurmurHash3 x86 32-bit algorithm internally.

**Rationale**: The Go backend evaluator (`backend/internal/sdk/evaluator.go`) explicitly uses a 32-bit MurmurHash3 algorithm on the concatenated string `flagKey + userIdentity` to modulo into a 1-10000 bucket. For rollouts to behave exactly the same on the server and the Node.js SDK, the hashing logic must be mathematically identical.

**Alternatives considered**:
- *SHA-256*: Not viable because it does not match the backend's MurmurHash3 bucket logic.
- *Node's native `crypto` module*: Does not ship with MurmurHash3 out-of-the-box. We must rely on an NPM package or inline the 32-bit Murmur3 function directly (which is <50 lines of code). Inlining is preferred to minimize dependencies.

## Topic 3: Real-Time Delta Updates (Streaming)

**Decision**: Use Server-Sent Events (SSE) via the standard `EventSource` polyfill (or native Node fetch streams) to consume the `StreamRulesets` backend equivalent.

**Rationale**: Since the backend `StreamRulesets` uses gRPC server-streaming, consuming gRPC directly in Node.js requires `@grpc/grpc-js`. Since we already built a REST wrapper, we can either build a gRPC client, or add an SSE wrapper in the Go backend. To minimize Node.js dependencies, we will use `@grpc/grpc-js` to directly connect to the backend `StreamRulesets` gRPC endpoint since the Go backend exposes a pure gRPC server on port 9090.

**Alternatives considered**:
- *SSE / REST Streaming*: Would require modifying the backend to expose an SSE REST endpoint. Since gRPC is already built and working (Feature 005), `@grpc/grpc-js` is the most direct path.
