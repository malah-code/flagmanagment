import { apiClient } from './apiClient';

export interface SlackWebhookConfig {
  id?: string;
  environment_id: string;
  webhook_url: string;
  enabled: boolean;
}

export const slackApi = {
  getSlackConfig: async (envId: string): Promise<SlackWebhookConfig> => {
    return await apiClient.get<SlackWebhookConfig>(`/environments/${envId}/slack`);
  },

  saveSlackConfig: async (envId: string, webhook_url: string, enabled: boolean): Promise<SlackWebhookConfig> => {
    return await apiClient.post<SlackWebhookConfig>(`/environments/${envId}/slack`, {
      webhook_url,
      enabled,
    });
  },

  deleteSlackConfig: async (envId: string): Promise<void> => {
    await apiClient.delete(`/environments/${envId}/slack`);
  },
};
