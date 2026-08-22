import { SDKOptions } from './types';

export interface SDKEvaluationMetric {
  flagKey: string;
  timestamp: string;
}

export class MetricsSync {
  private options: SDKOptions;
  private buffer: SDKEvaluationMetric[] = [];
  private syncInterval?: NodeJS.Timeout;

  constructor(options: SDKOptions) {
    this.options = options;
  }

  public start(intervalMs: number = 15000) {
    if (this.syncInterval) {
      clearInterval(this.syncInterval);
    }
    this.syncInterval = setInterval(() => {
      this.flush();
    }, intervalMs);
  }

  public stop() {
    if (this.syncInterval) {
      clearInterval(this.syncInterval);
      this.syncInterval = undefined;
    }
    this.flush();
  }

  public recordEvaluation(flagKey: string) {
    this.buffer.push({
      flagKey,
      timestamp: new Date().toISOString(),
    });
  }

  private async flush() {
    if (this.buffer.length === 0) return;

    // Capture the current buffer and clear it
    const payload = { evaluations: [...this.buffer] };
    this.buffer = [];

    const endpoint = this.options.endpoint || 'http://localhost:8080';
    const url = `${endpoint}/api/v1/sdk/metrics`;

    try {
      const res = await fetch(url, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${this.options.environmentToken}`,
        },
        body: JSON.stringify(payload),
      });

      if (!res.ok) {
        console.warn(`[FlagManagment] Failed to sync metrics: HTTP ${res.status}`);
      }
    } catch (err) {
      console.warn(`[FlagManagment] Failed to sync metrics: ${(err as Error).message}`);
    }
  }
}
