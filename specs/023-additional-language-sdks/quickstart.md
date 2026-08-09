# Quickstart & Validation Guide: Additional Language SDKs

## Prerequisites
- A running FlagManagment backend or edge proxy.
- A valid Environment API Key (`fm_sa_*`).
- Development environments for Go, Java, Python, .NET, React, iOS, Android.

---

## Language Snippets

### Go
```go
import "github.com/malah-code/flagmanagment/sdk/go"
import "github.com/open-feature/go-sdk/openfeature"

client := sdk.NewProvider("fm_sa_key", "http://localhost:8080/api/v1/sdk/stream")
openfeature.SetProvider(client)
openfeature.Init()

evalCtx := openfeature.NewEvaluationContext("user-123", nil)
boolValue, _ := openfeature.GetClient().BooleanValue("test-flag", false, evalCtx)
```

### Java
```java
import com.flagmanagment.sdk.*;
import dev.openfeature.sdk.*;

Client fmClient = new Client("fm_sa_key", "http://localhost:8080/api/v1/sdk/stream");
fmClient.connect();

OpenFeatureAPI.getInstance().setProvider(new Provider(fmClient));
Client ofClient = OpenFeatureAPI.getInstance().getClient();
boolean flagValue = ofClient.getBooleanValue("test-flag", false);
```

### Python
```python
from openfeature import api, EvaluationContext
from flagmanagment.client import Client
from flagmanagment.provider import FlagManagmentProvider

fm_client = Client("fm_sa_key", "http://localhost:8080/api/v1/sdk/stream")
fm_client.connect()

api.set_provider(FlagManagmentProvider(fm_client))
client = api.get_client()

flag_value = client.get_boolean_value("test-flag", False, EvaluationContext(targeting_key="user-123"))
```

### .NET
```csharp
using FlagManagment.Sdk;
using OpenFeature;

var fmClient = new Client("fm_sa_key", "http://localhost:8080/api/v1/sdk/stream");
fmClient.Connect();

Api.Instance.SetProviderAsync(new Provider(fmClient)).Wait();
var client = Api.Instance.GetClient();
var flagValue = client.GetBooleanValue("test-flag", false);
```

### React
```tsx
import { FlagProvider, useFlag } from '@flagmanagment/react-sdk';

function App() {
  return (
    <FlagProvider apiKey="fm_sa_key" streamUrl="http://localhost:8080/api/v1/sdk/stream" context={{targetingKey: 'user-123'}}>
      <FeatureComponent />
    </FlagProvider>
  );
}

function FeatureComponent() {
  const isEnabled = useFlag("test-flag", false);
  return isEnabled ? <NewFeature /> : <OldFeature />;
}
```

### iOS (Swift)
```swift
import FlagManagment
import OpenFeature

let client = FlagClient(apiKey: "fm_sa_key", streamUrl: URL(string: "http://localhost:8080/api/v1/sdk/stream")!)
client.connect()

OpenFeatureAPI.shared.setProvider(provider: FlagManagmentProvider(client: client))
let ofClient = OpenFeatureAPI.shared.getClient()

let isEnabled = ofClient.getBooleanValue(key: "test-flag", defaultValue: false)
```

### Android (Kotlin)
```kotlin
import com.flagmanagment.sdk.*
import dev.openfeature.sdk.*

val client = FlagClient(context, "fm_sa_key", "http://localhost:8080/api/v1/sdk/stream")
client.connect()

OpenFeatureAPI.setProvider(FlagManagmentProvider(client), null)
val ofClient = OpenFeatureAPI.getClient()

val isEnabled = ofClient.getBooleanValue("test-flag", false)
```

---

## Validation Scenarios

### Scenario 1: Go Server-Side In-Memory Evaluation
1. Initialize the Go SDK with an API key.
2. Ensure the SDK successfully connects to the backend streaming endpoint and receives the bootstrap ruleset.
3. Call the SDK's evaluation method (`GetBooleanValue("test-flag", context, false)`).
4. **Validation**: The evaluation must complete locally without network I/O and return the correct value based on targeting rules.
5. Create a load test script using `goroutines` to hit the evaluation method 10,000 times concurrently.
6. **Validation**: Verify p99 latency < 1ms.

### Scenario 2: Streaming Delta Sync
1. Start an SDK application (any language).
2. Through the backend REST API, modify the rollout percentage of a flag.
3. **Validation**: The SDK should receive the `flag_updated` SSE event and log the ruleset change within 2 seconds. Subsequent evaluations should reflect the new percentage immediately.

### Scenario 3: React Declarative UI
1. Wrap a sample React application with `<FlagProvider apiKey="..." userContext={{key:"user-1"}}>`.
2. Inside a component, use `const { value } = useFlag("new-checkout")`.
3. Change the flag value in the backend.
4. **Validation**: The component should automatically re-render with the new value.

### Scenario 4: Mobile Offline Cache (iOS/Android)
1. Launch the mobile application. Ensure it downloads the initial ruleset.
2. Put the device in Airplane Mode (no network).
3. Evaluate a flag via the SDK.
4. **Validation**: The SDK must evaluate the flag using the locally encrypted cache and return the correct value without hanging or blocking the main thread.
