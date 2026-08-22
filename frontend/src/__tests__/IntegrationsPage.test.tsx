import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { IntegrationsPage } from '../pages/IntegrationsPage';

const mockMutateAsync = vi.fn();

vi.mock('../hooks/useWebhooks', () => ({
  useWebhooks: () => ({
    data: [
      {
        id: 'wh-1',
        project_id: 'proj-1',
        name: 'My PostHog',
        integration_type: 'posthog',
        url: 'https://app.posthog.com',
        events: ['flag.updated'],
        is_active: true,
        created_at: '2026-08-22',
      },
    ],
    isLoading: false,
  }),
  useCreateWebhook: () => ({
    mutateAsync: mockMutateAsync,
    isPending: false,
  }),
}));

describe('IntegrationsPage', () => {
  it('renders integration catalog and active integrations', () => {
    const queryClient = new QueryClient();
    render(
      <QueryClientProvider client={queryClient}>
        <IntegrationsPage projectId="proj-1" />
      </QueryClientProvider>,
    );

    expect(screen.getByText('Integrations')).toBeInTheDocument();
    expect(screen.getByText('PostHog')).toBeInTheDocument();
    expect(screen.getByText('Datadog')).toBeInTheDocument();
    expect(screen.getByText('Custom Webhook')).toBeInTheDocument();

    // Check active integration item
    expect(screen.getByText('My PostHog')).toBeInTheDocument();
    expect(screen.getByText('Active')).toBeInTheDocument();
  });

  it('opens configuration modal when clicking connect', () => {
    const queryClient = new QueryClient();
    render(
      <QueryClientProvider client={queryClient}>
        <IntegrationsPage projectId="proj-1" />
      </QueryClientProvider>,
    );

    const connectButtons = screen.getAllByText('Connect');
    fireEvent.click(connectButtons[0]); // Click connect on PostHog

    expect(screen.getByText('Configure PostHog')).toBeInTheDocument();
    expect(screen.getByText('PostHog Host')).toBeInTheDocument();
  });
});
