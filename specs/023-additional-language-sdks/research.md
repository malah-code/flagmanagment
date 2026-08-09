# Research: Additional Language SDKs

## 1. Cross-Language Bucketing (MurmurHash3)
**Decision**: Standardize on MurmurHash3 32-bit x86 variant across all languages.
**Rationale**: Ensuring bucketing consistency across Go, Java, Python, .NET, Node, React, iOS, and Android requires a deterministic hashing algorithm. MurmurHash3 is fast, non-cryptographic, and widely available in all ecosystems.
**Alternatives considered**: SHA-256 (too slow for <1ms requirement), FNV-1a (higher collision rate for user IDs).

## 2. Server-Side Sync Mechanisms (Go, Java, Python, .NET)
**Decision**: Use Server-Sent Events (SSE) / HTTP Streaming as the primary synchronization channel, falling back to HTTP polling.
**Rationale**: SSE is natively supported by HTTP/1.1 and HTTP/2 clients in all major languages without requiring heavy gRPC libraries. It matches the requirements for lightweight delta updates and avoids complex gRPC compilation toolchains for every language.
**Alternatives considered**: gRPC bidir-streaming (complex to distribute as library dependencies across all languages), WebSockets (requires persistent stateful load balancing which complicates edge proxying).

## 3. OpenFeature Compliance
**Decision**: Implement the `Provider` interface specified by OpenFeature for each language ecosystem.
**Rationale**: Constitution explicitly requires OpenFeature interoperability. Each language has a standard OpenFeature SDK (e.g., `github.com/open-feature/go-sdk`, `@openfeature/react-sdk`).
**Alternatives considered**: Proprietary client interfaces (violates Constitution).

## 4. Mobile Offline Storage (iOS, Android)
**Decision**: Use platform-native secure enclaves for caching flags (Keychain for iOS, EncryptedSharedPreferences for Android).
**Rationale**: Prevents malicious actors with physical device access from altering feature flags (e.g., bypassing paywalls or unlocking features).
**Alternatives considered**: SQLite/CoreData without encryption (security risk), SharedPreferences/UserDefaults without encryption (security risk).

## 5. React Hooks & Re-renders
**Decision**: Use React 18+ `useSyncExternalStore` for robust state management.
**Rationale**: Avoids tearing in React concurrent mode and provides a seamless mechanism to subscribe to the internal FlagManagment SDK state without wrapping the entire app in complex context providers that trigger full-tree re-renders.
**Alternatives considered**: `useState` + `useEffect` (prone to stale closures and tearing).
