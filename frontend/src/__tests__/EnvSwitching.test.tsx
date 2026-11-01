import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import { ClientSideKeyCard } from '../components/environments/ClientSideKeyCard';

describe('Environment SDK Key Display & Switching', () => {
  it('displays the client-side environment key correctly', () => {
    render(<ClientSideKeyCard envName="Production" apiKey="env_prod_secret_token_123" />);

    expect(screen.getByText(/Client-side Environment Key/i)).toBeInTheDocument();
    expect(screen.getByText('env_prod_secret_token_123')).toBeInTheDocument();
  });
});
