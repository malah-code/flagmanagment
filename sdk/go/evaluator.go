package sdk

import (
	"fmt"
	"sort"

	"github.com/open-feature/go-sdk/openfeature"
	"github.com/spaolacci/murmur3"
)

// bucketUser hashes the targeting key and returns a bucket between 0 and 99.
func bucketUser(targetingKey string) int {
	hasher := murmur3.New32()
	_, _ = hasher.Write([]byte(targetingKey))
	hashValue := hasher.Sum32()
	return int(hashValue % 100)
}

// evaluateFlag evaluates a flag against the provided EvaluationContext
func evaluateFlag(flag map[string]interface{}, evalContext openfeature.EvaluationContext) (interface{}, string, openfeature.Reason, error) {
	// 1. Check if flag is enabled
	enabled, ok := flag["enabled"].(bool)
	if !ok || !enabled {
		defaultVariant, _ := flag["defaultVariant"].(string)
		return getVariantValue(flag, defaultVariant), defaultVariant, openfeature.DisabledReason, nil
	}

	// Extract targeting key
	targetingKey := evalContext.TargetingKey()
	if targetingKey == "" {
		defaultVariant, _ := flag["defaultVariant"].(string)
		return getVariantValue(flag, defaultVariant), defaultVariant, openfeature.DefaultReason, nil
	}

	// 2. Evaluate targeting rules
	rules, _ := flag["rules"].([]interface{})
	for _, ruleObj := range rules {
		rule, ok := ruleObj.(map[string]interface{})
		if !ok {
			continue
		}

		rollout, ok := rule["rollout"].(map[string]interface{})
		if ok {
			variant, err := getVariantByRollout(targetingKey, rollout)
			if err == nil {
				return getVariantValue(flag, variant), variant, openfeature.TargetingMatchReason, nil
			}
		}
	}

	// 3. Fallback to default variant
	defaultVariant, _ := flag["defaultVariant"].(string)
	return getVariantValue(flag, defaultVariant), defaultVariant, openfeature.DefaultReason, nil
}

func getVariantByRollout(targetingKey string, rollout map[string]interface{}) (string, error) {
	bucket := bucketUser(targetingKey)

	// Sort variant keys for deterministic iteration order.
	// Go map iteration is randomized — without sorting, the same user
	// could land in different variants across evaluations.
	keys := make([]string, 0, len(rollout))
	for k := range rollout {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	currentThreshold := 0
	for _, variant := range keys {
		percentageObj := rollout[variant]
		percentage, ok := percentageObj.(float64)
		if !ok {
			percentageInt, ok := percentageObj.(int)
			if !ok {
				continue
			}
			percentage = float64(percentageInt)
		}

		currentThreshold += int(percentage)
		if bucket < currentThreshold {
			return variant, nil
		}
	}

	return "", fmt.Errorf("no variant matched, bucket %d exceeded total percentage %d", bucket, currentThreshold)
}

func getVariantValue(flag map[string]interface{}, variant string) interface{} {
	variants, ok := flag["variants"].(map[string]interface{})
	if !ok {
		return nil
	}
	return variants[variant]
}
