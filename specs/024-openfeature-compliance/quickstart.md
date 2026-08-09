# Quickstart Validation: OpenFeature API Compliance

To validate the OpenFeature API compliance changes for the React, iOS, and Android SDKs, we will simulate evaluating a feature flag using the standard OpenFeature API in each language.

## 1. Validating React SDK

Create a minimal React test component that uses the standard OpenFeature React SDK client to evaluate a flag, proving that the FlagManagment provider correctly handles it.

```tsx
import { OpenFeatureProvider, useBooleanFlagValue } from '@openfeature/react-sdk';
import { FlagClient, FlagManagmentWebProvider } from 'flagmanagment-react-sdk';

// 1. Initialize FlagManagment client
const client = new FlagClient('test-key', 'http://localhost:8080/api/v1/sdk/stream');

// 2. Register provider
const provider = new FlagManagmentWebProvider(client);

// 3. Inject context
const context = { targetingKey: 'user-123' };

function App() {
  return (
    <OpenFeatureProvider provider={provider} context={context}>
      <FeatureComponent />
    </OpenFeatureProvider>
  );
}

function FeatureComponent() {
  // 4. Use standard OpenFeature hook
  const isEnabled = useBooleanFlagValue('sample-flag', false);
  
  return (
    <div>
      {isEnabled ? <p>Feature is ON!</p> : <p>Feature is OFF</p>}
    </div>
  );
}
```

## 2. Validating Android SDK

In Android Studio, inject the `FlagManagmentProvider` into the global OpenFeature API.

```kotlin
import dev.openfeature.sdk.*
import com.flagmanagment.sdk.*

// 1. Initialize Client
val fmClient = FlagClient(context, "test-key", "http://localhost:8080/api/v1/sdk/stream")
fmClient.connect()

// 2. Set OpenFeature Provider
val provider = FlagManagmentProvider(fmClient)
OpenFeatureAPI.getInstance().provider = provider

// 3. Set Evaluation Context
val evalContext = EvaluationContext(targetingKey = "user-123")
OpenFeatureAPI.getInstance().setEvaluationContext(evalContext)

// 4. Evaluate using standard OpenFeature Client
val client = OpenFeatureAPI.getInstance().getClient()
val isEnabled = client.getBooleanValue("sample-flag", false)

println("Feature enabled: $isEnabled")
```

## 3. Validating iOS SDK

In Xcode, inject the `FlagManagmentProvider` into the global OpenFeature API.

```swift
import OpenFeature
import FlagManagment

// 1. Initialize Client
let url = URL(string: "http://localhost:8080/api/v1/sdk/stream")!
let fmClient = FlagClient(apiKey: "test-key", streamUrl: url)
fmClient.connect()

// 2. Set OpenFeature Provider
let provider = FlagManagmentProvider(client: fmClient)
OpenFeatureAPI.shared.setProvider(provider: provider)

// 3. Set Evaluation Context
let evalContext = MutableContext(targetingKey: "user-123")
OpenFeatureAPI.shared.setEvaluationContext(evalContext: evalContext)

// 4. Evaluate using standard OpenFeature Client
let client = OpenFeatureAPI.shared.getClient()
let isEnabled = client.getBooleanValue(key: "sample-flag", defaultValue: false)

print("Feature enabled: \(isEnabled)")
```

### Success Criteria Verification

When evaluating `sample-flag` with targeting rules, the variants returned should correctly correspond to the MurmurHash3 bucket for `user-123`, proving that the OpenFeature `Provider` implementations are correctly parsing the context and executing local bucketing evaluation.
