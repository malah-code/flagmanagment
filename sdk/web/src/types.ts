export interface WebClientOptions {
  apiKey: string;
  apiUrl: string;
  streamUrl?: string;
  context?: Record<string, any>;
  enableLocalStorage?: boolean;
  storageKey?: string;
}

export interface FlagValue<T = any> {
  value: T;
  variant?: string;
  reason?: string;
}

export type FlagChangeListener = (flagKey: string, newValue: any, oldValue: any) => void;
export type ReadyListener = (flags: Record<string, FlagValue>) => void;
export type ErrorListener = (error: Error) => void;

