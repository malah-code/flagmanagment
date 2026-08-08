package hooks

import (
	"github.com/flagmanagment/backend/internal/models"
)

// HookContext provides context about the evaluation
type HookContext struct {
	FlagKey           string
	FlagType          string
	EvaluationContext *models.EvaluationContext
}

// EvaluationDetails provides details about the evaluation result
type EvaluationDetails struct {
	FlagKey string
	Value   interface{}
	Reason  string
}

// Hook represents the OpenFeature Hook interface
type Hook interface {
	After(ctx HookContext, details EvaluationDetails)
	Error(ctx HookContext, err error)
}
