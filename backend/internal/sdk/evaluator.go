package sdk

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"regexp"
	"strings"
	"sync"

	"github.com/flagmanagment/backend/internal/models"
	"github.com/spaolacci/murmur3"
)

// EvaluationContext represents the stringified context for internal evaluation matching
type EvaluationContext map[string]string

var (
	regexCache sync.Map
)

// HashPII creates a shallow copy of the EvaluationContext and hashes PII fields
func HashPII(ctx *models.EvaluationContext) *models.EvaluationContext {
	if ctx == nil {
		return nil
	}
	hashedCtx := &models.EvaluationContext{
		EntityKey:  ctx.EntityKey,
		Attributes: make(map[string]interface{}),
	}
	
	// List of known PII fields to hash based on Constitution VII
	piiFields := map[string]bool{
		"email": true,
		"ssn":   true,
		"phone": true,
	}

	for k, v := range ctx.Attributes {
		if piiFields[k] {
			if strVal, ok := v.(string); ok {
				hash := sha256.Sum256([]byte(strVal))
				hashedCtx.Attributes[k] = hex.EncodeToString(hash[:])
				continue
			}
		}
		hashedCtx.Attributes[k] = v
	}
	return hashedCtx
}

// EvaluateRolloutSplit evaluates percentage rollout using MurmurHash3
func EvaluateRolloutSplit(flagKey string, entityKey string, rolloutRulesJSON models.JSONB) (string, bool) {
	if rolloutRulesJSON == nil || entityKey == "" {
		return "", false
	}
	rulesInterface, ok := rolloutRulesJSON["rules"]
	if !ok {
		return "", false
	}
	rBytes, err := json.Marshal(rulesInterface)
	if err != nil {
		return "", false
	}
	var rollouts []models.RolloutRule
	if err := json.Unmarshal(rBytes, &rollouts); err != nil || len(rollouts) == 0 {
		return "", false
	}

	// MurmurHash3 32-bit calculation
	hash := murmur3.Sum32([]byte(flagKey + ":" + entityKey))
	bucket := int(hash % 10000) // 0 - 9999

	cumulative := 0
	for _, r := range rollouts {
		cumulative += r.Percentage
		if bucket < cumulative {
			return r.VariationID, true
		}
	}
	return "", false
}

// EvaluateFlag orchestrates the evaluation of a single flag against the context
func EvaluateFlag(flagRule *models.FlagRule, ctx *models.EvaluationContext, rulesMap map[string]*models.FlagRule) models.EvaluationResult {
	if !flagRule.Enabled {
		return models.EvaluationResult{
			Value:  flagRule.DefaultVariation,
			Reason: "FLAG_DISABLED",
		}
	}

	if flagRule.ParentFlagKey != "" && rulesMap != nil {
		parentRule, exists := rulesMap[flagRule.ParentFlagKey]
		if exists && parentRule != nil {
			parentResult := EvaluateFlag(parentRule, ctx, rulesMap)
			if parentResult.Reason == "FLAG_DISABLED" || parentResult.Value == "false" || parentResult.Value == false {
				return models.EvaluationResult{
					Value:  flagRule.DefaultVariation,
					Reason: "PARENT_FLAG_DISABLED",
				}
			}
		}
	}
	
	entityKey := ""
	if ctx != nil {
		entityKey = ctx.EntityKey
	}

	if flagRule.TargetingRules != nil {
		// Convert context for the rules engine
		evalCtx := make(EvaluationContext)
		if ctx != nil {
			for k, v := range ctx.Attributes {
				if strVal, ok := v.(string); ok {
					evalCtx[k] = strVal
				}
			}
		}

		// Convert JSONRawMessage to JSONB map structure for our internal evaluator
		var rulesMap models.JSONB
		if err := json.Unmarshal(flagRule.TargetingRules, &rulesMap); err == nil {
			variation, matched := EvaluateContextualRules(rulesMap, evalCtx)
			if matched {
				return models.EvaluationResult{
					Value:  variation,
					Reason: "TARGETING_MATCH",
				}
			}
		}
	}

	if flagRule.RolloutRules != nil && entityKey != "" {
		var rolloutMap models.JSONB
		if err := json.Unmarshal(flagRule.RolloutRules, &rolloutMap); err == nil {
			variationID, matched := EvaluateRolloutSplit(flagRule.Key, entityKey, rolloutMap)
			if matched {
				return models.EvaluationResult{
					Value:  variationID,
					Reason: "PERCENTAGE_ROLLOUT",
				}
			}
		}
	}

	return models.EvaluationResult{
		Value:  flagRule.DefaultVariation,
		Reason: "DEFAULT",
	}
}

// getCompiledRegex retrieves a compiled regex from the cache or compiles and caches it
func getCompiledRegex(pattern string) (*regexp.Regexp, error) {
	if cached, ok := regexCache.Load(pattern); ok {
		return cached.(*regexp.Regexp), nil
	}
	compiled, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	regexCache.Store(pattern, compiled)
	return compiled, nil
}

// EvaluateContextualRules evaluates targeting rules against a provided context.
// It returns a boolean and whether a rule matched.
func EvaluateContextualRules(targetingRulesJSON models.JSONB, context EvaluationContext) (bool, bool) {
	if targetingRulesJSON == nil {
		return false, false
	}
	rulesInterface, ok := targetingRulesJSON["rules"]
	if !ok {
		return false, false
	}
	
	rulesBytes, err := json.Marshal(rulesInterface)
	if err != nil {
		return false, false
	}
	
	var rules []models.TargetingRule
	if err := json.Unmarshal(rulesBytes, &rules); err != nil {
		return false, false
	}

	for _, rule := range rules {
		if evaluateRule(rule, context) {
			return rule.Variation, true
		}
	}
	return false, false
}

// evaluateRule evaluates a single targeting rule (AND logic for its conditions)
func evaluateRule(rule models.TargetingRule, context EvaluationContext) bool {
	if len(rule.Conditions) == 0 {
		return false
	}

	for _, condition := range rule.Conditions {
		if !evaluateCondition(condition, context) {
			return false
		}
	}
	return true
}

// evaluateCondition evaluates a single condition against the context
func evaluateCondition(condition models.TargetingCondition, context EvaluationContext) bool {
	contextValue, ok := context[condition.Attribute]
	if !ok {
		return false // attribute missing in context -> condition false
	}

	switch condition.Operator {
	case models.OperatorEquals:
		return contextValue == condition.Value
	case models.OperatorContains:
		return strings.Contains(contextValue, condition.Value)
	case models.OperatorRegex:
		re, err := getCompiledRegex(condition.Value)
		if err != nil {
			return false // invalid regex -> condition false
		}
		return re.MatchString(contextValue)
	default:
		return false // unknown operator -> condition false
	}
}
