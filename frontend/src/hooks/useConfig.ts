import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { configService } from '../services/config';
import type { SMTPConfig } from '../services/config';

export const CONFIG_KEYS = {
  smtp: ['config', 'smtp'] as const,
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
