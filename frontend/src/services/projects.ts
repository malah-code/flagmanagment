import { apiClient } from './apiClient';
import type { Project } from '../types';

export interface CreateProjectPayload {
  name: string;
  description?: string;
}

export interface UpdateProjectPayload {
  name?: string;
  description?: string;
}

export const projectService = {
  async getAll(): Promise<Project[]> {
    return apiClient.get<Project[]>('/projects');
  },

  async getById(id: string): Promise<Project> {
    return apiClient.get<Project>(`/projects/${id}`);
  },

  async create(data: CreateProjectPayload): Promise<Project> {
    return apiClient.post<Project>('/projects', data);
  },

  async update(id: string, data: UpdateProjectPayload): Promise<Project> {
    return apiClient.put<Project>(`/projects/${id}`, data);
  },

  async delete(id: string): Promise<void> {
    return apiClient.delete<void>(`/projects/${id}`);
  },
};
