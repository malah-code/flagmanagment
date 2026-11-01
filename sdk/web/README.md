# FlagManagment — Web & Vanilla JavaScript SDK

The official client-side JavaScript / TypeScript SDK for **FlagManagment**, designed for modern web browsers, SPAs, and frameworks (Vue, Svelte, Angular, SolidJS, or vanilla script tags).

---

## Features

- **Zero Runtime Dependencies**: Ultra-lightweight footprint (<4 kB gzipped).
- **Instant Page Loads**: Synchronous flag reads backed by in-memory cache and optional `localStorage` persistence.
- **Real-Time Updates**: Live Server-Sent Events (SSE) delta updates.
- **Event-Driven**: Listen for flag changes with `client.on('change', ...)`.
- **CNCF OpenFeature**: Includes `OpenFeatureWebProvider` for standardized flag resolution.

---

## Installation

```bash
npm install @flagmanagment/web-sdk
```

---

## Quickstart

```typescript
import { initFlagManagment } from '@flagmanagment/web-sdk';

// Initialize the web client
const client = initFlagManagment({
  apiKey: 'flg_client_your_client_key',
  apiUrl: 'https://flags.example.com',
  context: {
    targetingKey: 'user-123',
    tier: 'premium',
  },
});

// Evaluate flags synchronously (0 latency)
const isV2Checkout = client.getBooleanValue('new-checkout-v2', false);
const bannerText = client.getStringValue('promotional-banner', 'Welcome!');

// Subscribe to real-time flag changes
client.on('change', (flagKey, newValue, oldValue) => {
  console.log(`Flag ${flagKey} changed from ${oldValue} to ${newValue}`);
});
```

---

## Building & Testing

```bash
npm run build
npm test
```
