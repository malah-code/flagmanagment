import { apiClient } from './apiClient';
import type { ChangeRequest } from '../types';

function mapChangeRequest(apiData: any): ChangeRequest {
  return {
    id: apiData.id,
    projectId: apiData.project_id || apiData.projectId,
    environmentId: apiData.environment_id || apiData.environmentId,
    title: apiData.title,
    description: apiData.description,
    status: apiData.status,
    proposedChanges: apiData.proposed_changes || apiData.proposedChanges || {},
    currentState: apiData.current_state || apiData.currentState || {},
    createdBy: apiData.created_by || apiData.createdBy,
    appliedBy: apiData.applied_by || apiData.appliedBy,
    createdAt: apiData.created_at || apiData.createdAt,
    updatedAt: apiData.updated_at || apiData.updatedAt,
  };
}

export const changeRequestApi = {
  listByEnvironment: async (envId: string, status?: string): Promise<ChangeRequest[]> => {
    const params = status ? `?status=${status}` : '';
    const response = await apiClient.get<{ data: any[] }>(
      `/environments/${envId}/change-requests${params}`,
    );
    return (response.data || []).map(mapChangeRequest);
  },

  getById: async (id: string): Promise<ChangeRequest> => {
    const response = await apiClient.get<any>(`/change-requests/${id}`);
    return mapChangeRequest(response);
  },

  approve: async (id: string, comment?: string): Promise<ChangeRequest> => {
    const response = await apiClient.post<any>(`/change-requests/${id}/approve`, { comment });
    return mapChangeRequest(response);
  },

  reject: async (id: string, reason?: string): Promise<ChangeRequest> => {
    const response = await apiClient.post<any>(`/change-requests/${id}/reject`, { reason });
    return mapChangeRequest(response);
  },
};
