import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { OpenFeaturePlayground } from '../pages/OpenFeaturePlayground';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

// Mock useFlags and useEnvironments
vi.mock('../hooks/useFlags', () => ({
  useFlags: () => ({
    data: [
      { id: 'flag-1', key: 'checkout-v2', name: 'Checkout V2', type: 'boolean' },
      { id: 'flag-2', key: 'search-model', name: 'Search Model', type: 'string' },
    ],
    isLoading: false,
  }),
}));

vi.mock('../hooks/useEnvironments', () => ({
  useEnvironments: () => ({
    data: [
      { id: 'env-prod', name: 'Production', isProtected: true },
      { id: 'env-stage', name: 'Staging', isProtected: false },
    ],
    isLoading: false,
  }),
}));

const queryClient = new QueryClient();

describe('OpenFeaturePlayground Component', () => {
  const renderComponent = () =>
    render(
      <QueryClientProvider client={queryClient}>
        <OpenFeaturePlayground projectId="test-proj-123" />
      </QueryClientProvider>,
    );

  it('renders OpenFeature playground header and conformance checklist', () => {
    renderComponent();
    expect(screen.getByText(/OpenFeature Playground & Conformance Hub/i)).toBeInTheDocument();
    expect(screen.getByText(/CNCF OpenFeature Standardized/i)).toBeInTheDocument();
    expect(screen.getByText(/Deterministic MurmurHash3 Bucketing/i)).toBeInTheDocument();
  });

  it('allows selecting evaluation call types and running evaluation', () => {
    renderComponent();

    // Verify method buttons
    expect(screen.getByRole('button', { name: /getBooleanValue/i })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /getStringValue/i })).toBeInTheDocument();

    // Click Run Evaluation
    const runBtn = screen.getByRole('button', { name: /Run OpenFeature Evaluation/i });
    fireEvent.click(runBtn);

    // Verify ResolutionDetails card appears
    expect(screen.getByText('ResolutionDetails')).toBeInTheDocument();
    expect(screen.getByText(/TARGETING_MATCH|PERCENTAGE_ROLLOUT/)).toBeInTheDocument();
  });

  it('switches SDK languages and renders matching code snippets', () => {
    renderComponent();

    // Default is React
    expect(screen.getByText(/FlagManagmentWebProvider/i)).toBeInTheDocument();

    // Switch to Node.js
    const nodeTab = screen.getByRole('button', { name: 'Node.js' });
    fireEvent.click(nodeTab);
    expect(screen.getByText(/@openfeature\/server-sdk/i)).toBeInTheDocument();

    // Switch to Python
    const pythonTab = screen.getByRole('button', { name: 'Python' });
    fireEvent.click(pythonTab);
    expect(screen.getByText(/openfeature-sdk/i)).toBeInTheDocument();

    // Switch to Go
    const goTab = screen.getByRole('button', { name: 'Go' });
    fireEvent.click(goTab);
    expect(screen.getByText(/github\.com\/open-feature\/go-sdk/i)).toBeInTheDocument();
  });
});
