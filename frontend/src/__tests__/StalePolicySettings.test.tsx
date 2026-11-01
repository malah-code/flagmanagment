import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { StalePolicySettings } from '../components/settings/StalePolicySettings';

vi.mock('../services/api', () => ({
  getStalePolicy: vi.fn().mockResolvedValue({ stale_after_days: 30 }),
  setStalePolicy: vi.fn().mockResolvedValue({ stale_after_days: 45 }),
  runStaleScan: vi.fn().mockResolvedValue({
    flags_scanned: 10,
    flags_marked_stale: 3,
    change_requests_created: 2,
    notifications_sent: 2,
  }),
}));

describe('StalePolicySettings Component', () => {
  it('renders policy threshold input and governance bot information', async () => {
    render(<StalePolicySettings projectId="test-proj-123" />);

    expect(await screen.findByText('Stale Flag Detection Policy')).toBeInTheDocument();
    expect(screen.getByText('Automated Governance Bot')).toBeInTheDocument();
    expect(screen.getByText('Active')).toBeInTheDocument();
    expect(screen.getByText('Auto Deprecation Proposals')).toBeInTheDocument();
    expect(screen.getByText('Slack & Webhook Alerts')).toBeInTheDocument();
  });

  it('triggers on-demand governance scan and displays results', async () => {
    render(<StalePolicySettings projectId="test-proj-123" />);

    const runBtn = await screen.findByRole('button', { name: /Run Governance Bot Scan/i });
    fireEvent.click(runBtn);

    await waitFor(() => {
      expect(screen.getByText('Latest Governance Scan Result')).toBeInTheDocument();
      expect(screen.getByText('10')).toBeInTheDocument(); // scanned
      expect(screen.getByText('3')).toBeInTheDocument(); // marked stale
      expect(screen.getByText('CRS DRAFTED')).toBeInTheDocument();
      expect(screen.getAllByText('2').length).toBeGreaterThanOrEqual(1);
    });
  });
});
