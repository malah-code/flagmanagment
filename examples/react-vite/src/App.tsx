import { useState } from 'react';
import { FlagProvider, useFlag, useFlags } from '@flagmanagment/react-sdk';

const USERS = [
  { id: 'usr-standard-1', name: 'Alice (Standard)', tier: 'free', country: 'US' },
  { id: 'usr-pro-2', name: 'Bob (Pro Tier)', tier: 'pro', country: 'CA' },
  { id: 'usr-enterprise-3', name: 'Charlie (Enterprise VIP)', tier: 'enterprise', country: 'DE' },
];

export default function App() {
  const [currentUser, setCurrentUser] = useState(USERS[0]);

  return (
    <FlagProvider
      apiKey="flg_client_sample_react_key"
      apiUrl="http://localhost:8080"
      context={{
        targetingKey: currentUser.id,
        tier: currentUser.tier,
        country: currentUser.country,
      }}
    >
      <div style={{ maxWidth: 800, margin: '40px auto', padding: '0 20px' }}>
        <header style={{ borderBottom: '1px solid #e2e8f0', paddingBottom: 20, marginBottom: 30 }}>
          <h1 style={{ fontSize: 24, fontWeight: 700, margin: 0, color: '#4f46e5' }}>
            FlagManagment — React SDK Sample
          </h1>
          <p style={{ color: '#64748b', fontSize: 14, marginTop: 6 }}>
            Demonstrates context-aware targeting, dynamic hooks, and live re-evaluation.
          </p>
        </header>

        {/* User Switcher to demonstrate context changing */}
        <section style={{ background: '#fff', border: '1px solid #e2e8f0', borderRadius: 12, padding: 20, marginBottom: 24 }}>
          <label style={{ fontSize: 12, fontWeight: 700, textTransform: 'uppercase', color: '#64748b', display: 'block', marginBottom: 8 }}>
            Simulate Evaluation Context (User Identity & Tier)
          </label>
          <div style={{ display: 'flex', gap: 12 }}>
            {USERS.map((user) => (
              <button
                key={user.id}
                onClick={() => setCurrentUser(user)}
                style={{
                  padding: '8px 16px',
                  borderRadius: 8,
                  fontSize: 13,
                  fontWeight: 600,
                  cursor: 'pointer',
                  border: currentUser.id === user.id ? '2px solid #4f46e5' : '1px solid #cbd5e1',
                  background: currentUser.id === user.id ? '#eef2ff' : '#fff',
                  color: currentUser.id === user.id ? '#4338ca' : '#334155',
                }}
              >
                {user.name}
              </button>
            ))}
          </div>
        </section>

        {/* Feature Component Demonstrating Hooks */}
        <FeatureShowcase currentUser={currentUser} />
      </div>
    </FlagProvider>
  );
}

function FeatureShowcase({ currentUser }: { currentUser: typeof USERS[0] }) {
  // Evaluates boolean feature flag with local default fallback
  const isOneClickCheckout = useFlag<boolean>('new-checkout-v2', false);

  // Evaluates multivariate string flag
  const bannerMessage = useFlag<string>('promotional-banner', 'Welcome to our platform!');

  // Evaluates number limit flag
  const downloadLimit = useFlag<number>('max-downloads', 10);

  // Get full flag state dictionary
  const allFlags = useFlags();

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: 20 }}>
      {/* Dynamic Banner */}
      <div style={{ background: '#312e81', color: '#fff', padding: 16, borderRadius: 12, fontSize: 14, fontWeight: 500 }}>
        📢 <strong>Live Promotion:</strong> {bannerMessage}
      </div>

      {/* Feature Cards Grid */}
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(280px, 1fr))', gap: 16 }}>
        {/* Checkout Card */}
        <div style={{ background: '#fff', border: '1px solid #e2e8f0', borderRadius: 12, padding: 20 }}>
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 12 }}>
            <h3 style={{ margin: 0, fontSize: 16, fontWeight: 600 }}>Checkout Experience</h3>
            <span style={{
              fontSize: 11,
              fontWeight: 700,
              padding: '2px 8px',
              borderRadius: 20,
              background: isOneClickCheckout ? '#dcfce7' : '#f1f5f9',
              color: isOneClickCheckout ? '#15803d' : '#475569',
            }}>
              {isOneClickCheckout ? 'V2 Enabled' : 'V1 Standard'}
            </span>
          </div>
          <p style={{ fontSize: 13, color: '#64748b', lineHeight: 1.5, margin: '0 0 16px 0' }}>
            Flag <code>new-checkout-v2</code> evaluated for user <strong>{currentUser.id}</strong> (Tier: {currentUser.tier}).
          </p>
          <button style={{
            width: '100%',
            padding: '10px 0',
            background: isOneClickCheckout ? '#10b981' : '#4f46e5',
            color: '#fff',
            border: 'none',
            borderRadius: 8,
            fontWeight: 600,
            cursor: 'pointer',
          }}>
            {isOneClickCheckout ? '⚡ 1-Click Instant Checkout' : 'Proceed to Regular Checkout'}
          </button>
        </div>

        {/* Quota Card */}
        <div style={{ background: '#fff', border: '1px solid #e2e8f0', borderRadius: 12, padding: 20 }}>
          <h3 style={{ margin: '0 0 12px 0', fontSize: 16, fontWeight: 600 }}>Tier Quota Limits</h3>
          <p style={{ fontSize: 13, color: '#64748b', lineHeight: 1.5, margin: '0 0 16px 0' }}>
            Flag <code>max-downloads</code> integer limit.
          </p>
          <div style={{ fontSize: 28, fontWeight: 700, color: '#0f172a' }}>
            {downloadLimit} <span style={{ fontSize: 14, fontWeight: 400, color: '#64748b' }}>files/day</span>
          </div>
        </div>
      </div>

      {/* Raw Snapshot Inspection */}
      <div style={{ background: '#fff', border: '1px solid #e2e8f0', borderRadius: 12, padding: 20 }}>
        <h4 style={{ margin: '0 0 8px 0', fontSize: 13, fontWeight: 700, color: '#475569' }}>
          Real-time Cached Flags Snapshot ({Object.keys(allFlags).length} flags active)
        </h4>
        <pre style={{ background: '#f8fafc', padding: 12, borderRadius: 8, fontSize: 12, overflowX: 'auto', margin: 0, color: '#334155' }}>
          {JSON.stringify(allFlags, null, 2)}
        </pre>
      </div>
    </div>
  );
}

