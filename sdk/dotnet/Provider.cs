using OpenFeature.Model;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text.Json;
using System.Threading.Tasks;

namespace FlagManagment.Sdk
{
    /// <summary>
    /// OpenFeature provider for FlagManagment. Evaluates flags locally
    /// using the in-memory cache maintained by Client and MurmurHash3
    /// bucketing via Evaluator.
    /// </summary>
    public class Provider : OpenFeature.FeatureProvider
    {
        private readonly Client _client;

        public Provider(Client client)
        {
            _client = client;
        }

        public override Metadata GetMetadata()
        {
            return new Metadata("FlagManagment-.NET-Provider");
        }

        public override Task<ResolutionDetails<bool>> ResolveBooleanValueAsync(
            string flagKey, bool defaultValue, EvaluationContext context = null)
        {
            var result = Evaluate(flagKey, context);
            if (result == null)
            {
                return Task.FromResult(new ResolutionDetails<bool>(
                    flagKey, defaultValue, ErrorType.FlagNotFound, reason: "ERROR"));
            }

            if (result.Value.Value is bool boolVal)
            {
                return Task.FromResult(new ResolutionDetails<bool>(
                    flagKey, boolVal, variant: result.Value.Variant, reason: result.Value.Reason));
            }
            return Task.FromResult(new ResolutionDetails<bool>(
                flagKey, defaultValue, reason: "ERROR"));
        }

        public override Task<ResolutionDetails<string>> ResolveStringValueAsync(
            string flagKey, string defaultValue, EvaluationContext context = null)
        {
            var result = Evaluate(flagKey, context);
            if (result == null)
            {
                return Task.FromResult(new ResolutionDetails<string>(flagKey, defaultValue, reason: "ERROR"));
            }

            if (result.Value.Value is string strVal)
            {
                return Task.FromResult(new ResolutionDetails<string>(
                    flagKey, strVal, variant: result.Value.Variant, reason: result.Value.Reason));
            }
            return Task.FromResult(new ResolutionDetails<string>(flagKey, defaultValue, reason: "ERROR"));
        }

        public override Task<ResolutionDetails<int>> ResolveIntegerValueAsync(
            string flagKey, int defaultValue, EvaluationContext context = null)
        {
            var result = Evaluate(flagKey, context);
            if (result?.Value.Value is JsonElement je && je.ValueKind == JsonValueKind.Number)
            {
                return Task.FromResult(new ResolutionDetails<int>(
                    flagKey, je.GetInt32(), variant: result.Value.Variant, reason: result.Value.Reason));
            }
            return Task.FromResult(new ResolutionDetails<int>(flagKey, defaultValue, reason: "ERROR"));
        }

        public override Task<ResolutionDetails<double>> ResolveDoubleValueAsync(
            string flagKey, double defaultValue, EvaluationContext context = null)
        {
            var result = Evaluate(flagKey, context);
            if (result?.Value.Value is JsonElement je && je.ValueKind == JsonValueKind.Number)
            {
                return Task.FromResult(new ResolutionDetails<double>(
                    flagKey, je.GetDouble(), variant: result.Value.Variant, reason: result.Value.Reason));
            }
            return Task.FromResult(new ResolutionDetails<double>(flagKey, defaultValue, reason: "ERROR"));
        }

        public override Task<ResolutionDetails<Value>> ResolveStructureValueAsync(
            string flagKey, Value defaultValue, EvaluationContext context = null)
        {
            return Task.FromResult(new ResolutionDetails<Value>(flagKey, defaultValue, reason: "DEFAULT"));
        }

        // --- Internal evaluation ---

        private record struct EvalResult(object Value, string Variant, string Reason);

        private EvalResult? Evaluate(string flagKey, EvaluationContext context)
        {
            var flagJson = _client.GetFlag(flagKey);
            if (flagJson == null) return null;

            var flag = flagJson.Value;

            // Check enabled
            bool enabled = flag.TryGetProperty("enabled", out var enabledEl) && enabledEl.GetBoolean();
            if (!enabled)
            {
                string dv = GetStringProp(flag, "defaultVariant");
                return new EvalResult(GetVariantValue(flag, dv), dv, "DISABLED");
            }

            // Extract targeting key
            string targetingKey = context?.TargetingKey;

            if (string.IsNullOrEmpty(targetingKey))
            {
                string dv = GetStringProp(flag, "defaultVariant");
                return new EvalResult(GetVariantValue(flag, dv), dv, "DEFAULT");
            }

            // Evaluate rules
            if (flag.TryGetProperty("rules", out var rulesEl) && rulesEl.ValueKind == JsonValueKind.Array)
            {
                foreach (var rule in rulesEl.EnumerateArray())
                {
                    if (rule.TryGetProperty("rollout", out var rolloutEl))
                    {
                        string variant = GetVariantByRollout(targetingKey, rolloutEl);
                        if (variant != null)
                        {
                            return new EvalResult(GetVariantValue(flag, variant), variant, "TARGETING_MATCH");
                        }
                    }
                }
            }

            // Fallback
            string defaultVariant = GetStringProp(flag, "defaultVariant");
            return new EvalResult(GetVariantValue(flag, defaultVariant), defaultVariant, "DEFAULT");
        }

        private string GetVariantByRollout(string targetingKey, JsonElement rollout)
        {
            int bucket = Evaluator.BucketUser(targetingKey);

            // Sort keys for deterministic iteration
            var sorted = rollout.EnumerateObject()
                .OrderBy(p => p.Name)
                .ToList();

            int currentThreshold = 0;
            foreach (var prop in sorted)
            {
                if (prop.Value.TryGetInt32(out int percentage))
                {
                    currentThreshold += percentage;
                    if (bucket < currentThreshold) return prop.Name;
                }
            }
            return null;
        }

        private object GetVariantValue(JsonElement flag, string variant)
        {
            if (string.IsNullOrEmpty(variant)) return null;
            if (!flag.TryGetProperty("variants", out var variants)) return null;
            if (!variants.TryGetProperty(variant, out var val)) return null;

            return val.ValueKind switch
            {
                JsonValueKind.True => (object)true,
                JsonValueKind.False => false,
                JsonValueKind.String => val.GetString(),
                JsonValueKind.Number => val,
                _ => val
            };
        }

        private static string GetStringProp(JsonElement el, string prop)
        {
            return el.TryGetProperty(prop, out var val) ? val.GetString() ?? "" : "";
        }
    }
}
