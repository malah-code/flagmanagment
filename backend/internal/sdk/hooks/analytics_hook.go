package hooks

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// AnalyticsHook is a reference implementation of an OpenFeature Hook that logs
// evaluation events to standard output (simulating an external analytics provider).
type AnalyticsHook struct {
	Prefix string
	Writer io.Writer
}

// NewAnalyticsHook creates a new AnalyticsHook with the given log prefix
func NewAnalyticsHook(prefix string, writers ...io.Writer) *AnalyticsHook {
	var w io.Writer = os.Stdout
	if len(writers) > 0 && writers[0] != nil {
		w = writers[0]
	}
	return &AnalyticsHook{
		Prefix: prefix,
		Writer: w,
	}
}

// After is called after the flag evaluation is complete.
// It formats the EvaluationDetails into an analytics event payload and logs it.
func (h *AnalyticsHook) After(ctx HookContext, details EvaluationDetails) {
	event := map[string]interface{}{
		"event":     "flag_evaluation",
		"flag_key":  details.FlagKey,
		"flag_type": ctx.FlagType,
		"value":     details.Value,
		"reason":    details.Reason,
	}

	if ctx.EvaluationContext != nil {
		event["entity_key"] = ctx.EvaluationContext.EntityKey
	}

	b, err := json.Marshal(event)
	if err != nil {
		fmt.Fprintf(h.Writer, "[%s] Error marshaling analytics event: %v\n", h.Prefix, err)
		return
	}

	fmt.Fprintf(h.Writer, "[%s] ANALYTICS_EVENT: %s\n", h.Prefix, string(b))
}

// Error is called if the evaluation encounters an error
func (h *AnalyticsHook) Error(ctx HookContext, err error) {
	fmt.Fprintf(h.Writer, "[%s] ERROR_EVENT: flag_key=%s err=%v\n", h.Prefix, ctx.FlagKey, err)
}
