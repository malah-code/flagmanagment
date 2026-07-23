import { apiClient } from './apiClient';
import type { FeatureFlag, FlagType } from '../types';

export interface CreateFlagPayload {
  projectId: string;
  key: string;
  description?: string;
  type: FlagType;
}

export const flagService = {
  async getByProject(projectId: string): Promise<FeatureFlag[]> {
    return apiClient.get<FeatureFlag[]>(`/projects/${projectId}/flags`);
  },

  async create(data: CreateFlagPayload): Promise<FeatureFlag> {
    return apiClient.post<FeatureFlag>(`/projects/${data.projectId}/flags`, {
      key: data.key,
      description: data.description,
      type: data.type,
    });
  },

  async delete(id: string): Promise<void> {
    return apiClient.delete<void>(`/flags/${id}`);
  },
};
