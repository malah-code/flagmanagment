import { useQuery } from '@tanstack/react-query';
import { auditService } from '../services/audit';

export const AUDIT_KEYS = {
  byProject: (projectId: string) => ['auditLogs', 'project', projectId] as const,
};

export function useAuditLogs(projectId: string, pageSize: number = 50, pageToken?: string) {
  return useQuery({
    queryKey: [...AUDIT_KEYS.byProject(projectId), pageSize, pageToken],
    queryFn: () => auditService.getByProject(projectId, pageSize, pageToken),
    enabled: !!projectId,
  });
}
