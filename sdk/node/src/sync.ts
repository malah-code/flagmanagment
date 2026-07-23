import { SDKOptions, RulesetSnapshot } from './types';
const EventSource = require('eventsource');

export class SyncService {
  private options: SDKOptions;
  private es?: any;
  private reconnectAttempt = 0;
  private reconnectTimeout?: NodeJS.Timeout;

  constructor(options: SDKOptions) {
    this.options = options;
  }

  public async fetchSnapshot(): Promise<RulesetSnapshot> {
    const endpoint = this.options.endpoint || 'http://localhost:8080';
    const url = `${endpoint}/api/v1/sdk/evaluate`; // Snapshot / REST evaluate URL

    const res = await fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${this.options.environmentToken}`,
      },
      body: JSON.stringify({}),
    });

    if (res.status === 401) {
      throw new Error('Unauthorized environment token');
    }

    if (!res.ok) {
      throw new Error(`Failed to fetch ruleset snapshot: HTTP ${res.status}`);
    }

    const data = (await res.json()) as any;
    return {
      version: data.version || '1.0.0',
      flags: data.flags || [],
    };
  }

  public startStreaming(onDelta: (version: string) => void, onError: (err: any) => void): () => void {
    const endpoint = this.options.endpoint || 'http://localhost:8080';
    const url = `${endpoint}/api/v1/sdk/stream`;

    const connect = () => {
      const es = new EventSource(url, {
        headers: {
          'Authorization': `Bearer ${this.options.environmentToken}`
        }
      });
      this.es = es;

      es.onopen = () => {
        this.reconnectAttempt = 0;
      };

      es.onmessage = async (e: any) => {
        try {
          // Whenever we receive a message (a ping or a delta), fetch the latest snapshot
          // In a real delta implementation, we'd apply patches. For now we resync on update.
          const snapshot = await this.fetchSnapshot();
          onDelta(snapshot.version);
        } catch (err) {
          onError(err);
        }
      };

      es.onerror = (err: any) => {
        onError(err);
        es.close();
        
        // Exponential backoff
        const delay = Math.min(1000 * Math.pow(2, this.reconnectAttempt), 30000);
        this.reconnectAttempt++;
        this.reconnectTimeout = setTimeout(connect, delay);
      };
    };

    connect();

    return () => {
      if (this.reconnectTimeout) clearTimeout(this.reconnectTimeout);
      if (this.es) this.es.close();
    };
  }
}
