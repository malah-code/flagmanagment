import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { flagService } from '../services/flags';
import type { CreateFlagPayload } from '../services/flags';

export const FLAG_KEYS = {
  all: ['flags'] as const,
  byProject: (projectId: string) => [...FLAG_KEYS.all, 'project', projectId] as const,
};

export function useFlags(projectId: string) {
  return useQuery({
    queryKey: FLAG_KEYS.byProject(projectId),
    queryFn: () => flagService.getByProject(projectId),
    enabled: Boolean(projectId),
  });
}

export function useCreateFlag() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (payload: CreateFlagPayload) => flagService.create(payload),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: FLAG_KEYS.byProject(variables.projectId) });
    },
  });
}

export function useDeleteFlag(projectId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => flagService.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: FLAG_KEYS.byProject(projectId) });
    },
  });
}

export function useUpdateFlag(projectId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ flagId, payload }: { flagId: string; payload: import('../services/flags').UpdateFlagPayload }) => 
      flagService.update(projectId, flagId, payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: FLAG_KEYS.byProject(projectId) });
    },
  });
}
