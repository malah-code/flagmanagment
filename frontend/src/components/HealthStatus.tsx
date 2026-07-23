import React, { useEffect, useState } from 'react';
import { fetchHealth, type HealthResponse } from '../services/api';

const HealthStatus: React.FC = () => {
  const [health, setHealth] = useState<HealthResponse | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState<boolean>(true);

  const checkHealth = async () => {
    try {
      setLoading(true);
      const data = await fetchHealth();
      setHealth(data);
      setError(null);
    } catch (err: unknown) {
      if (err instanceof Error) {
        setError(err.message || 'Failed to fetch health status');
      } else {
        setError('Failed to fetch health status');
      }
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    // eslint-disable-next-line react-hooks/set-state-in-effect
    checkHealth();
    const interval = setInterval(checkHealth, 10000); // Check every 10s
    return () => clearInterval(interval);
  }, []);

  if (loading && !health) {
    return <div>Loading system status...</div>;
  }

  if (error) {
    return (
      <div
        style={{
          color: '#dc3545',
          padding: '1rem',
          border: '1px solid #dc3545',
          borderRadius: '4px',
        }}
      >
        <strong>Error connecting to backend:</strong> {error}
      </div>
    );
  }

  if (!health) return null;

  const isHealthy = health.status === 'healthy';

  return (
    <div>
      <div
        style={{
          display: 'flex',
          alignItems: 'center',
          gap: '0.5rem',
          marginBottom: '1rem',
        }}
      >
        <div
          style={{
            width: '12px',
            height: '12px',
            borderRadius: '50%',
            backgroundColor: isHealthy ? '#28a745' : '#dc3545',
          }}
        />
        <strong style={{ fontSize: '1.1rem' }}>System is {health.status.toUpperCase()}</strong>
      </div>

      <div
        style={{
          display: 'grid',
          gap: '1rem',
          gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))',
        }}
      >
        <div style={cardStyle}>
          <h4>Uptime</h4>
          <p>{health.uptime_seconds}s</p>
        </div>

        <div style={cardStyle}>
          <h4>Version</h4>
          <p>{health.version}</p>
        </div>

        {Object.entries(health.checks || {}).map(([service, check]) => (
          <div
            key={service}
            style={{
              ...cardStyle,
              borderLeft: `4px solid ${check.status === 'healthy' ? '#28a745' : '#dc3545'}`,
            }}
          >
            <h4 style={{ textTransform: 'capitalize' }}>{service}</h4>
            <p>Status: {check.status}</p>
            {check.latency_ms !== undefined && <p>Latency: {check.latency_ms}ms</p>}
            {check.error && <p style={{ color: '#dc3545', fontSize: '0.9em' }}>{check.error}</p>}
          </div>
        ))}
      </div>
    </div>
  );
};

const cardStyle: React.CSSProperties = {
  background: 'white',
  padding: '1rem',
  borderRadius: '4px',
  boxShadow: '0 1px 3px rgba(0,0,0,0.1)',
};

export default HealthStatus;
