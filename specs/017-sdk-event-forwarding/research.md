# Research: SDK Event Forwarding for Analytics

## Topic: Hook / Interceptor Interface Design

### Context
We need to design a lightweight interface for SDKs to broadcast when a flag evaluation occurs, so developers can wire these events to product analytics like Amplitude, PostHog, or Datadog without adding latency to the local evaluation.

### Decision
Implement the standard **OpenFeature Provider Hooks** pattern. OpenFeature already defines a robust `Hook` specification, including `after` and `error` stages, which fits perfectly with our goal of non-blocking evaluation interceptors.

### Rationale
- Conforms to OpenFeature standard (already a PRD requirement).
- Developers are already familiar with it.
- OpenFeature provides an `after` hook which fires right before the evaluation result is returned to the caller. This allows us to inspect the variant assigned.
- Async execution can be easily implemented inside the hook's `after` method by spawning a goroutine (in Go) or using Promises (in Node.js/JS) without blocking the main evaluation flow.

### Alternatives considered
- Creating a proprietary `EventForwarder` interface: Rejected because it violates the OpenFeature standardization requirement.
- Sending events to the FlagManagment backend for ingestion: Rejected because it adds latency, requires scaling the backend for event ingestion, and contradicts the assumption that forwarding happens locally in the SDK.
