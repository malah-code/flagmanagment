# @flagmanagment/node-sdk

Official Node.js / TypeScript SDK for FlagManagment feature flag and remote configuration platform.

## Features
- In-memory local flag evaluation (<1ms latency)
- CNCF OpenFeature Provider compliant
- SHA-256 user PII hashing per Constitution VII
- Real-time delta updates via background synchronization

## Installation

```bash
npm install @flagmanagment/node-sdk @openfeature/server-sdk
```

## Quick Start

```typescript
import { FlagManagmentClient, FlagManagmentProvider } from '@flagmanagment/node-sdk';
import { OpenFeature } from '@openfeature/server-sdk';

const client = new FlagManagmentClient({
  environmentToken: 'your-env-token',
  endpoint: 'http://localhost:8080'
});

await client.init();

// Use OpenFeature API
OpenFeature.setProvider(new FlagManagmentProvider(client));
const featureClient = OpenFeature.getClient();

const isEnabled = await featureClient.getBooleanValue('my-feature', false, { targetingKey: 'user-123' });
console.log('Feature state:', isEnabled);
```
