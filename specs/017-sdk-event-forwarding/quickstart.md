# Quickstart: SDK Event Forwarding for Analytics

This guide demonstrates how to configure the FlagManagment Go SDK to automatically intercept flag evaluations and print the results (simulating forwarding to an analytics provider).

## Prerequisites

- A running instance of the FlagManagment backend.
- A Go application with the FlagManagment SDK installed.

## Validation Scenarios

### Scenario 1: Intercepting Flag Evaluation with a Logging Hook

1. **Create the Custom Hook**:
   In your Go application, define a struct that implements the OpenFeature `Hook` interface.

   ```go
   package main

   import (
       "fmt"
       "github.com/open-feature/go-sdk/openfeature"
   )

   type AnalyticsHook struct{}

   func (h AnalyticsHook) Before(ctx openfeature.HookContext, hints openfeature.HookHints) (*openfeature.EvaluationContext, error) {
       return nil, nil
   }

   func (h AnalyticsHook) After(ctx openfeature.HookContext, details openfeature.InterfaceEvaluationDetails, hints openfeature.HookHints) error {
       // Asynchronously "forward" the event
       go func() {
           fmt.Printf("[ANALYTICS] User '%s' evaluated flag '%s' and received variant: %v\n", 
               ctx.ClientMetadata().Name(), ctx.FlagKey(), details.Value)
       }()
       return nil
   }

   func (h AnalyticsHook) Error(ctx openfeature.HookContext, err error, hints openfeature.HookHints) {}
   func (h AnalyticsHook) Finally(ctx openfeature.HookContext, hints openfeature.HookHints) {}
   ```

2. **Register the Hook and Evaluate**:
   Register the hook with the FlagManagment provider (or at the global OpenFeature level) and evaluate a flag.

   ```go
   openfeature.AddHooks(AnalyticsHook{})
   client := openfeature.NewClient("my-app")
   
   evalCtx := openfeature.NewEvaluationContext("user-123", nil)
   val, _ := client.BooleanValue("new-checkout-flow", false, evalCtx)
   ```

3. **Verify the Output**:
   When you run the application, you should immediately see the analytics log line printed to standard output without delaying the flag evaluation result.

   **Expected Console Output**:
   ```
   [ANALYTICS] User 'my-app' evaluated flag 'new-checkout-flow' and received variant: true
   ```
