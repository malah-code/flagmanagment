package com.flagmanagment.sdk

import dev.openfeature.sdk.*
import org.json.JSONObject

class FlagManagmentProvider(private val client: FlagClient) : FeatureProvider {
    override val hooks: List<Hook<*>> = listOf()
    override val metadata: Metadata = FlagManagmentMetadata()

    private class FlagManagmentMetadata : Metadata {
        override val name: String = "FlagManagment-Android-Provider"
    }

    override fun getBooleanEvaluation(
        key: String,
        defaultValue: Boolean,
        ctx: EvaluationContext?
    ): ProviderEvaluation<Boolean> {
        val flag = client.getFlag(key) as? JSONObject ?: return ProviderEvaluation(defaultValue, reason = Reason.ERROR.name)
        
        val defaultVariant = flag.optString("defaultVariant")
        val variants = flag.optJSONObject("variants")
        val value = variants?.optBoolean(defaultVariant, defaultValue) ?: defaultValue
        
        return ProviderEvaluation(value, reason = Reason.STATIC.name)
    }

    override fun getStringEvaluation(
        key: String,
        defaultValue: String,
        ctx: EvaluationContext?
    ): ProviderEvaluation<String> {
        val flag = client.getFlag(key) as? JSONObject ?: return ProviderEvaluation(defaultValue, reason = Reason.ERROR.name)
        
        val defaultVariant = flag.optString("defaultVariant")
        val variants = flag.optJSONObject("variants")
        val value = variants?.optString(defaultVariant, defaultValue) ?: defaultValue
        
        return ProviderEvaluation(value, reason = Reason.STATIC.name)
    }

    override fun getIntegerEvaluation(
        key: String,
        defaultValue: Int,
        ctx: EvaluationContext?
    ): ProviderEvaluation<Int> {
        val flag = client.getFlag(key) as? JSONObject ?: return ProviderEvaluation(defaultValue, reason = Reason.ERROR.name)
        
        val defaultVariant = flag.optString("defaultVariant")
        val variants = flag.optJSONObject("variants")
        val value = variants?.optInt(defaultVariant, defaultValue) ?: defaultValue
        
        return ProviderEvaluation(value, reason = Reason.STATIC.name)
    }

    override fun getDoubleEvaluation(
        key: String,
        defaultValue: Double,
        ctx: EvaluationContext?
    ): ProviderEvaluation<Double> {
        val flag = client.getFlag(key) as? JSONObject ?: return ProviderEvaluation(defaultValue, reason = Reason.ERROR.name)
        
        val defaultVariant = flag.optString("defaultVariant")
        val variants = flag.optJSONObject("variants")
        val value = variants?.optDouble(defaultVariant, defaultValue) ?: defaultValue
        
        return ProviderEvaluation(value, reason = Reason.STATIC.name)
    }

    override fun getObjectEvaluation(
        key: String,
        defaultValue: Value,
        ctx: EvaluationContext?
    ): ProviderEvaluation<Value> {
        return ProviderEvaluation(defaultValue, reason = Reason.ERROR.name)
    }
}
