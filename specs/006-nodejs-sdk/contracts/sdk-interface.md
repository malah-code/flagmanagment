# API Contract: FlagManagment Node.js SDK

## Initialization

```typescript
import { FlagManagmentClient } from '@flagmanagment/node-sdk';

const client = new FlagManagmentClient({
  environmentToken: 'env_...',
  endpoint: 'http://localhost:8080'
});

// Await initialization (fetches snapshot)
await client.init();
```

## Direct Evaluation (Proprietary API)

```typescript
// Evaluate a boolean flag
const isEnabled = client.evaluateBoolean('new-checkout-flow', { identity: 'user-123' }, false);

// Evaluate a multivariate flag
const variation = client.evaluateString('button-color', { identity: 'user-123' }, 'blue');
```

## OpenFeature Provider Interface

```typescript
import { OpenFeature } from '@openfeature/server-sdk';
import { FlagManagmentProvider } from '@flagmanagment/node-sdk';

// Register the provider
OpenFeature.setProvider(new FlagManagmentProvider(client));

// Get standard OpenFeature client
const featureClient = OpenFeature.getClient();

// Evaluate using standard OpenFeature API
const isEnabled = await featureClient.getBooleanValue('new-checkout-flow', false, { targetingKey: 'user-123' });
```
