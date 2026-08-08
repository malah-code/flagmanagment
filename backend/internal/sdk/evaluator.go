package sdk

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"regexp"
	"strings"
	"sync"

	"github.com/flagmanagment/backend/internal/models"
	"github.com/flagmanagment/backend/internal/sdk/hooks"
	"github.com/spaolacci/murmur3"
)

var (
	regexCache     sync.Map
	hooksMu        sync.RWMutex
	registeredHooks []hooks.Hook
)

// RegisterHook adds an OpenFeature Hook to the evaluator pipeline
func RegisterHook(hook hooks.Hook) {
	hooksMu.Lock()
	defer hooksMu.Unlock()
	registeredHooks = append(registeredHooks, hook)
}

// ClearHooks clears all registered hooks (primarily for testing)
func ClearHooks() {
	hooksMu.Lock()
	defer hooksMu.Unlock()
	registeredHooks = nil
}

// executeAfterHooks runs the After hooks asynchronously
func executeAfterHooks(ctx hooks.HookContext, details hooks.EvaluationDetails) {
	hooksMu.RLock()
	hooksToRun := make([]hooks.Hook, len(registeredHooks))
	copy(hooksToRun, registeredHooks)
	hooksMu.RUnlock()

	for _, hook := range hooksToRun {
		go func(h hooks.Hook) {
			defer func() {
				if r := recover(); r != nil {
					// recover from hook panic safely
				}
			}()
			h.After(ctx, details)
		}(hook)
	}
}

// executeErrorHooks runs the Error hooks asynchronously
func executeErrorHooks(ctx hooks.HookContext, err error) {
	hooksMu.RLock()
	hooksToRun := make([]hooks.Hook, len(registeredHooks))
	copy(hooksToRun, registeredHooks)
	hooksMu.RUnlock()

	for _, hook := range hooksToRun {
		go func(h hooks.Hook) {
			defer func() {
				if r := recover(); r != nil {
					// recover from hook panic safely
				}
			}()
			h.Error(ctx, err)
		}(hook)
	}
}

// EvaluationContext represents the stringified context for internal evaluation matching
type EvaluationContext map[string]string

// HashPII creates a shallow copy of the EvaluationContext and hashes PII fields using SHA-256 and the provided salt
func HashPII(ctx *models.EvaluationContext, salt string) *models.EvaluationContext {
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
				hash := sha256.Sum256([]byte(salt + strVal))
				hashedCtx.Attributes[k] = hex.EncodeToString(hash[:])
				continue
			}
		}
		hashedCtx.Attributes[k] = v
	}
	return hashedCtx
}

// EvaluateRolloutSplit evaluates percentage rollout using MurmurHash3 and an environment salt
func EvaluateRolloutSplit(flagKey string, entityKey string, rolloutRulesJSON models.JSONB, salt string) (string, bool) {
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

	// MurmurHash3 32-bit calculation with environment salt
	hash := murmur3.Sum32([]byte(flagKey + ":" + entityKey + ":" + salt))
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
func EvaluateFlag(flagRule *models.FlagRule, ctx *models.EvaluationContext, rulesMap map[string]*models.FlagRule, salt ...string) models.EvaluationResult {
	sVal := ""
	if len(salt) > 0 {
		sVal = salt[0]
	}

	hookCtx := hooks.HookContext{
		FlagKey:           flagRule.Key,
		FlagType:          flagRule.Type,
		EvaluationContext: ctx,
	}

	if !flagRule.Enabled {
		res := models.EvaluationResult{
			Value:  flagRule.DefaultVariation,
			Reason: "FLAG_DISABLED",
		}
		executeAfterHooks(hookCtx, hooks.EvaluationDetails{
			FlagKey: flagRule.Key,
			Value:   res.Value,
			Reason:  res.Reason,
		})
		return res
	}

	if flagRule.ParentFlagKey != "" && rulesMap != nil {
		parentRule, exists := rulesMap[flagRule.ParentFlagKey]
		if exists && parentRule != nil {
			parentResult := EvaluateFlag(parentRule, ctx, rulesMap, sVal)
			if parentResult.Reason == "FLAG_DISABLED" || parentResult.Value == "false" || parentResult.Value == false {
				res := models.EvaluationResult{
					Value:  flagRule.DefaultVariation,
					Reason: "PARENT_FLAG_DISABLED",
				}
				executeAfterHooks(hookCtx, hooks.EvaluationDetails{
					FlagKey: flagRule.Key,
					Value:   res.Value,
					Reason:  res.Reason,
				})
				return res
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
				res := models.EvaluationResult{
					Value:  variation,
					Reason: "TARGETING_MATCH",
				}
				executeAfterHooks(hookCtx, hooks.EvaluationDetails{
					FlagKey: flagRule.Key,
					Value:   res.Value,
					Reason:  res.Reason,
				})
				return res
			}
		}
	}

	if flagRule.RolloutRules != nil && entityKey != "" {
		var rolloutMap models.JSONB
		if err := json.Unmarshal(flagRule.RolloutRules, &rolloutMap); err == nil {
			variationID, matched := EvaluateRolloutSplit(flagRule.Key, entityKey, rolloutMap, sVal)
			if matched {
				res := models.EvaluationResult{
					Value:  variationID,
					Reason: "PERCENTAGE_ROLLOUT",
				}
				executeAfterHooks(hookCtx, hooks.EvaluationDetails{
					FlagKey: flagRule.Key,
					Value:   res.Value,
					Reason:  res.Reason,
				})
				return res
			}
		}
	}

	res := models.EvaluationResult{
		Value:  flagRule.DefaultVariation,
		Reason: "DEFAULT",
	}
	executeAfterHooks(hookCtx, hooks.EvaluationDetails{
		FlagKey: flagRule.Key,
		Value:   res.Value,
		Reason:  res.Reason,
	})
	return res
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
