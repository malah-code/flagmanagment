type Listener = () => void;

export class FlagClient {
  private apiKey: string;
  private apiUrl: string;
  private streamUrl?: string;
  private flags: Record<string, any> = {};
  private listeners: Set<Listener> = new Set();
  private eventSource: EventSource | null = null;
  private isConnecting: boolean = false;
  private context: Record<string, any> = {};

  constructor(apiKey: string, apiUrl: string, streamUrl?: string) {
    this.apiKey = apiKey;
    this.apiUrl = apiUrl;
    this.streamUrl = streamUrl;
  }

  public async setContext(context: Record<string, any>) {
    this.context = context;
    await this.fetchFlags();
  }

  public getContext() {
    return this.context;
  }

  private async fetchFlags() {
    try {
      const url = `${this.apiUrl}/api/v1/client/evaluate`;
      const res = await fetch(url, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${this.apiKey}`,
        },
        body: JSON.stringify({ context: this.context }),
      });

      if (!res.ok) {
        console.error(`[FlagManagment] Failed to fetch flags: HTTP ${res.status}`);
        return;
      }

      const data = await res.json();
      this.flags = data || {};
      this.notifyListeners();
    } catch (e) {
      console.error(`[FlagManagment] Error fetching flags:`, e);
    }
  }

  public async connect() {
    await this.fetchFlags();

    if (!this.streamUrl) return;

    if (this.eventSource || this.isConnecting) return;
    this.isConnecting = true;

    // Standard EventSource for MVP
    this.eventSource = new EventSource(`${this.streamUrl}?apiKey=${this.apiKey}`);

    this.eventSource.addEventListener("bootstrap", () => {
      this.fetchFlags();
    });

    this.eventSource.addEventListener("flag_updated", () => {
      // Whenever a flag is updated on the server, we re-fetch for this context
      this.fetchFlags();
    });

    this.eventSource.onerror = () => {
      this.eventSource?.close();
      this.eventSource = null;
      this.isConnecting = false;
      // Reconnect logic
      setTimeout(() => this.connect(), 5000);
    };
  }

  public disconnect() {
    if (this.eventSource) {
      this.eventSource.close();
      this.eventSource = null;
    }
    this.isConnecting = false;
  }

  public getFlag(key: string) {
    return this.flags[key];
  }

  public getFlags() {
    return this.flags;
  }

  public subscribe(listener: Listener): () => void {
    this.listeners.add(listener);
    return () => {
      this.listeners.delete(listener);
    };
  }

  private notifyListeners() {
    this.listeners.forEach((listener) => listener());
  }
}
