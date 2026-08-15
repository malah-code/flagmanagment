import { apiClient } from './apiClient';
import type { FlagState } from '../types';

export interface UpdateFlagStatePayload {
  isEnabled?: boolean;
  defaultVariation?: string;
  targetingRules?: { rules: any[] } | null;
  remoteConfig?: Record<string, any>;
  rolloutRules?: { rules: any[] } | null;
}

function mapApiFlagStateToUI(apiData: any): FlagState {
  return {
    ...apiData,
    id: apiData.id || `${apiData.featureFlagId}-${apiData.environmentId}`,
    flagId: apiData.featureFlagId,
    isEnabled: apiData.enabled,
    rules: apiData.rules || [],
  };
}

function mapUIFlagStateToApi(data: UpdateFlagStatePayload): any {
  return {
    ...data,
    enabled: data.isEnabled,
  };
}

export const flagStateService = {
  async getByEnvironment(projectId: string, environmentId: string): Promise<FlagState[]> {
    const res = await apiClient.get<{ data: any[] }>(`/projects/${projectId}/environments/${environmentId}/flag-states`);
    return (res.data || []).map(mapApiFlagStateToUI);
  },

  async update(projectId: string, environmentId: string, flagId: string, data: UpdateFlagStatePayload): Promise<FlagState> {
    const apiData = mapUIFlagStateToApi(data);
    const res = await apiClient.put<any>(`/projects/${projectId}/environments/${environmentId}/flags/${flagId}/state`, apiData);
    return mapApiFlagStateToUI(res);
  },

  async createOrUpdateByEnvAndFlag(projectId: string, environmentId: string, flagId: string, data: UpdateFlagStatePayload): Promise<FlagState> {
    const apiData = mapUIFlagStateToApi(data);
    const res = await apiClient.put<any>(`/projects/${projectId}/environments/${environmentId}/flags/${flagId}/state`, apiData);
    return mapApiFlagStateToUI(res);
  },
};
