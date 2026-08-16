import { apiClient } from './apiClient';

export interface AuditLog {
  id: string;
  project_id?: string;
  environment_id?: string;
  actor_id: string;
  action: string;
  target_type: string;
  target_id: string;
  previous_state?: Record<string, any>;
  new_state?: Record<string, any>;
  actor_ip?: string;
  created_at: string;
}

export interface AuditLogsResponse {
  data: AuditLog[];
  next_page_token?: string;
}

export const auditService = {
  async getByProject(
    projectId: string,
    pageSize: number = 50,
    pageToken?: string,
  ): Promise<AuditLogsResponse> {
    const params = new URLSearchParams({ page_size: pageSize.toString() });
    if (pageToken) params.append('page_token', pageToken);
    return apiClient.get<AuditLogsResponse>(
      `/projects/${projectId}/audit-logs?${params.toString()}`,
    );
  },
};
