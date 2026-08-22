import { Provider, ResolutionDetails, EvaluationContext, JsonValue, StandardResolutionReasons, ErrorCode } from '@openfeature/react-sdk';
import { FlagClient } from './client';

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
      const flagState = this.client.getFlag(flagKey);
      if (!flagState) {
        return {
          value: defaultValue,
          reason: StandardResolutionReasons.DEFAULT,
          errorCode: ErrorCode.FLAG_NOT_FOUND,
        };
      }

      if (flagState.value !== undefined && flagState.value !== null) {
        const expectedType = typeof defaultValue;
        const actualType = typeof flagState.value;

        if (expectedType !== 'object' && actualType !== expectedType) {
          return {
            value: defaultValue,
            reason: StandardResolutionReasons.ERROR,
            errorCode: ErrorCode.TYPE_MISMATCH,
          };
        }
      }

      return {
        value: (flagState.value as unknown as T) ?? defaultValue,
        reason: flagState.reason,
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
