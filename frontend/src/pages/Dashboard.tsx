import React from 'react';
import HealthStatus from '../components/HealthStatus';

const Dashboard: React.FC = () => {
  return (
    <div style={{ padding: '2rem', fontFamily: 'system-ui, sans-serif' }}>
      <header style={{ marginBottom: '2rem' }}>
        <h1>FlagManagment Dashboard</h1>
        <p style={{ color: '#666' }}>System Status Overview</p>
      </header>

      <main>
        <section
          style={{
            background: '#f8f9fa',
            padding: '1.5rem',
            borderRadius: '8px',
            border: '1px solid #dee2e6',
          }}
        >
          <h2>Infrastructure Health</h2>
          <HealthStatus />
        </section>
      </main>
    </div>
  );
};

export default Dashboard;
