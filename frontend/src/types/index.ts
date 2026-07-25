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
  apiKey?: string; // Only returned on creation
  createdAt: string;
  updatedAt: string;
}

export type FlagType = 'BOOLEAN' | 'STRING' | 'NUMBER' | 'JSON';

export interface FeatureFlag {
  id: string;
  projectId: string;
  key: string;
  description: string;
  type: FlagType;
  createdAt: string;
  updatedAt: string;
}

export interface FlagRule {
  id: string;
  property: string;
  operator: 'EQUALS' | 'CONTAINS' | 'GREATER_THAN' | 'LESS_THAN' | 'IN';
  value: string;
}

export interface FlagState {
  id: string;
  flagId: string;
  environmentId: string;
  isEnabled: boolean;
  rules: FlagRule[];
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
  proposedChanges: Record<string, any>;
  currentState?: Record<string, any>;
  createdBy: string;
  appliedBy?: string;
  createdAt: string;
  updatedAt: string;
}
