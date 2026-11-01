import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import { FlagTelemetryPanel } from '../components/flags/FlagTelemetryPanel';

describe('FlagTelemetryPanel Component', () => {
  const sampleVariations = [
    { id: 'v1', name: 'Treatment A', value: true },
    { id: 'v2', name: 'Control B', value: false },
  ];

  it('renders telemetry metrics and default 24h stats', () => {
    render(
      <FlagTelemetryPanel
        flagKey="checkout-v2"
        flagType="boolean"
        currentEnvName="Production"
        variations={sampleVariations}
      />,
    );

    expect(screen.getByText('Total Evaluations')).toBeInTheDocument();
    expect(screen.getByText('1,158,400')).toBeInTheDocument();
    expect(screen.getByText('Evaluation Rate')).toBeInTheDocument();
    expect(screen.getByText('14.8/s')).toBeInTheDocument();
    expect(screen.getByText('P99 Eval Latency')).toBeInTheDocument();
    expect(screen.getByText('Cache Hit Ratio')).toBeInTheDocument();
    expect(screen.getByText('99.98%')).toBeInTheDocument();
    expect(screen.getByText(/Production/i)).toBeInTheDocument();
  });

  it('renders variation traffic distribution with custom variations', () => {
    render(
      <FlagTelemetryPanel
        flagKey="checkout-v2"
        flagType="boolean"
        currentEnvName="Production"
        variations={sampleVariations}
      />,
    );

    expect(screen.getByText('Variation Traffic Distribution')).toBeInTheDocument();
    expect(screen.getByText('Treatment A')).toBeInTheDocument();
    expect(screen.getByText('Control B')).toBeInTheDocument();
  });

  it('updates metrics when switching time ranges', () => {
    render(
      <FlagTelemetryPanel flagKey="checkout-v2" flagType="boolean" currentEnvName="Production" />,
    );

    // Switch to 1H
    const btn1h = screen.getByRole('button', { name: '1H' });
    fireEvent.click(btn1h);
    expect(screen.getByText('48,290')).toBeInTheDocument();
    expect(screen.getByText('13.4/s')).toBeInTheDocument();

    // Switch to 7D
    const btn7d = screen.getByRole('button', { name: '7D' });
    fireEvent.click(btn7d);
    expect(screen.getByText('8,420,150')).toBeInTheDocument();
  });

  it('toggles live streaming between paused and active states', () => {
    render(
      <FlagTelemetryPanel flagKey="checkout-v2" flagType="boolean" currentEnvName="Production" />,
    );

    expect(screen.getByText('Live Stream Active')).toBeInTheDocument();
    const pauseBtn = screen.getByRole('button', { name: /pause/i });
    fireEvent.click(pauseBtn);

    expect(screen.getByText('Live Stream Paused')).toBeInTheDocument();
    const resumeBtn = screen.getByRole('button', { name: /resume/i });
    fireEvent.click(resumeBtn);
    expect(screen.getByText('Live Stream Active')).toBeInTheDocument();
  });

  it('renders recent evaluation audit stream table with latency and context keys', () => {
    render(
      <FlagTelemetryPanel flagKey="search-algorithm" flagType="string" currentEnvName="Staging" />,
    );

    expect(screen.getByText('Recent Evaluation Stream')).toBeInTheDocument();
    expect(screen.getByText('usr_89f01ab2')).toBeInTheDocument();
    expect(screen.getByText('usr_premium_9921')).toBeInTheDocument();
    expect(screen.getByText('TARGETING_MATCH')).toBeInTheDocument();
  });
});
