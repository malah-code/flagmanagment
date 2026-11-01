import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import { ChangeRequestDiff } from '../components/ChangeRequestDiff';

describe('ChangeRequestDiff', () => {
  it('renders current and proposed configuration diff headers', () => {
    const current = { enabled: false, rollout: 0 };
    const proposed = { enabled: true, rollout: 50 };

    render(
      <ChangeRequestDiff currentState={current} proposedChanges={proposed} />,
    );

    expect(screen.getByText('Current Configuration')).toBeInTheDocument();
    expect(screen.getByText('Proposed Configuration')).toBeInTheDocument();
  });
});

