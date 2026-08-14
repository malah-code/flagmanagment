import { apiClient } from './apiClient';
import type { FeatureFlag, FlagType, Variation } from '../types';

export interface CreateFlagPayload {
  projectId: string;
  key: string;
  name?: string;
  description?: string;
  type: FlagType;
  variations?: Variation[];
  tags?: string[];
  parentFlagId?: string;
  enabledByDefault?: boolean;
}

export const flagService = {
  async getByProject(projectId: string): Promise<FeatureFlag[]> {
    const res = await apiClient.get<{ data: FeatureFlag[] }>(`/projects/${projectId}/flags`);
    return res.data || [];
  },

  async create(data: CreateFlagPayload): Promise<FeatureFlag> {
    return apiClient.post<FeatureFlag>(`/projects/${data.projectId}/flags`, {
      key: data.key,
      name: data.name || data.key,
      description: data.description,
      type: data.type,
      enabledByDefault: data.enabledByDefault,
      variations: data.variations,
      tags: data.tags,
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
