import { apiClient } from './apiClient';

export interface LoginResponse {
  token: string;
  user: {
    id: string;
    email: string;
  };
}

export const authService = {
  async login(email: string, password: string): Promise<LoginResponse> {
    const data = await apiClient.post<LoginResponse>('/api/v1/auth/login', {
      email,
      password,
    });
    localStorage.setItem('auth_token', data.token);
    return data;
  },

  logout() {
    localStorage.removeItem('auth_token');
  },

  getToken(): string | null {
    return localStorage.getItem('auth_token');
  },

  isAuthenticated(): boolean {
    return !!this.getToken();
  },
};
