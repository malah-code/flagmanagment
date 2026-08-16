// Shared Types for FlagManagement

export interface Project {
  id: string;
  name: string;
  description: string;
  createdAt: string;
  updatedAt: string;
}

export interface Environment {
  id: string;
  projectId: string;
  name: string;
  isProtected?: boolean;
  apiKey?: string;
  clientKey?: string;
  sdkSettings?: Record<string, any>;
  createdAt: string;
  updatedAt: string;
}

export interface ServerKey {
  id: string;
  environmentId: string;
  name: string;
  createdAt: string;
  lastUsedAt?: string;
}

export interface CreateServerKeyResponse {
  id: string;
  environmentId: string;
  name: string;
  key: string;
  createdAt: string;
}

export type FlagType = 'BOOLEAN' | 'MULTIVARIATE' | 'STRING' | 'NUMBER' | 'JSON';

export interface Variation {
  id: string;
  name: string;
  description?: string;
  value: unknown;
}

export interface FeatureFlag {
  id: string;
  projectId: string;
  key: string;
  name?: string;
  description: string;
  type: FlagType;
  variations?: Variation[];
  tags?: string[];
  parentFlagId?: string;
  createdAt: string;
  updatedAt: string;
}

export interface FlagRule {
  id: string;
  property: string;
  operator: 'EQUALS' | 'CONTAINS' | 'GREATER_THAN' | 'LESS_THAN' | 'IN';
  value: string;
}

export interface RolloutRule {
  variationId: string;
  percentage: number; // 0 to 10000 (basis points)
}

export type LifecycleState = 'ACTIVE' | 'STALE' | 'DEPRECATED' | 'ARCHIVED';

export interface StaleFlagPolicy {
  id: string;
  projectId: string;
  environmentId?: string;
  staleAfterDays: number;
  createdAt: string;
  updatedAt: string;
}

export interface FlagState {
  id: string;
  flagId: string;
  environmentId: string;
  isEnabled: boolean;
  lifecycleState?: LifecycleState;
  lastEvaluatedAt?: string;
  lastStateChangeAt?: string;
  rules: FlagRule[];
  targetingRules?: { rules: unknown[] };
  remoteConfig?: Record<string, unknown>;
  defaultVariation?: string;
  rolloutRules?: RolloutRule[];
  createdAt: string;
  updatedAt: string;
}

export type ChangeRequestStatus = 'PENDING' | 'APPROVED' | 'REJECTED' | 'APPLIED';

export interface ChangeRequest {
  id: string;
  projectId: string;
  environmentId: string;
  title: string;
  description: string;
  status: ChangeRequestStatus;
  proposedChanges: Record<string, unknown>;
  currentState?: Record<string, unknown>;
  createdBy: string;
  appliedBy?: string;
  createdAt: string;
  updatedAt: string;
}

export * from './scheduledChange';
