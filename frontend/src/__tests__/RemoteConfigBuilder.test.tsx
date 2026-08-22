import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { RemoteConfigBuilder } from '../components/flagStates/RemoteConfigBuilder';

const mockMutateAsync = vi.fn();

vi.mock('../../hooks/useFlagStates', () => ({
  useUpdateFlagState: () => ({
    mutateAsync: mockMutateAsync,
    isPending: false,
  }),
}));

describe('RemoteConfigBuilder', () => {
  it('renders modal and initial fields in visual mode', () => {
    const queryClient = new QueryClient();
    render(
      <QueryClientProvider client={queryClient}>
        <RemoteConfigBuilder
          isOpen={true}
          onClose={vi.fn()}
          envId="env-1"
          projectId="proj-1"
          flagId="flag-1"
          flagKey="beta_feature"
          initialConfig={{ primaryColor: 'indigo', timeoutMs: 5000, enabled: true }}
        />
      </QueryClientProvider>,
    );

    expect(screen.getByText('Remote Config Payload')).toBeInTheDocument();
    expect(screen.getByText('beta_feature')).toBeInTheDocument();
    expect(screen.getByDisplayValue('primaryColor')).toBeInTheDocument();
    expect(screen.getByDisplayValue('indigo')).toBeInTheDocument();
    expect(screen.getByDisplayValue('timeoutMs')).toBeInTheDocument();
    expect(screen.getByDisplayValue('5000')).toBeInTheDocument();
  });

  it('allows switching to raw JSON mode', () => {
    const queryClient = new QueryClient();
    render(
      <QueryClientProvider client={queryClient}>
        <RemoteConfigBuilder
          isOpen={true}
          onClose={vi.fn()}
          envId="env-1"
          projectId="proj-1"
          flagId="flag-1"
          flagKey="beta_feature"
          initialConfig={{ maxRetries: 3 }}
        />
      </QueryClientProvider>,
    );

    const jsonModeButton = screen.getByText('Raw JSON');
    fireEvent.click(jsonModeButton);

    const textarea = screen.getByRole('textbox');
    expect(textarea).toBeInTheDocument();
    expect(textarea).toHaveValue(JSON.stringify({ maxRetries: 3 }, null, 2));
  });
});
