import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { environmentService } from '../services/environments';
import type { CreateEnvironmentPayload } from '../services/environments';

export const ENVIRONMENT_KEYS = {
  all: ['environments'] as const,
  byProject: (projectId: string) => [...ENVIRONMENT_KEYS.all, 'project', projectId] as const,
};

export function useEnvironments(projectId: string) {
  return useQuery({
    queryKey: ENVIRONMENT_KEYS.byProject(projectId),
    queryFn: () => environmentService.getByProject(projectId),
    enabled: Boolean(projectId),
  });
}

export function useCreateEnvironment() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (payload: CreateEnvironmentPayload) => environmentService.create(payload),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ENVIRONMENT_KEYS.byProject(variables.projectId) });
    },
  });
}

export function useUpdateEnvironment(projectId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (payload: import('../services/environments').UpdateEnvironmentPayload) => environmentService.update(payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ENVIRONMENT_KEYS.byProject(projectId) });
    },
  });
}

export function useDeleteEnvironment(projectId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => environmentService.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ENVIRONMENT_KEYS.byProject(projectId) });
    },
  });
}

export function useCloneEnvironment(projectId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ envId, name }: { envId: string; name: string }) => 
      environmentService.clone(projectId, envId, { name }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ENVIRONMENT_KEYS.byProject(projectId) });
    },
  });
}
