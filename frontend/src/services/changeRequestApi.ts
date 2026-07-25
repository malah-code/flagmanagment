import { apiClient } from './apiClient';
import type { ChangeRequest } from '../types';

export const changeRequestApi = {
  listByEnvironment: async (envId: string, status?: string): Promise<ChangeRequest[]> => {
    const params = status ? `?status=${status}` : '';
    const response = await apiClient.get<{ data: ChangeRequest[] }>(`/environments/${envId}/change-requests${params}`);
    return response.data;
  },

  getById: async (id: string): Promise<ChangeRequest> => {
    return apiClient.get<ChangeRequest>(`/change-requests/${id}`);
  },

  approve: async (id: string, comment?: string): Promise<ChangeRequest> => {
    return apiClient.post<ChangeRequest>(`/change-requests/${id}/approve`, { comment });
  },

  reject: async (id: string, reason?: string): Promise<ChangeRequest> => {
    return apiClient.post<ChangeRequest>(`/change-requests/${id}/reject`, { reason });
  },
};
