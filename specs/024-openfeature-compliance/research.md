# Phase 0: Outline & Research

## Context and Findings

The "OpenFeature API Compliance" feature aims to ensure that all 7 language SDKs in FlagManagment correctly implement the OpenFeature `Provider` API, specifically parsing the `EvaluationContext` to determine the `targetingKey`, applying MurmurHash3 rollout bucketing, and returning standardized OpenFeature `Reason` enums.

During a code audit of the current SDK implementations, we discovered:

1. **Go, Java, Python, and .NET**: These SDKs were recently updated during the "Additional Language SDKs" feature. Their OpenFeature `Provider` implementations are correctly wired to their respective local evaluators, supporting full context-aware targeting, deterministic bucketing, and returning appropriate OpenFeature `ResolutionDetails`.

2. **React SDK (`sdk/react/src/provider.ts`)**: The current `FlagManagmentWebProvider` is a stub. It completely ignores the `EvaluationContext` (meaning targeting rules are ignored via this interface) and always returns the default variant from the raw flag object. The local evaluation logic exists in `hooks.ts`, but it needs to be abstracted into the provider.

3. **iOS SDK (`sdk/ios/Sources/FlagManagment/Provider.swift`)**: The `FlagManagmentProvider` is a stub. It ignores `EvaluationContext`, ignores the `enabled` state of the flag, and forcefully attempts to cast and return the default variant.

4. **Android SDK (`sdk/android/src/main/kotlin/com/flagmanagment/sdk/Provider.kt`)**: The `FlagManagmentProvider` is a stub. It behaves exactly like the iOS provider, ignoring context and rollout logic.

## Decisions

- **Decision 1**: We will extract the `evaluateLocally` logic from the React SDK's `hooks.ts` into a standalone `evaluator.ts` and wire it directly into `provider.ts` so that `client.getBooleanValue()` calls via OpenFeature behave identically to `useFlag()`.
- **Decision 2**: We will refactor the iOS `Provider.swift` to use the `MurmurHash3.bucketUser` function introduced in the previous feature, correctly applying targeting rules and sorting variant keys for determinism.
- **Decision 3**: We will refactor the Android `Provider.kt` to use the `MurmurHash3.bucketUser` function, matching the iOS logic.
- **Decision 4**: We will map FlagManagment's evaluation results to standard OpenFeature `Reason` strings (`STATIC`/`TARGETING_MATCH`/`DISABLED`/`DEFAULT`/`ERROR`).
