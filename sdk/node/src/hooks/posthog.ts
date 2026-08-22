import { Hook, HookContext, EvaluationDetails } from '@openfeature/server-sdk';

export interface PostHogHookOptions {
  apiKey: string;
  host?: string;
  flushIntervalMs?: number;
  batchSize?: number;
}

export class PostHogHook implements Hook {
  public name = 'PostHogHook';
  private apiKey: string;
  private host: string;
  private flushIntervalMs: number;
  private batchSize: number;
  
  private buffer: any[] = [];
  private syncInterval?: NodeJS.Timeout;

  constructor(options: PostHogHookOptions) {
    this.apiKey = options.apiKey;
    this.host = options.host || 'https://us.i.posthog.com';
    this.flushIntervalMs = options.flushIntervalMs || 15000;
    this.batchSize = options.batchSize || 100;

    this.start();
  }

  private start() {
    if (this.syncInterval) {
      clearInterval(this.syncInterval);
    }
    this.syncInterval = setInterval(() => {
      this.flush();
    }, this.flushIntervalMs);
  }

  public close() {
    if (this.syncInterval) {
      clearInterval(this.syncInterval);
    }
    this.flush();
  }

  after(hookContext: HookContext, evaluationDetails: EvaluationDetails<any>): void {
    // Only capture events if the flag was successfully evaluated and we have an identity
    const identity = hookContext.context.targetingKey || hookContext.context.identity;
    
    // Ignore default fallback events if desired, but here we track them as "DEFAULT"
    const isDefault = evaluationDetails.reason === 'DEFAULT';

    const event = {
      event: '$feature_flag_called',
      properties: {
        distinct_id: identity || 'anonymous',
        $feature_flag: hookContext.flagKey,
        $feature_flag_response: evaluationDetails.value,
        $feature_flag_reason: evaluationDetails.reason,
        timestamp: new Date().toISOString(),
      },
    };

    this.buffer.push(event);

    if (this.buffer.length >= this.batchSize) {
      this.flush();
    }
  }

  error(hookContext: HookContext, error: Error): void {
    const identity = hookContext.context.targetingKey || hookContext.context.identity;
    const event = {
      event: '$feature_flag_error',
      properties: {
        distinct_id: identity || 'anonymous',
        $feature_flag: hookContext.flagKey,
        error_message: error.message,
        timestamp: new Date().toISOString(),
      },
    };
    this.buffer.push(event);
  }

  private async flush() {
    if (this.buffer.length === 0) return;

    const payload = {
      api_key: this.apiKey,
      batch: [...this.buffer],
    };
    this.buffer = [];

    try {
      const res = await fetch(`${this.host}/capture/`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(payload),
      });

      if (!res.ok) {
        console.warn(`[PostHogHook] Failed to flush events: HTTP ${res.status}`);
      }
    } catch (err) {
      console.warn(`[PostHogHook] Failed to flush events: ${(err as Error).message}`);
    }
  }
}
