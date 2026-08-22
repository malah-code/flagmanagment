import { apiClient } from './apiClient';

export interface SMTPConfig {
  host: string;
  port: number;
  username?: string;
  password?: string;
}

export interface OIDCConfig {
  enabled: boolean;
  issuer_url: string;
  client_id: string;
  client_secret?: string;
}

export interface SAMLConfig {
  enabled: boolean;
  idp_metadata_url: string;
  idp_entity_id: string;
  idp_sso_url: string;
  idp_cert: string;
  sp_entity_id: string;
}

export interface SSOConfig {
  oidc: OIDCConfig;
  saml: SAMLConfig;
}

export interface SSOProviders {
  oidc_enabled: boolean;
  saml_enabled: boolean;
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
  },
  async getSSO(): Promise<SSOConfig> {
    return apiClient.get<SSOConfig>('/config/sso');
  },
  async updateSSO(data: SSOConfig): Promise<SSOConfig> {
    return apiClient.put<SSOConfig>('/config/sso', data);
  },
  async getSSOProviders(): Promise<SSOProviders> {
    return apiClient.get<SSOProviders>('/auth/sso/providers');
  },
};
