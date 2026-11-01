import { WebClientOptions, FlagValue, FlagChangeListener, ReadyListener, ErrorListener } from './types';

export class FlagWebClient {
  private apiKey: string;
  private apiUrl: string;
  private streamUrl?: string;
  private context: Record<string, any>;
  private enableLocalStorage: boolean;
  private storageKey: string;

  private flags: Record<string, FlagValue> = {};
  private eventSource: any = null;
  private isReady = false;

  private changeListeners: Set<FlagChangeListener> = new Set();
  private readyListeners: Set<ReadyListener> = new Set();
  private errorListeners: Set<ErrorListener> = new Set();

  constructor(options: WebClientOptions) {
    this.apiKey = options.apiKey;
    this.apiUrl = options.apiUrl.replace(/\/+$/, '');
    this.streamUrl = options.streamUrl;
    this.context = options.context || {};
    this.enableLocalStorage = options.enableLocalStorage ?? true;
    this.storageKey = options.storageKey || `flagmanagment_cache_${this.apiKey.slice(-8)}`;

    this.loadFromStorage();
  }

  /**
   * Initializes the client, evaluates flags for the current context,
   * and optionally starts real-time SSE streaming.
   */
  public async connect(): Promise<Record<string, FlagValue>> {
    try {
      await this.fetchEvaluations();
      this.isReady = true;
      this.notifyReady();

      if (typeof window !== 'undefined' && 'EventSource' in window) {
        this.startStreaming();
      }

      return this.flags;
    } catch (err) {
      const error = err instanceof Error ? err : new Error(String(err));
      this.notifyError(error);
      return this.flags;
    }
  }

  /**
   * Disconnects active SSE streams.
   */
  public disconnect(): void {
    if (this.eventSource) {
      this.eventSource.close();
      this.eventSource = null;
    }
  }

  /**
   * Updates user targeting context and triggers re-evaluation.
   */
  public async setContext(newContext: Record<string, any>): Promise<void> {
    this.context = { ...this.context, ...newContext };
    await this.fetchEvaluations();
  }

  /**
   * Returns boolean value for a flag.
   */
  public getBooleanValue(flagKey: string, defaultValue: boolean): boolean {
    const flag = this.flags[flagKey];
    if (flag && typeof flag.value === 'boolean') {
      return flag.value;
    }
    return defaultValue;
  }

  /**
   * Returns string value for a flag.
   */
  public getStringValue(flagKey: string, defaultValue: string): string {
    const flag = this.flags[flagKey];
    if (flag && typeof flag.value === 'string') {
      return flag.value;
    }
    return defaultValue;
  }

  /**
   * Returns number value for a flag.
   */
  public getNumberValue(flagKey: string, defaultValue: number): number {
    const flag = this.flags[flagKey];
    if (flag && typeof flag.value === 'number') {
      return flag.value;
    }
    return defaultValue;
  }

  /**
   * Returns JSON object value for a flag.
   */
  public getObjectValue<T = any>(flagKey: string, defaultValue: T): T {
    const flag = this.flags[flagKey];
    if (flag && typeof flag.value === 'object' && flag.value !== null) {
      return flag.value as T;
    }
    return defaultValue;
  }

  /**
   * Returns snapshot of all currently evaluated flags.
   */
  public getFlags(): Record<string, FlagValue> {
    return { ...this.flags };
  }

  /**
   * Event registration.
   */
  public on(event: 'ready', callback: ReadyListener): () => void;
  public on(event: 'change', callback: FlagChangeListener): () => void;
  public on(event: 'error', callback: ErrorListener): () => void;
  public on(event: 'ready' | 'change' | 'error', callback: any): () => void {
    if (event === 'ready') {
      this.readyListeners.add(callback);
      if (this.isReady) callback(this.flags);
      return () => this.readyListeners.delete(callback);
    }
    if (event === 'change') {
      this.changeListeners.add(callback);
      return () => this.changeListeners.delete(callback);
    }
    if (event === 'error') {
      this.errorListeners.add(callback);
      return () => this.errorListeners.delete(callback);
    }
    return () => {};
  }

  private async fetchEvaluations(): Promise<void> {
    const url = `${this.apiUrl}/api/v1/sdk/evaluations`;
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      'X-Client-Key': this.apiKey,
    };

    const res = await fetch(url, {
      method: 'POST',
      headers,
      body: JSON.stringify({ context: this.context }),
    });

    if (!res.ok) {
      throw new Error(`Failed to fetch flag evaluations: ${res.statusText}`);
    }

    const data = await res.json();
    const newFlags: Record<string, FlagValue> = data.flags || data;
    this.updateFlags(newFlags);
  }

  private updateFlags(newFlags: Record<string, FlagValue>): void {
    const oldFlags = { ...this.flags };
    this.flags = newFlags;
    this.saveToStorage();

    // Trigger change events
    for (const [key, newFlag] of Object.entries(newFlags)) {
      const oldFlag = oldFlags[key];
      if (!oldFlag || oldFlag.value !== newFlag.value) {
        this.notifyChange(key, newFlag.value, oldFlag?.value);
      }
    }
  }

  private startStreaming(): void {
    if (this.eventSource) {
      this.eventSource.close();
    }

    const streamEndpoint = this.streamUrl || `${this.apiUrl}/api/v1/sdk/stream`;
    try {
      this.eventSource = new EventSource(`${streamEndpoint}?token=${this.apiKey}`);

      this.eventSource.onmessage = (event: MessageEvent) => {
        try {
          const payload = JSON.parse(event.data);
          if (payload.flagKey && payload.flag) {
            const oldVal = this.flags[payload.flagKey]?.value;
            this.flags[payload.flagKey] = {
              value: payload.flag.value ?? payload.flag,
              variant: payload.flag.variant,
              reason: 'STREAM_UPDATE',
            };
            this.saveToStorage();
            this.notifyChange(payload.flagKey, this.flags[payload.flagKey].value, oldVal);
          }
        } catch (e) {
          // ignore parse errors on keepalive ping
        }
      };

      this.eventSource.onerror = () => {
        // EventSource will automatically attempt reconnection
      };
    } catch (e) {
      // streaming unavailable
    }
  }

  private loadFromStorage(): void {
    if (!this.enableLocalStorage || typeof localStorage === 'undefined') return;
    try {
      const raw = localStorage.getItem(this.storageKey);
      if (raw) {
        this.flags = JSON.parse(raw);
      }
    } catch {
      // fallback
    }
  }

  private saveToStorage(): void {
    if (!this.enableLocalStorage || typeof localStorage === 'undefined') return;
    try {
      localStorage.setItem(this.storageKey, JSON.stringify(this.flags));
    } catch {
      // storage quota or private browsing fallback
    }
  }

  private notifyReady(): void {
    this.readyListeners.forEach((fn) => fn(this.flags));
  }

  private notifyChange(key: string, newVal: any, oldVal: any): void {
    this.changeListeners.forEach((fn) => fn(key, newVal, oldVal));
  }

  private notifyError(err: Error): void {
    this.errorListeners.forEach((fn) => fn(err));
  }
}
