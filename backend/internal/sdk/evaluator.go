package sdk

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/flagmanagment/backend/internal/models"
)

// EvaluateFlag performs local evaluation of a flag rule against context.
func EvaluateFlag(flag *models.FlagRule, ctx *models.EvaluationContext) models.EvaluationResult {
	if !flag.Enabled {
		return models.EvaluationResult{
			Value:  flag.DefaultVariation,
			Reason: "DISABLED",
		}
	}

	// Basic evaluation logic (MVP defaults to DefaultVariation)
	return models.EvaluationResult{
		Value:  flag.DefaultVariation,
		Reason: "DEFAULT",
	}
}

// HashPII hashes sensitive user context attributes per Constitution VII.
func HashPII(ctx *models.EvaluationContext) *models.EvaluationContext {
	if ctx == nil {
		return nil
	}

	hashedCtx := &models.EvaluationContext{
		EntityKey:  hashString(ctx.EntityKey),
		Attributes: make(map[string]interface{}),
	}

	for k, v := range ctx.Attributes {
		if strVal, ok := v.(string); ok {
			hashedCtx.Attributes[k] = hashString(strVal)
		} else {
			hashedCtx.Attributes[k] = v
		}
	}

	return hashedCtx
}

func hashString(s string) string {
	if s == "" {
		return ""
	}
	hash := sha256.Sum256([]byte(s))
	return hex.EncodeToString(hash[:])
}
