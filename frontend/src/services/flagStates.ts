import { apiClient } from './apiClient';
import type { FlagState } from '../types';

export interface UpdateFlagStatePayload {
  isEnabled: boolean;
  rules?: Record<string, unknown>[];
}

export const flagStateService = {
  async getByEnvironment(environmentId: string): Promise<FlagState[]> {
    return apiClient.get<FlagState[]>(`/environments/${environmentId}/flag-states`);
  },

  async update(id: string, data: UpdateFlagStatePayload): Promise<FlagState> {
    return apiClient.put<FlagState>(`/flag-states/${id}`, data);
  },
};
