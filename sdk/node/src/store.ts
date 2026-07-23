import { FlagRule, RulesetSnapshot } from './types';

export class RuleStore {
  private flags: Map<string, FlagRule> = new Map();
  private version: string = '';

  public setSnapshot(snapshot: RulesetSnapshot): void {
    this.version = snapshot.version;
    this.flags.clear();
    for (const flag of snapshot.flags) {
      this.flags.set(flag.key, flag);
    }
  }

  public getFlag(key: string): FlagRule | undefined {
    return this.flags.get(key);
  }

  public getVersion(): string {
    return this.version;
  }

  public updateFlag(flag: FlagRule): void {
    this.flags.set(flag.key, flag);
  }
}
