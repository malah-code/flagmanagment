import { render, screen } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { BrowserRouter } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { ProjectsList } from '../pages/ProjectsList';

vi.mock('../hooks/useProjects', () => ({
  useProjects: () => ({
    data: [
      { id: '1', name: 'Test Project 1', description: 'Description 1', createdAt: '2026-01-01' },
    ],
    isLoading: false,
    isError: false,
  }),
  useDeleteProject: () => ({
    mutateAsync: vi.fn(),
  }),
  useCreateProject: () => ({
    mutateAsync: vi.fn(),
    isPending: false,
  }),
}));

describe('ProjectsList', () => {
  it('renders projects correctly', () => {
    const queryClient = new QueryClient();
    render(
      <QueryClientProvider client={queryClient}>
        <BrowserRouter>
          <ProjectsList />
        </BrowserRouter>
      </QueryClientProvider>,
    );

    expect(screen.getByText('Test Project 1')).toBeInTheDocument();
    expect(screen.getByText('Description 1')).toBeInTheDocument();
  });
});
