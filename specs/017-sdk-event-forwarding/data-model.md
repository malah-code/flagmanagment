# Data Model: SDK Event Forwarding for Analytics

As this feature focuses on SDK-side event interception and not backend persistence, the data model defines the internal structures passed within the SDK evaluation lifecycle rather than database tables.

## Core Entities

### EvaluationContext (Provided by OpenFeature)
- **TargetingKey**: String (e.g., User ID)
- **Attributes**: Map<String, interface{}> (e.g., tenant_id, email, region)

### EvaluationEvent / HookContext
- **FlagKey**: String
- **FlagType**: Boolean, String, Number, Object
- **DefaultValue**: interface{}
- **EvaluationContext**: The context passed into the evaluation
- **ProviderMetadata**: Metadata about the FlagManagment provider

### EvaluationDetails (Result)
- **FlagKey**: String
- **Value**: interface{} (The assigned variant value)
- **Variant**: String (The assigned variant ID or Name)
- **Reason**: String (e.g., "TARGETING_MATCH", "FALLBACK")
- **ErrorCode**: String (if an error occurred)
- **ErrorMessage**: String (if an error occurred)

## Interface Contracts (SDK Hooks)

The SDK must implement the standard OpenFeature `Hook` interface. The primary method utilized for event forwarding will be `After`:

```go
type Hook interface {
	Before(ctx HookContext, hookHints HookHints) (EvaluationContext, error)
	After(ctx HookContext, details EvaluationDetails, hookHints HookHints) error
	Error(ctx HookContext, err error, hookHints HookHints)
	Finally(ctx HookContext, hookHints HookHints)
}
```

- When `EvaluateFlag` is called, the SDK processes the targeting rules.
- Once a variant is resolved, the SDK calls the `After` method of all registered hooks.
- An Analytics Provider Hook (e.g., `PostHogHook`) implements the `After` method, formats the `EvaluationDetails` into an analytics event, and asynchronously transmits it via its respective client without blocking the SDK.
