# gRPC Interface Contract: Edge Proxy

**Feature**: `019-edge-proxy-relay`

---

## Overview

The Edge Proxy implements the **exact same** `SDKService` gRPC interface as the FlagManagment backend. SDK clients require **no code changes** — only the endpoint address changes from the backend to the proxy.

The proto source lives at: `backend/api/proto/sdk/v1/sdk.proto`
The generated Go code lives at: `backend/pkg/gen/sdk/v1/`

---

## Service Definition

```proto
// source: api/proto/sdk/v1/sdk.proto
service SDKService {
  // Server-streaming RPC: client subscribes to real-time delta updates
  rpc StreamRulesets(StreamRequest) returns (stream RulesetDelta);

  // Unary RPC: client fetches the complete current ruleset snapshot
  rpc FetchSnapshot(SnapshotRequest) returns (RulesetSnapshot);
}
```

---

## Proxy Behaviour vs. Backend Behaviour

| Behaviour | FlagManagment Backend | Edge Proxy |
|-----------|----------------------|------------|
| `FetchSnapshot` source | PostgreSQL + Redis cache | In-memory `RulesetStore` |
| `StreamRulesets` source | Redis pub/sub | In-memory broadcaster channel |
| Authentication | SDK env token validated against DB | Same SDK env token accepted (proxy validates it is the same token it was configured with) |
| Upstream | N/A | Connects to backend `SDKService` as a client |
| Delta push | Redis pub/sub → all connected SDKs | Upstream stream → in-memory broadcaster → all connected SDKs |

---

## Authentication

The Edge Proxy authenticates downstream SDK clients using the same `authorization` gRPC metadata header mechanism as the backend. The proxy validates that the environment token in the client request matches `FM_PROXY_ENV_TOKEN`. A mismatched token returns `codes.Unauthenticated`.

```
Metadata: authorization: Bearer env_<token>
```

---

## SDK Client Configuration Change

To redirect an SDK client to use the Edge Proxy, change only the endpoint:

| Config | Direct to Backend | Via Edge Proxy |
|--------|------------------|----------------|
| gRPC Endpoint | `flagmanagment.example.com:9090` | `edge-proxy.internal:9091` |
| Environment Token | unchanged | unchanged |
| Any SDK code | unchanged | unchanged |

---

## Error Codes

| Code | Condition |
|------|-----------|
| `UNAUTHENTICATED` | Token in request does not match `FM_PROXY_ENV_TOKEN` |
| `UNAVAILABLE` | Proxy is in `starting_up` state (no ruleset loaded yet) |
| `INTERNAL` | Unexpected internal error in proxy |
