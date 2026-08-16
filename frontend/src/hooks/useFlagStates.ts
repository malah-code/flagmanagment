import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { flagStateService } from '../services/flagStates';
import type { UpdateFlagStatePayload } from '../services/flagStates';

export const FLAG_STATE_KEYS = {
  all: ['flagStates'] as const,
  byEnvironment: (environmentId: string) =>
    [...FLAG_STATE_KEYS.all, 'environment', environmentId] as const,
};

export function useFlagStates(projectId: string, environmentId: string) {
  return useQuery({
    queryKey: FLAG_STATE_KEYS.byEnvironment(environmentId),
    queryFn: () => flagStateService.getByEnvironment(projectId, environmentId),
    enabled: Boolean(environmentId) && Boolean(projectId),
  });
}

export function useUpdateFlagState(projectId: string, environmentId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ flagId, payload }: { flagId: string; payload: UpdateFlagStatePayload }) =>
      flagStateService.update(projectId, environmentId, flagId, payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: FLAG_STATE_KEYS.byEnvironment(environmentId) });
    },
  });
}

export function useInitFlagState(projectId: string, environmentId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ flagId, payload }: { flagId: string; payload: UpdateFlagStatePayload }) =>
      flagStateService.createOrUpdateByEnvAndFlag(projectId, environmentId, flagId, payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: FLAG_STATE_KEYS.byEnvironment(environmentId) });
    },
  });
}
