class EvaluationResult<T> {
  final T value;
  final String variant;
  final String reason;

  EvaluationResult({
    required this.value,
    required this.variant,
    required this.reason,
  });
}

class FlagEvaluator {
  /// Simple deterministic hashing for percentage bucketing (0..99).
  static int bucketUser(String targetingKey) {
    var hash = 0;
    for (var i = 0; i < targetingKey.length; i++) {
      hash = (31 * hash + targetingKey.codeUnitAt(i)) & 0x7FFFFFFF;
    }
    return hash % 100;
  }

  /// Evaluates a cached flag against user identity.
  static EvaluationResult<T> evaluate<T>({
    required Map<String, dynamic> flag,
    required T defaultValue,
    String? targetingKey,
  }) {
    final bool isEnabled = flag['enabled'] ?? false;
    final String defaultVariant = flag['defaultVariant'] ?? 'off';

    if (!isEnabled) {
      final val = flag['variants']?[defaultVariant] as T? ?? defaultValue;
      return EvaluationResult<T>(
        value: val,
        variant: defaultVariant,
        reason: 'DISABLED',
      );
    }

    if (targetingKey == null || targetingKey.isEmpty) {
      final val = flag['variants']?[defaultVariant] as T? ?? defaultValue;
      return EvaluationResult<T>(
        value: val,
        variant: defaultVariant,
        reason: 'DEFAULT',
      );
    }

    final rules = flag['rules'] as List<dynamic>?;
    if (rules != null) {
      for (final rule in rules) {
        final rollout = rule['rollout'] as Map<String, dynamic>?;
        if (rollout != null) {
          final bucket = bucketUser(targetingKey);
          double cumulative = 0.0;
          for (final entry in rollout.entries) {
            cumulative += (entry.value as num).toDouble();
            if (bucket < cumulative) {
              final val = flag['variants']?[entry.key] as T? ?? defaultValue;
              return EvaluationResult<T>(
                value: val,
                variant: entry.key,
                reason: 'TARGETING_MATCH',
              );
            }
          }
        }
      }
    }

    final val = flag['variants']?[defaultVariant] as T? ?? defaultValue;
    return EvaluationResult<T>(
      value: val,
      variant: defaultVariant,
      reason: 'DEFAULT',
    );
  }
}
