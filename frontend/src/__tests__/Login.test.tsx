import { render, screen } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { BrowserRouter } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { Login } from '../pages/Login';

vi.mock('../hooks/useConfig', () => ({
  useSSOProviders: () => ({
    data: {
      oidc_enabled: true,
      saml_enabled: true,
    },
    isLoading: false,
  }),
}));

describe('Login', () => {
  it('renders login form and enabled SSO buttons', () => {
    const queryClient = new QueryClient();
    render(
      <QueryClientProvider client={queryClient}>
        <BrowserRouter>
          <Login />
        </BrowserRouter>
      </QueryClientProvider>,
    );

    expect(screen.getByText('FlagManagment')).toBeInTheDocument();
    expect(screen.getByText('Sign in to your account')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Sign In' })).toBeInTheDocument();

    // Verify dynamic SSO buttons
    expect(screen.getByText('Log in with SSO (OIDC)')).toBeInTheDocument();
    expect(screen.getByText('Log in with SAML 2.0')).toBeInTheDocument();
  });
});
