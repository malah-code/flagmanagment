import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';

interface WebhookIntegration {
  id: string;
  project_id: string;
  name: string;
  integration_type: string;
  url: string;
  events: string[];
  is_active: boolean;
  created_at: string;
}

export const useWebhooks = (projectId: string) => {
  return useQuery<WebhookIntegration[]>({
    queryKey: ['webhooks', projectId],
    queryFn: async () => {
      const res = await fetch(`/api/v1/projects/${projectId}/webhooks`);
      if (!res.ok) throw new Error('Failed to fetch webhooks');
      return res.json();
    },
    enabled: !!projectId,
  });
};

export const useCreateWebhook = (projectId: string) => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (data: {
      name: string;
      integration_type: string;
      url: string;
      secret_key?: string;
      events: string[];
      is_active: boolean;
    }) => {
      const res = await fetch(`/api/v1/projects/${projectId}/webhooks`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data),
      });
      if (!res.ok) throw new Error('Failed to create webhook');
      return res.json();
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['webhooks', projectId] });
    },
  });
};
