import { SDKOptions, EvaluationContext, EvaluationResult } from './types';
import { RuleStore } from './store';
import { SyncService } from './sync';
import { evaluateFlag, hashPII } from './evaluator';

export class FlagManagmentClient {
  private store: RuleStore;
  private sync: SyncService;
  private stopStream?: () => void;

  constructor(options: SDKOptions) {
    this.store = new RuleStore();
    this.sync = new SyncService(options);
  }

  public async init(): Promise<void> {
    const snapshot = await this.sync.fetchSnapshot();
    this.store.setSnapshot(snapshot);

    this.stopStream = this.sync.startStreaming(
      async () => {
        const latest = await this.sync.fetchSnapshot();
        this.store.setSnapshot(latest);
      },
      (err) => {
        console.warn('[FlagManagment] Sync stream error:', err.message);
      }
    );
  }

  public close(): void {
    if (this.stopStream) {
      this.stopStream();
    }
  }

  public evaluate(flagKey: string, context?: EvaluationContext, defaultValue: any = false): EvaluationResult {
    const flag = this.store.getFlag(flagKey);
    if (!flag) {
      return {
        value: defaultValue,
        reason: 'DEFAULT',
      };
    }

    const hashedCtx = hashPII(context);
    return evaluateFlag(flag, hashedCtx);
  }

  public getBooleanValue(flagKey: string, defaultValue: boolean, context?: EvaluationContext): boolean {
    const res = this.evaluate(flagKey, context, defaultValue);
    return typeof res.value === 'boolean' ? res.value : defaultValue;
  }

  public getStringValue(flagKey: string, defaultValue: string, context?: EvaluationContext): string {
    const res = this.evaluate(flagKey, context, defaultValue);
    return typeof res.value === 'string' ? res.value : defaultValue;
  }
}
