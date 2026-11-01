import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import { StaleBadge } from '../components/flags/StaleBadge';

describe('Flag Targeting & Lifecycle Badges', () => {
  it('renders stale badge for flags flagged as STALE', () => {
    render(<StaleBadge state="STALE" />);
    expect(screen.getByText('Stale')).toBeInTheDocument();
  });

  it('renders deprecated badge for flags flagged as DEPRECATED', () => {
    render(<StaleBadge state="DEPRECATED" />);
    expect(screen.getByText('Deprecated')).toBeInTheDocument();
  });
});

