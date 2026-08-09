package com.flagmanagment.sdk;

import dev.openfeature.sdk.*;
import org.json.JSONObject;

import java.util.TreeMap;

/**
 * OpenFeature provider for FlagManagment. Evaluates flags locally
 * using the in-memory cache maintained by {@link Client} and
 * MurmurHash3 bucketing via {@link Evaluator}.
 */
public class Provider implements FeatureProvider {
    private final Client client;

    public Provider(Client client) {
        this.client = client;
    }

    @Override
    public Metadata getMetadata() {
        return () -> "FlagManagment-Java-Provider";
    }

    @Override
    public ProviderEvaluation<Boolean> getBooleanEvaluation(String key, Boolean defaultValue, EvaluationContext ctx) {
        EvalResult result = evaluate(key, ctx);
        if (result == null) {
            return ProviderEvaluation.<Boolean>builder()
                .value(defaultValue)
                .reason(Reason.ERROR.toString())
                .build();
        }
        if (result.value instanceof Boolean) {
            return ProviderEvaluation.<Boolean>builder()
                .value((Boolean) result.value)
                .variant(result.variant)
                .reason(result.reason)
                .build();
        }
        return ProviderEvaluation.<Boolean>builder()
            .value(defaultValue)
            .reason(Reason.ERROR.toString())
            .build();
    }

    @Override
    public ProviderEvaluation<String> getStringEvaluation(String key, String defaultValue, EvaluationContext ctx) {
        EvalResult result = evaluate(key, ctx);
        if (result == null) {
            return ProviderEvaluation.<String>builder().value(defaultValue).reason(Reason.ERROR.toString()).build();
        }
        if (result.value instanceof String) {
            return ProviderEvaluation.<String>builder().value((String) result.value).variant(result.variant).reason(result.reason).build();
        }
        return ProviderEvaluation.<String>builder().value(defaultValue).reason(Reason.ERROR.toString()).build();
    }

    @Override
    public ProviderEvaluation<Integer> getIntegerEvaluation(String key, Integer defaultValue, EvaluationContext ctx) {
        EvalResult result = evaluate(key, ctx);
        if (result == null) {
            return ProviderEvaluation.<Integer>builder().value(defaultValue).reason(Reason.ERROR.toString()).build();
        }
        if (result.value instanceof Number) {
            return ProviderEvaluation.<Integer>builder().value(((Number) result.value).intValue()).variant(result.variant).reason(result.reason).build();
        }
        return ProviderEvaluation.<Integer>builder().value(defaultValue).reason(Reason.ERROR.toString()).build();
    }

    @Override
    public ProviderEvaluation<Double> getDoubleEvaluation(String key, Double defaultValue, EvaluationContext ctx) {
        EvalResult result = evaluate(key, ctx);
        if (result == null) {
            return ProviderEvaluation.<Double>builder().value(defaultValue).reason(Reason.ERROR.toString()).build();
        }
        if (result.value instanceof Number) {
            return ProviderEvaluation.<Double>builder().value(((Number) result.value).doubleValue()).variant(result.variant).reason(result.reason).build();
        }
        return ProviderEvaluation.<Double>builder().value(defaultValue).reason(Reason.ERROR.toString()).build();
    }

    @Override
    public ProviderEvaluation<Value> getObjectEvaluation(String key, Value defaultValue, EvaluationContext ctx) {
        return ProviderEvaluation.<Value>builder().value(defaultValue).reason(Reason.DEFAULT.toString()).build();
    }

    // --- Internal evaluation logic ---

    private static class EvalResult {
        final Object value;
        final String variant;
        final String reason;

        EvalResult(Object value, String variant, String reason) {
            this.value = value;
            this.variant = variant;
            this.reason = reason;
        }
    }

    private EvalResult evaluate(String flagKey, EvaluationContext ctx) {
        Object flagObj = client.getFlag(flagKey);
        if (flagObj == null) return null;

        JSONObject flag;
        if (flagObj instanceof JSONObject) {
            flag = (JSONObject) flagObj;
        } else {
            return null;
        }

        // Check enabled
        if (!flag.optBoolean("enabled", false)) {
            String defaultVariant = flag.optString("defaultVariant", "");
            Object value = getVariantValue(flag, defaultVariant);
            return new EvalResult(value, defaultVariant, Reason.DISABLED.toString());
        }

        // Extract targeting key
        String targetingKey = ctx != null ? ctx.getTargetingKey() : null;
        if (targetingKey == null || targetingKey.isEmpty()) {
            String defaultVariant = flag.optString("defaultVariant", "");
            Object value = getVariantValue(flag, defaultVariant);
            return new EvalResult(value, defaultVariant, Reason.DEFAULT.toString());
        }

        // Evaluate rules with rollout
        org.json.JSONArray rules = flag.optJSONArray("rules");
        if (rules != null) {
            for (int i = 0; i < rules.length(); i++) {
                JSONObject rule = rules.optJSONObject(i);
                if (rule == null) continue;

                JSONObject rollout = rule.optJSONObject("rollout");
                if (rollout != null) {
                    String variant = getVariantByRollout(targetingKey, rollout);
                    if (variant != null) {
                        Object value = getVariantValue(flag, variant);
                        return new EvalResult(value, variant, Reason.TARGETING_MATCH.toString());
                    }
                }
            }
        }

        // Fallback to default
        String defaultVariant = flag.optString("defaultVariant", "");
        Object value = getVariantValue(flag, defaultVariant);
        return new EvalResult(value, defaultVariant, Reason.DEFAULT.toString());
    }

    private String getVariantByRollout(String targetingKey, JSONObject rollout) {
        int bucket = Evaluator.bucketUser(targetingKey);

        // Sort keys for deterministic iteration
        TreeMap<String, Integer> sorted = new TreeMap<>();
        for (String key : rollout.keySet()) {
            sorted.put(key, rollout.optInt(key, 0));
        }

        int currentThreshold = 0;
        for (var entry : sorted.entrySet()) {
            currentThreshold += entry.getValue();
            if (bucket < currentThreshold) {
                return entry.getKey();
            }
        }
        return null;
    }

    private Object getVariantValue(JSONObject flag, String variant) {
        JSONObject variants = flag.optJSONObject("variants");
        if (variants == null || variant == null) return null;
        return variants.opt(variant);
    }
}
