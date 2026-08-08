import { apiClient } from './apiClient';
import type { FeatureFlag, FlagType, Variation } from '../types';

export interface CreateFlagPayload {
  projectId: string;
  key: string;
  description?: string;
  type: FlagType;
  variations?: Variation[];
  parentFlagId?: string;
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
      variations: data.variations,
      parentFlagId: data.parentFlagId,
    });
  },

  async delete(id: string): Promise<void> {
    return apiClient.delete<void>(`/flags/${id}`);
  },

  async promote(projectId: string, flagId: string, sourceEnvId: string, targetEnvId: string): Promise<any> {
    return apiClient.post<any>(`/projects/${projectId}/flags/${flagId}/promote`, {
      source_env_id: sourceEnvId,
      target_env_id: targetEnvId,
    });
  },
};
