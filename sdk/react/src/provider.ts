import { Provider, ResolutionDetails, EvaluationContext, JsonValue, StandardResolutionReasons, ErrorCode } from '@openfeature/react-sdk';
import { FlagClient } from './client';
import { evaluateFlag } from './evaluator';

export class FlagManagmentWebProvider implements Provider {
  readonly metadata = {
    name: 'FlagManagment-React-Provider',
  };

  private client: FlagClient;

  constructor(client: FlagClient) {
    this.client = client;
  }

  private evaluateFlag<T>(flagKey: string, defaultValue: T, context: EvaluationContext): ResolutionDetails<T> {
    try {
      const flag = this.client.getFlag(flagKey);
      if (!flag) {
        return {
          value: defaultValue,
          reason: StandardResolutionReasons.DEFAULT,
          errorCode: ErrorCode.FLAG_NOT_FOUND,
        };
      }

      const targetingKey = context?.targetingKey ?? this.client.getContext()?.targetingKey;
      const result = evaluateFlag(flag, targetingKey);

      if (result.value !== undefined && result.value !== null) {
        const expectedType = typeof defaultValue;
        const actualType = typeof result.value;

        if (expectedType !== 'object' && actualType !== expectedType) {
          return {
            value: defaultValue,
            reason: StandardResolutionReasons.ERROR,
            errorCode: ErrorCode.TYPE_MISMATCH,
          };
        }
      }

      return {
        value: (result.value as unknown as T) ?? defaultValue,
        variant: result.variant,
        reason: result.reason,
      };
    } catch (e) {
      return {
        value: defaultValue,
        reason: StandardResolutionReasons.ERROR,
        errorCode: ErrorCode.GENERAL,
      };
    }
  }

  resolveBooleanEvaluation(flagKey: string, defaultValue: boolean, context: EvaluationContext): ResolutionDetails<boolean> {
    return this.evaluateFlag(flagKey, defaultValue, context);
  }

  resolveStringEvaluation(flagKey: string, defaultValue: string, context: EvaluationContext): ResolutionDetails<string> {
    return this.evaluateFlag(flagKey, defaultValue, context);
  }

  resolveNumberEvaluation(flagKey: string, defaultValue: number, context: EvaluationContext): ResolutionDetails<number> {
    return this.evaluateFlag(flagKey, defaultValue, context);
  }

  resolveObjectEvaluation<U extends JsonValue>(flagKey: string, defaultValue: U, context: EvaluationContext): ResolutionDetails<U> {
    return this.evaluateFlag<U>(flagKey, defaultValue, context);
  }
}
