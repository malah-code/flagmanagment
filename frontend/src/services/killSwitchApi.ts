import { apiClient } from './apiClient';

export interface KillSwitchRule {
  id: string;
  flag_id: string;
  environment_id: string;
  alert_identifier: string;
  action: string;
  created_at: string;
}

export const killSwitchApi = {
  list: async (envId: string, flagId: string): Promise<KillSwitchRule[]> => {
    return await apiClient.get<KillSwitchRule[]>(
      `/environments/${envId}/flags/${flagId}/kill-switches`,
    );
  },

  create: async (
    envId: string,
    flagId: string,
    alertIdentifier: string,
  ): Promise<KillSwitchRule> => {
    return await apiClient.post<KillSwitchRule>(
      `/environments/${envId}/flags/${flagId}/kill-switches`,
      {
        alert_identifier: alertIdentifier,
      },
    );
  },

  delete: async (envId: string, flagId: string, id: string): Promise<void> => {
    await apiClient.delete(`/environments/${envId}/flags/${flagId}/kill-switches/${id}`);
  },
};
