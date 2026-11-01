const express = require('express');
const { OpenFeature } = require('@openfeature/server-sdk');
const { FlagManagmentProvider } = require('@flagmanagment/node-sdk');

const app = express();
app.use(express.json());

const PORT = process.env.PORT || 4000;
const SERVER_KEY = process.env.FLAGMANAGMENT_SERVER_KEY || 'sample-server-sdk-key';
const BASE_URL = process.env.FLAGMANAGMENT_API_URL || 'http://localhost:8080';

// Initialize the FlagManagment OpenFeature Provider
const provider = new FlagManagmentProvider({
  serverKey: SERVER_KEY,
  baseUrl: BASE_URL,
  refreshInterval: 10000, // 10s polling fallback if SSE inactive
});

OpenFeature.setProvider(provider);
const client = OpenFeature.getClient();

// Health Check
app.get('/health', (req, res) => {
  res.json({ status: 'ok', service: 'node-express-sample' });
});

// Checkout Flow Route - Demonstrates Boolean Flag with Contextual Targeting
app.get('/api/checkout', async (req, res) => {
  const userId = req.query.userId || 'user-123';
  const tier = req.query.tier || 'standard';
  const country = req.query.country || 'US';

  // Build OpenFeature Evaluation Context
  const context = {
    targetingKey: userId,
    tier,
    country,
  };

  try {
    // Evaluate boolean flag
    const isV2Enabled = await client.getBooleanValue('new-checkout-v2', false, context);

    if (isV2Enabled) {
      return res.json({
        flow: 'v2-one-click-checkout',
        features: ['apple-pay', 'google-pay', 'instant-address-autofill'],
        context,
      });
    }

    return res.json({
      flow: 'v1-standard-checkout',
      features: ['standard-credit-card'],
      context,
    });
  } catch (error) {
    return res.status(500).json({ error: error.message });
  }
});

// Multivariate Banner Route - Demonstrates String / Multivariate Flag
app.get('/api/banner', async (req, res) => {
  const userId = req.query.userId || 'anon-visitor';
  const context = { targetingKey: userId };

  try {
    const bannerVariant = await client.getStringValue(
      'promotional-banner',
      'control',
      context
    );

    res.json({
      variant: bannerVariant,
      message:
        bannerVariant === 'discount-20'
          ? 'Special Flash Deal: 20% off all orders!'
          : bannerVariant === 'free-shipping'
            ? 'Free Priority Shipping on orders over $50!'
            : 'Welcome to our store!',
      context,
    });
  } catch (error) {
    res.status(500).json({ error: error.message });
  }
});

// Generic Flag Evaluation Endpoint - Demonstrates full resolution details
app.post('/api/evaluate', async (req, res) => {
  const { flagKey, defaultValue = false, context = {} } = req.body;

  if (!flagKey) {
    return res.status(400).json({ error: 'flagKey is required' });
  }

  try {
    const details = await client.getBooleanDetails(flagKey, defaultValue, context);
    res.json({
      flagKey,
      value: details.value,
      reason: details.reason,
      variant: details.variant,
    });
  } catch (error) {
    res.status(500).json({ error: error.message });
  }
});

if (require.main === module) {
  app.listen(PORT, () => {
    console.log(`🚀 FlagManagment Node.js Express Sample listening on http://localhost:${PORT}`);
    console.log(`👉 Try: http://localhost:${PORT}/api/checkout?userId=vip-42&tier=premium`);
    console.log(`👉 Try: http://localhost:${PORT}/api/banner?userId=user-99`);
  });
}

module.exports = app;

