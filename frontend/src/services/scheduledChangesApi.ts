import { apiClient } from './apiClient';
import type {
  ScheduledChange,
  CreateScheduledChangeRequest,
  UpdateScheduledChangeRequest,
} from '../types/scheduledChange';

export const scheduledChangesApi = {
  create: async (envId: string, data: CreateScheduledChangeRequest): Promise<ScheduledChange> => {
    return await apiClient.post<ScheduledChange>(`/environments/${envId}/scheduled-changes`, data);
  },

  list: async (
    envId: string,
    status?: string,
  ): Promise<{ data: ScheduledChange[]; nextPageToken?: string }> => {
    const url = status
      ? `/environments/${envId}/scheduled-changes?status=${encodeURIComponent(status)}`
      : `/environments/${envId}/scheduled-changes`;
    return await apiClient.get<{ data: ScheduledChange[]; nextPageToken?: string }>(url);
  },

  getByID: async (id: string): Promise<ScheduledChange> => {
    return await apiClient.get<ScheduledChange>(`/scheduled-changes/${id}`);
  },

  update: async (id: string, data: UpdateScheduledChangeRequest): Promise<ScheduledChange> => {
    return await apiClient.patch<ScheduledChange>(`/scheduled-changes/${id}`, data);
  },

  cancel: async (id: string): Promise<ScheduledChange> => {
    return await apiClient.delete<ScheduledChange>(`/scheduled-changes/${id}`);
  },
};
