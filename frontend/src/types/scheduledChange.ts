export type ScheduledChangeTargetType = 'FLAG' | 'CHANGE_REQUEST';
export type ScheduledChangeAction = 'ENABLE' | 'DISABLE' | 'APPLY';
export type ScheduledChangeStatus = 'PENDING' | 'EXECUTED' | 'CANCELLED';

export interface ScheduledChange {
  id: string;
  project_id: string;
  environment_id: string;
  target_type: ScheduledChangeTargetType;
  target_id: string;
  action: ScheduledChangeAction;
  scheduled_for: string; // ISO-8601 UTC
  status: ScheduledChangeStatus;
  created_by: string;
  executed_at?: string | null;
  cancelled_at?: string | null;
  created_at: string;
  updated_at: string;
}

export interface CreateScheduledChangeRequest {
  target_type: ScheduledChangeTargetType;
  target_id: string;
  action: ScheduledChangeAction;
  scheduled_for: string; // ISO-8601 UTC
}

export interface UpdateScheduledChangeRequest {
  scheduled_for: string; // ISO-8601 UTC
}
