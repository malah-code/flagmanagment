import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { flagStateService } from '../services/flagStates';
import type { UpdateFlagStatePayload } from '../services/flagStates';

export const FLAG_STATE_KEYS = {
  all: ['flagStates'] as const,
  byEnvironment: (environmentId: string) => [...FLAG_STATE_KEYS.all, 'environment', environmentId] as const,
};

export function useFlagStates(environmentId: string) {
  return useQuery({
    queryKey: FLAG_STATE_KEYS.byEnvironment(environmentId),
    queryFn: () => flagStateService.getByEnvironment(environmentId),
    enabled: Boolean(environmentId),
  });
}

export function useUpdateFlagState(environmentId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ flagId, payload }: { flagId: string; payload: UpdateFlagStatePayload }) =>
      flagStateService.update(environmentId, flagId, payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: FLAG_STATE_KEYS.byEnvironment(environmentId) });
    },
  });
}

export function useInitFlagState(environmentId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ flagId, payload }: { flagId: string; payload: UpdateFlagStatePayload }) =>
      flagStateService.createOrUpdateByEnvAndFlag(environmentId, flagId, payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: FLAG_STATE_KEYS.byEnvironment(environmentId) });
    },
  });
}
