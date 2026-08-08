package hooks

import (
	"encoding/json"
	"fmt"
	"log"
)

// AnalyticsHook is a reference implementation of an OpenFeature Hook that logs
// evaluation events to standard output (simulating an external analytics provider).
type AnalyticsHook struct {
	Prefix string
}

// NewAnalyticsHook creates a new AnalyticsHook with the given log prefix
func NewAnalyticsHook(prefix string) *AnalyticsHook {
	return &AnalyticsHook{
		Prefix: prefix,
	}
}

// After is called after the flag evaluation is complete.
// It formats the EvaluationDetails into an analytics event payload and logs it.
func (h *AnalyticsHook) After(ctx HookContext, details EvaluationDetails) {
	event := map[string]interface{}{
		"event":       "flag_evaluation",
		"flag_key":    details.FlagKey,
		"flag_type":   ctx.FlagType,
		"value":       details.Value,
		"reason":      details.Reason,
	}

	if ctx.EvaluationContext != nil {
		event["entity_key"] = ctx.EvaluationContext.EntityKey
		// Be careful not to log full PII in a real implementation
		// Here it's already hashed by HashPII
	}

	b, err := json.Marshal(event)
	if err != nil {
		log.Printf("[%s] Error marshaling analytics event: %v", h.Prefix, err)
		return
	}

	// In a real implementation (e.g. PostHog, Amplitude), this is where the API call happens.
	// We simulate this network call by printing to standard output.
	fmt.Printf("[%s] ANALYTICS_EVENT: %s\n", h.Prefix, string(b))
}

// Error is called if the evaluation encounters an error
func (h *AnalyticsHook) Error(ctx HookContext, err error) {
	log.Printf("[%s] ERROR_EVENT: flag_key=%s err=%v\n", h.Prefix, ctx.FlagKey, err)
}
