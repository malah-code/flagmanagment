# Data Model: Additional Language SDKs

## Key Entities

The internal SDK representations are unified across all supported languages (Go, Java, Python, .NET, React, iOS, Android).

### 1. SDK Configuration (Client Configuration)
Object initialized by the application to configure the SDK.
- `APIKey` (string, required): Environment-specific API key (`fm_sa_*` or standard).
- `BaseURL` (string, required): Management or Edge Proxy URL.
- `SyncMode` (enum: STREAMING | POLLING, default: STREAMING): How the SDK receives updates.
- `PollingIntervalSeconds` (int, default: 30): Active only if `SyncMode` == POLLING.
- `ConnectTimeout` (int, default: 5s).
- `OfflineMode` (boolean, default: false): If true, do not attempt to contact network.

### 2. Evaluation Context (OpenFeature standard)
Used for targeting rules and percentages.
- `TargetingKey` (string, required): Unique identifier for the user or entity.
- `Attributes` (map[string]any): Arbitrary metadata (e.g., email, country, app_version).

### 3. Flag Ruleset Snapshot (Internal Data Structure)
The JSON payload representing the environment's current flags.
- `Version` (int): Ruleset monotonic version.
- `Flags` (map[string]FlagDefinition): Dictionary of all flags in the environment.
- `Segments` (map[string]SegmentDefinition): Reusable targeting segments.

### 4. Flag Definition
- `Key` (string): Flag identifier.
- `Type` (enum: BOOLEAN, STRING, NUMBER, OBJECT).
- `Enabled` (boolean): Master kill switch.
- `DefaultVariant` (string).
- `TargetingRules` (array of TargetingRule).

### 5. Evaluation Result
The output of a flag evaluation request.
- `FlagKey` (string).
- `Value` (any): The resolved value (bool, string, int, double, json).
- `Variant` (string): Which variant was selected.
- `Reason` (enum: TARGETING_MATCH, DEFAULT, ERROR, OFF, SPLIT).
- `ErrorCode` (string, optional).
