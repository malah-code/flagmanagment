import { apiClient } from './apiClient';

export interface SMTPConfig {
  host: string;
  port: number;
  username?: string;
  password?: string;
}

export const configService = {
  async getSMTP(): Promise<SMTPConfig> {
    return apiClient.get<SMTPConfig>('/config/smtp');
  },
  async updateSMTP(data: SMTPConfig): Promise<SMTPConfig> {
    return apiClient.put<SMTPConfig>('/config/smtp', data);
  },
  async testSMTP(email: string): Promise<void> {
    return apiClient.post('/config/smtp/test', { email });
  }
};
