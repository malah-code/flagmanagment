export interface FlagRule {
  key: string;
  type: string;
  enabled: boolean;
  defaultVariation: string;
  targetingRulesJson?: any;
}

export interface RulesetSnapshot {
  version: string;
  flags: FlagRule[];
}

export interface EvaluationContext {
  identity?: string;
  attributes?: Record<string, any>;
}

export interface EvaluationResult {
  value: boolean | string | number | Record<string, any>;
  reason: string;
  variation?: string;
}

export interface SDKOptions {
  environmentToken: string;
  endpoint?: string;
}
