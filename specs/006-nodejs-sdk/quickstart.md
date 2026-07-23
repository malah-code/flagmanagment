# Quickstart: Node.js SDK

## Prerequisites
- Node.js 20+
- A running FlagManagment backend server (`make up` from repo root)
- An active Environment Token from the FlagManagment database.

## Setup Project

```bash
mkdir test-sdk && cd test-sdk
npm init -y
npm install typescript @types/node ts-node --save-dev
npm install @openfeature/server-sdk
# (Once published/linked) npm install @flagmanagment/node-sdk
```

Create a `test.ts` file:

```typescript
import { FlagManagmentClient, FlagManagmentProvider } from '@flagmanagment/node-sdk';
import { OpenFeature } from '@openfeature/server-sdk';

async function run() {
  // 1. Initialize client
  const client = new FlagManagmentClient({
    environmentToken: process.env.FLAGMANAGMENT_TOKEN || 'test-token',
    endpoint: 'http://localhost:8080'
  });
  
  await client.init();
  console.log("SDK connected and initialized!");

  // 2. Set OpenFeature provider
  OpenFeature.setProvider(new FlagManagmentProvider(client));
  const featureClient = OpenFeature.getClient();

  // 3. Evaluate a flag
  const result = await featureClient.getBooleanValue('test-feature', false, { targetingKey: 'user-1' });
  console.log("Evaluation Result:", result);

  // Keep alive to test streaming
  console.log("Waiting for realtime updates...");
  setInterval(() => {}, 1000);
}

run().catch(console.error);
```

## Running the validation

1. Start the Go backend and ensure it has a valid environment and a `test-feature` flag.
2. Get the environment token from the PostgreSQL database (`environment_tokens` table).
3. Run the script:
   ```bash
   export FLAGMANAGMENT_TOKEN="your-token"
   npx ts-node test.ts
   ```
4. Verify the output correctly evaluates to the flag state. Change the flag state in the database/UI and verify that the SDK stream handles it (if you evaluate again inside the `setInterval`).
