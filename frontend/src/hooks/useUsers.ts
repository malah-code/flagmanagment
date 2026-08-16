import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { userService } from '../services/users';
import type { InviteUserPayload, UpdateAccessPayload } from '../services/users';

export const USER_KEYS = {
  all: ['users'] as const,
};

export function useUsers(limit: number = 50, offset: number = 0) {
  return useQuery({
    queryKey: [...USER_KEYS.all, limit, offset],
    queryFn: () => userService.getAll(limit, offset),
  });
}

export function useInviteUser() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (payload: InviteUserPayload) => userService.invite(payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: USER_KEYS.all });
    },
  });
}

export function useUpdateUserAccess() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, payload }: { id: string; payload: UpdateAccessPayload }) =>
      userService.updateAccess(id, payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: USER_KEYS.all });
    },
  });
}
