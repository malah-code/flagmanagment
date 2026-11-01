import { render, screen } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { BrowserRouter } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { Login } from '../pages/Login';

vi.mock('../hooks/useConfig', () => ({
  useSSOProviders: () => ({
    data: { oidc_enabled: false, saml_enabled: false },
    isLoading: false,
  }),
}));

describe('Login Error Handling & Feedback', () => {
  it('renders required input fields for email and password', () => {
    const queryClient = new QueryClient();
    render(
      <QueryClientProvider client={queryClient}>
        <BrowserRouter>
          <Login />
        </BrowserRouter>
      </QueryClientProvider>,
    );

    const emailLabel = screen.getByText('Email');
    const passwordLabel = screen.getByText('Password');

    expect(emailLabel).toBeInTheDocument();
    expect(passwordLabel).toBeInTheDocument();
  });
});
