type Listener = () => void;

export class FlagClient {
  private apiKey: string;
  private streamUrl: string;
  private flags: Record<string, any> = {};
  private listeners: Set<Listener> = new Set();
  private eventSource: EventSource | null = null;
  private isConnecting: boolean = false;
  private context: Record<string, any> = {};

  constructor(apiKey: string, streamUrl: string) {
    this.apiKey = apiKey;
    this.streamUrl = streamUrl;
  }

  public setContext(context: Record<string, any>) {
    this.context = context;
    // Real-world implementation might re-fetch or re-evaluate based on new context.
  }

  public getContext() {
    return this.context;
  }

  public connect() {
    if (this.eventSource || this.isConnecting) return;
    this.isConnecting = true;

    // Passing authorization in standard EventSource is limited in browsers.
    // A production solution typically uses a polyfill or fetch-based SSE 
    // to pass the API Key in the Authorization header.
    // Assuming standard EventSource for MVP:
    this.eventSource = new EventSource(`${this.streamUrl}?apiKey=${this.apiKey}`);

    this.eventSource.addEventListener("bootstrap", (event: MessageEvent) => {
      try {
        const data = JSON.parse(event.data);
        if (data.flags) {
          this.flags = data.flags;
          this.notifyListeners();
        }
      } catch (e) {
        console.error("Failed to parse bootstrap data:", e);
      }
    });

    this.eventSource.addEventListener("flag_updated", (event: MessageEvent) => {
      try {
        const data = JSON.parse(event.data);
        if (data.flagKey && data.flag) {
          this.flags = {
            ...this.flags,
            [data.flagKey]: data.flag,
          };
          this.notifyListeners();
        }
      } catch (e) {
        console.error("Failed to parse flag update:", e);
      }
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
