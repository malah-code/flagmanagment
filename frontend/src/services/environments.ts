import { apiClient } from './apiClient';
import type { Environment } from '../types';

export interface CreateEnvironmentPayload {
  projectId: string;
  name: string;
}

export const environmentService = {
  async getByProject(projectId: string): Promise<Environment[]> {
    return apiClient.get<Environment[]>(`/projects/${projectId}/environments`);
  },

  async create(data: CreateEnvironmentPayload): Promise<Environment> {
    return apiClient.post<Environment>(`/projects/${data.projectId}/environments`, {
      name: data.name,
    });
  },

  async delete(id: string): Promise<void> {
    return apiClient.delete<void>(`/environments/${id}`);
  },
};
