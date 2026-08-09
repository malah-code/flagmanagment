package com.flagmanagment.sdk

import dev.openfeature.sdk.*
import org.json.JSONObject

class FlagManagmentProvider(private val client: FlagClient) : FeatureProvider {
    override val hooks: List<Hook<*>> = listOf()
    override val metadata: Metadata = FlagManagmentMetadata()

    private class FlagManagmentMetadata : Metadata {
        override val name: String = "FlagManagment-Android-Provider"
    }

    private fun <T> evaluateFlagInternal(
        key: String,
        defaultValue: T,
        ctx: EvaluationContext?,
        extractValue: (JSONObject?, String) -> T?
    ): ProviderEvaluation<T> {
        val flag = client.getFlag(key) as? JSONObject
            ?: return ProviderEvaluation(defaultValue, reason = Reason.ERROR.name)
            
        val enabled = flag.optBoolean("enabled", false)
        val defaultVariant = flag.optString("defaultVariant")
        val variants = flag.optJSONObject("variants")

        if (!enabled) {
            val value = extractValue(variants, defaultVariant)
                ?: return ProviderEvaluation(defaultValue, reason = Reason.ERROR.name)
            return ProviderEvaluation(value, reason = Reason.DISABLED.name)
        }
        
        val targetingKey = ctx?.targetingKey
        if (targetingKey != null) {
            val rules = flag.optJSONArray("rules")
            if (rules != null) {
                for (i in 0 until rules.length()) {
                    val rule = rules.optJSONObject(i)
                    if (rule != null) {
                        val rollout = rule.optJSONObject("rollout")
                        if (rollout != null) {
                            val bucket = MurmurHash3.bucketUser(targetingKey)
                            var currentThreshold = 0
                            
                            val iter = rollout.keys()
                            val keys = mutableListOf<String>()
                            while (iter.hasNext()) keys.add(iter.next())
                            keys.sort()
                            
                            for (variant in keys) {
                                val weight = rollout.optInt(variant, 0)
                                currentThreshold += weight
                                if (bucket < currentThreshold) {
                                    val value = extractValue(variants, variant)
                                        ?: return ProviderEvaluation(defaultValue, reason = Reason.ERROR.name)
                                    return ProviderEvaluation(value, reason = Reason.TARGETING_MATCH.name)
                                }
                            }
                        }
                    }
                }
            }
        }
        
        val value = extractValue(variants, defaultVariant)
            ?: return ProviderEvaluation(defaultValue, reason = Reason.ERROR.name)
        return ProviderEvaluation(value, reason = Reason.STATIC.name)
    }

    override fun getBooleanEvaluation(
        key: String,
        defaultValue: Boolean,
        ctx: EvaluationContext?
    ): ProviderEvaluation<Boolean> {
        return evaluateFlagInternal(key, defaultValue, ctx) { variants, variant ->
            if (variants?.has(variant) == true) variants.optBoolean(variant, defaultValue) else null
        }
    }

    override fun getStringEvaluation(
        key: String,
        defaultValue: String,
        ctx: EvaluationContext?
    ): ProviderEvaluation<String> {
        return evaluateFlagInternal(key, defaultValue, ctx) { variants, variant ->
            if (variants?.has(variant) == true) variants.optString(variant, defaultValue) else null
        }
    }

    override fun getIntegerEvaluation(
        key: String,
        defaultValue: Int,
        ctx: EvaluationContext?
    ): ProviderEvaluation<Int> {
        return evaluateFlagInternal(key, defaultValue, ctx) { variants, variant ->
            if (variants?.has(variant) == true) variants.optInt(variant, defaultValue) else null
        }
    }

    override fun getDoubleEvaluation(
        key: String,
        defaultValue: Double,
        ctx: EvaluationContext?
    ): ProviderEvaluation<Double> {
        return evaluateFlagInternal(key, defaultValue, ctx) { variants, variant ->
            if (variants?.has(variant) == true) variants.optDouble(variant, defaultValue) else null
        }
    }

    override fun getObjectEvaluation(
        key: String,
        defaultValue: Value,
        ctx: EvaluationContext?
    ): ProviderEvaluation<Value> {
        return ProviderEvaluation(defaultValue, reason = Reason.ERROR.name)
    }
}
