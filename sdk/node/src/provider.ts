import { Provider, ProviderMetadata, ResolutionDetails, EvaluationContext as OFContext, JsonValue } from '@openfeature/server-sdk';
import { FlagManagmentClient } from './client';
import { EvaluationContext } from './types';

export class FlagManagmentProvider implements Provider {
  readonly metadata: ProviderMetadata = {
    name: 'FlagManagmentProvider',
  };

  private client: FlagManagmentClient;

  constructor(client: FlagManagmentClient) {
    this.client = client;
  }

  private mapContext(ofCtx?: OFContext): EvaluationContext | undefined {
    if (!ofCtx) return undefined;
    const { targetingKey, ...attributes } = ofCtx;
    return {
      identity: targetingKey,
      attributes,
    };
  }

  async resolveBooleanEvaluation(flagKey: string, defaultValue: boolean, evalCtx?: OFContext): Promise<ResolutionDetails<boolean>> {
    const ctx = this.mapContext(evalCtx);
    const res = this.client.evaluate(flagKey, ctx, defaultValue);
    return {
      value: typeof res.value === 'boolean' ? res.value : defaultValue,
      reason: res.reason,
    };
  }

  async resolveStringEvaluation(flagKey: string, defaultValue: string, evalCtx?: OFContext): Promise<ResolutionDetails<string>> {
    const ctx = this.mapContext(evalCtx);
    const res = this.client.evaluate(flagKey, ctx, defaultValue);
    return {
      value: typeof res.value === 'string' ? res.value : defaultValue,
      reason: res.reason,
    };
  }

  async resolveNumberEvaluation(flagKey: string, defaultValue: number, evalCtx?: OFContext): Promise<ResolutionDetails<number>> {
    const ctx = this.mapContext(evalCtx);
    const res = this.client.evaluate(flagKey, ctx, defaultValue);
    return {
      value: typeof res.value === 'number' ? res.value : defaultValue,
      reason: res.reason,
    };
  }

  async resolveObjectEvaluation<T extends JsonValue>(flagKey: string, defaultValue: T, evalCtx?: OFContext): Promise<ResolutionDetails<T>> {
    const ctx = this.mapContext(evalCtx);
    const res = this.client.evaluate(flagKey, ctx, defaultValue);
    return {
      value: typeof res.value === 'object' && res.value !== null ? (res.value as T) : defaultValue,
      reason: res.reason,
    };
  }
}
