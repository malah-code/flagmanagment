import { FlagWebClient } from './client';

export interface OpenFeatureResolution<T> {
  value: T;
  variant?: string;
  reason?: string;
  errorCode?: string;
  errorMessage?: string;
}

/**
 * OpenFeature Client Provider for Web/Browser applications.
 */
export class OpenFeatureWebProvider {
  readonly metadata = {
    name: 'FlagManagment-Web-Provider',
  };

  private client: FlagWebClient;

  constructor(client: FlagWebClient) {
    this.client = client;
  }

  async initialize(context?: Record<string, any>): Promise<void> {
    if (context) {
      await this.client.setContext(context);
    } else {
      await this.client.connect();
    }
  }

  onContextChange(newContext: Record<string, any>): Promise<void> {
    return this.client.setContext(newContext);
  }

  resolveBooleanEvaluation(flagKey: string, defaultValue: boolean): OpenFeatureResolution<boolean> {
    const val = this.client.getBooleanValue(flagKey, defaultValue);
    const raw = this.client.getFlags()[flagKey];
    return {
      value: val,
      variant: raw?.variant,
      reason: raw?.reason || 'TARGETING_MATCH',
    };
  }

  resolveStringEvaluation(flagKey: string, defaultValue: string): OpenFeatureResolution<string> {
    const val = this.client.getStringValue(flagKey, defaultValue);
    const raw = this.client.getFlags()[flagKey];
    return {
      value: val,
      variant: raw?.variant,
      reason: raw?.reason || 'TARGETING_MATCH',
    };
  }

  resolveNumberEvaluation(flagKey: string, defaultValue: number): OpenFeatureResolution<number> {
    const val = this.client.getNumberValue(flagKey, defaultValue);
    const raw = this.client.getFlags()[flagKey];
    return {
      value: val,
      variant: raw?.variant,
      reason: raw?.reason || 'TARGETING_MATCH',
    };
  }

  resolveObjectEvaluation<T = any>(flagKey: string, defaultValue: T): OpenFeatureResolution<T> {
    const val = this.client.getObjectValue<T>(flagKey, defaultValue);
    const raw = this.client.getFlags()[flagKey];
    return {
      value: val,
      variant: raw?.variant,
      reason: raw?.reason || 'TARGETING_MATCH',
    };
  }
}

