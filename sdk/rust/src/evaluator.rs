use std::collections::HashMap;
use serde_json::Value;
use crate::murmur::bucket_user;

#[derive(Debug, Clone, PartialEq)]
pub struct EvaluationResult {
    pub value: Value,
    pub variant: String,
    pub reason: String,
}

pub struct Evaluator;

impl Evaluator {
    /// Evaluates a feature flag definition against a targeting key and context.
    pub fn evaluate(
        flag: &Value,
        targeting_key: Option<&str>,
        _context: Option<&HashMap<String, Value>>,
    ) -> EvaluationResult {
        // 1. Check enabled state
        let is_enabled = flag.get("enabled").and_then(|v| v.as_bool()).unwrap_or(false);
        let default_variant = flag
            .get("defaultVariant")
            .and_then(|v| v.as_str())
            .unwrap_or("off")
            .to_string();

        if !is_enabled {
            let val = Self::get_variant_value(flag, &default_variant);
            return EvaluationResult {
                value: val,
                variant: default_variant,
                reason: "DISABLED".to_string(),
            };
        }

        // 2. Check targeting key
        let key = match targeting_key {
            Some(k) if !k.is_empty() => k,
            _ => {
                let val = Self::get_variant_value(flag, &default_variant);
                return EvaluationResult {
                    value: val,
                    variant: default_variant,
                    reason: "DEFAULT".to_string(),
                };
            }
        };

        // 3. Rollout bucketing
        if let Some(rules) = flag.get("rules").and_then(|r| r.as_array()) {
            for rule in rules {
                if let Some(rollout) = rule.get("rollout").and_then(|r| r.as_object()) {
                    let bucket = bucket_user(key);
                    let mut cumulative: f64 = 0.0;

                    for (variant, percentage) in rollout {
                        if let Some(pct) = percentage.as_f64() {
                            cumulative += pct;
                            if (bucket as f64) < cumulative {
                                let val = Self::get_variant_value(flag, variant);
                                return EvaluationResult {
                                    value: val,
                                    variant: variant.clone(),
                                    reason: "TARGETING_MATCH".to_string(),
                                };
                            }
                        }
                    }
                }
            }
        }

        let val = Self::get_variant_value(flag, &default_variant);
        EvaluationResult {
            value: val,
            variant: default_variant,
            reason: "DEFAULT".to_string(),
        }
    }

    fn get_variant_value(flag: &Value, variant: &str) -> Value {
        flag.get("variants")
            .and_then(|vars| vars.get(variant))
            .cloned()
            .unwrap_or(Value::Null)
    }
}
