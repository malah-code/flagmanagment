import { apiClient } from './apiClient';
import type { Environment } from '../types';

export interface CreateEnvironmentPayload {
  projectId: string;
  name: string;
}

export interface UpdateEnvironmentPayload {
  envId: string;
  name: string;
  isProtected: boolean;
  sdkSettings?: Record<string, any>;
}

export const environmentService = {
  async getByProject(projectId: string): Promise<Environment[]> {
    const res = await apiClient.get<{ data: Environment[] }>(`/projects/${projectId}/environments`);
    return res.data || [];
  },

  async create(data: CreateEnvironmentPayload): Promise<Environment> {
    return apiClient.post<Environment>(`/projects/${data.projectId}/environments`, {
      name: data.name,
    });
  },

  async update(data: UpdateEnvironmentPayload): Promise<Environment> {
    return apiClient.put<Environment>(`/environments/${data.envId}`, {
      name: data.name,
      isProtected: data.isProtected,
      sdkSettings: data.sdkSettings,
    });
  },

  async delete(id: string): Promise<void> {
    return apiClient.delete<void>(`/environments/${id}`);
  },

  async listServerKeys(projectId: string, envId: string) {
    const res = await apiClient.get<{ data: import('../types').ServerKey[] }>(
      `/projects/${projectId}/environments/${envId}/server-keys`
    );
    return res.data || [];
  },

  async createServerKey(projectId: string, envId: string, name: string) {
    return apiClient.post<import('../types').CreateServerKeyResponse>(
      `/projects/${projectId}/environments/${envId}/server-keys`,
      { name }
    );
  },

  async deleteServerKey(projectId: string, envId: string, keyId: string): Promise<void> {
    return apiClient.delete<void>(
      `/projects/${projectId}/environments/${envId}/server-keys/${keyId}`
    );
  },
};
