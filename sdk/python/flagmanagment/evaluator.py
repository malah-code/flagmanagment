try:
    import mmh3
except ImportError:
    import zlib
    class mmh3:
        @staticmethod
        def hash(key: str, signed: bool = False) -> int:
            return zlib.crc32(key.encode('utf-8'))



class Evaluator:
    """Local flag evaluation engine using MurmurHash3 for deterministic bucketing."""

    @staticmethod
    def bucket_user(targeting_key: str) -> int:
        """Hash the targeting key with MurmurHash3 32-bit and return a bucket 0-99."""
        hash_val = mmh3.hash(targeting_key, signed=False)
        return hash_val % 100

    @staticmethod
    def evaluate_flag(flag: dict, targeting_key: str | None = None) -> tuple:
        """Evaluate a flag definition and return (value, variant, reason).

        Returns:
            Tuple of (value, variant_key, reason_string)
        """
        # 1. Check if flag is enabled
        if not flag.get("enabled", False):
            default_variant = flag.get("defaultVariant", "")
            value = Evaluator._get_variant_value(flag, default_variant)
            return value, default_variant, "DISABLED"

        # 2. No targeting key → default variant
        if not targeting_key:
            default_variant = flag.get("defaultVariant", "")
            value = Evaluator._get_variant_value(flag, default_variant)
            return value, default_variant, "DEFAULT"

        # 3. Evaluate rules with rollout percentages
        rules = flag.get("rules", [])
        for rule in rules:
            if not isinstance(rule, dict):
                continue
            rollout = rule.get("rollout")
            if rollout and isinstance(rollout, dict):
                variant = Evaluator._get_variant_by_rollout(targeting_key, rollout)
                if variant is not None:
                    value = Evaluator._get_variant_value(flag, variant)
                    return value, variant, "TARGETING_MATCH"

        # 4. Fallback to default
        default_variant = flag.get("defaultVariant", "")
        value = Evaluator._get_variant_value(flag, default_variant)
        return value, default_variant, "DEFAULT"

    @staticmethod
    def _get_variant_by_rollout(targeting_key: str, rollout: dict) -> str | None:
        """Determine variant based on MurmurHash3 bucketing and rollout percentages."""
        bucket = Evaluator.bucket_user(targeting_key)

        # Sort keys for deterministic iteration (matches Go, Java, .NET)
        current_threshold = 0
        for variant in sorted(rollout.keys()):
            percentage = rollout[variant]
            if not isinstance(percentage, (int, float)):
                continue
            current_threshold += int(percentage)
            if bucket < current_threshold:
                return variant

        return None

    @staticmethod
    def _get_variant_value(flag: dict, variant: str):
        """Look up the value for a variant key."""
        variants = flag.get("variants", {})
        return variants.get(variant)
