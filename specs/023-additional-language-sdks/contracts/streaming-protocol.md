# SDK Streaming Protocol Contract

## Overview
SDKs connect to the Proxy or Backend via Server-Sent Events (SSE). 

## 1. Initial Handshake & Bootstrap
**Endpoint**: `GET /api/v1/sdk/stream`
**Headers**:
- `Authorization`: `Bearer fm_sa_...`
- `Accept`: `text/event-stream`

### Initial Snapshot Payload (Event: `bootstrap`)
The server immediately sends the full environment ruleset upon connection.

```json
event: bootstrap
data: {
  "version": 42,
  "flags": {
    "feature-x": {
      "key": "feature-x",
      "enabled": true,
      "type": "BOOLEAN",
      "defaultVariant": "on",
      "variants": { "on": true, "off": false },
      "rules": []
    }
  },
  "segments": {}
}
```

## 2. Delta Updates
When a flag is changed on the backend, a delta is pushed to active streams.

### Flag Updated Payload (Event: `flag_updated`)
```json
event: flag_updated
data: {
  "version": 43,
  "flagKey": "feature-x",
  "flag": { /* updated flag definition */ }
}
```

### Flag Deleted Payload (Event: `flag_deleted`)
```json
event: flag_deleted
data: {
  "version": 44,
  "flagKey": "feature-x"
}
```

## 3. Keep-Alive / Heartbeat
Sent every 30 seconds to prevent connection drops.
```json
event: ping
data: {"time": 1690000000}
```
