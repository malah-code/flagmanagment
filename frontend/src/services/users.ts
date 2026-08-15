import { apiClient } from './apiClient';

export interface UserResponse {
  id: string;
  email: string;
  auth_provider: string;
  external_id?: string;
  created_at: string;
  updated_at: string;
  roles?: string[];
  projects?: string[];
}

export interface GetUsersResponse {
  users: UserResponse[];
  total: number;
}

export interface InviteUserPayload {
  email: string;
  role: string;
  project_ids: string[];
}

export interface UpdateAccessPayload {
  role: string;
  project_ids: string[];
}

export const userService = {
  async getAll(limit: number = 50, offset: number = 0): Promise<GetUsersResponse> {
    return apiClient.get<GetUsersResponse>(`/users?limit=${limit}&offset=${offset}`);
  },

  async invite(data: InviteUserPayload): Promise<void> {
    return apiClient.post('/users/invite', data);
  },

  async updateAccess(id: string, data: UpdateAccessPayload): Promise<void> {
    return apiClient.put(`/users/${id}/access`, data);
  },
};
