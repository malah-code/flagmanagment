import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { configService } from '../services/config';
import type { SMTPConfig, SSOConfig } from '../services/config';

export const CONFIG_KEYS = {
  smtp: ['config', 'smtp'] as const,
  sso: ['config', 'sso'] as const,
  ssoProviders: ['auth', 'ssoProviders'] as const,
};

export function useSMTPConfig() {
  return useQuery({
    queryKey: CONFIG_KEYS.smtp,
    queryFn: () => configService.getSMTP(),
    retry: false,
  });
}

export function useUpdateSMTPConfig() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: SMTPConfig) => configService.updateSMTP(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: CONFIG_KEYS.smtp });
    },
  });
}

export function useTestSMTP() {
  return useMutation({
    mutationFn: (email: string) => configService.testSMTP(email),
  });
}

export function useSSOConfig() {
  return useQuery({
    queryKey: CONFIG_KEYS.sso,
    queryFn: () => configService.getSSO(),
    retry: false,
  });
}

export function useUpdateSSOConfig() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: SSOConfig) => configService.updateSSO(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: CONFIG_KEYS.sso });
      queryClient.invalidateQueries({ queryKey: CONFIG_KEYS.ssoProviders });
    },
  });
}

export function useSSOProviders() {
  return useQuery({
    queryKey: CONFIG_KEYS.ssoProviders,
    queryFn: () => configService.getSSOProviders(),
    retry: false,
  });
}
