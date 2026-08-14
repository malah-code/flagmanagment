const http = require('http');

const API_BASE = process.env.API_BASE || 'http://localhost:3000/api/v1';

async function request(path, method = 'GET', body = null, token = '') {
  const url = new URL(API_BASE + path);
  const options = {
    method,
    headers: {
      'Content-Type': 'application/json',
      'Accept': 'application/json',
      ...(token ? { 'Authorization': `Bearer ${token}` } : {})
    }
  };

  return new Promise((resolve, reject) => {
    const req = http.request(url, options, (res) => {
      let data = '';
      res.on('data', chunk => data += chunk);
      res.on('end', () => {
        try {
          const parsed = data ? JSON.parse(data) : {};
          resolve({ status: res.statusCode, data: parsed, raw: data });
        } catch (e) {
          resolve({ status: res.statusCode, data: null, raw: data });
        }
      });
    });
    req.on('error', reject);
    if (body) req.write(JSON.stringify(body));
    req.end();
  });
}

async function seedData(token) {
  console.log('🌱 Populating database with rich, production-grade sample projects...\n');

  const projectsData = [
    {
      name: 'E-Commerce & Checkout Platform',
      description: 'Unified payment gateway, checkout V2 funnels, discount engine, and BNPL integrations',
      environments: ['Development', 'Staging', 'Production'],
      flags: [
        {
          key: 'new-stripe-elements-ui',
          name: 'New Stripe Elements UI',
          description: 'Enable Stripe Elements V3 unified checkout form',
          type: 'BOOLEAN'
        },
        {
          key: 'express-apple-pay',
          name: 'Express Apple Pay Checkout',
          description: 'One-click Apple Pay button on cart drawer',
          type: 'BOOLEAN',
          parentKey: 'new-stripe-elements-ui'
        },
        {
          key: 'checkout-funnel-variant',
          name: 'Checkout Funnel Variant',
          description: 'A/B Test 1-Step Single Page vs 2-Step vs Accordion Checkout Flow',
          type: 'MULTIVARIATE',
          variations: [
            { id: 'var_1step', name: '1-Step Single Page', value: '1-step-single-page' },
            { id: 'var_2step', name: '2-Step Multi Page', value: '2-step-multi-page' },
            { id: 'var_accordion', name: 'Collapsible Accordion', value: 'accordion-view' }
          ]
        },
        {
          key: 'cart-capacity-limits',
          name: 'Cart Capacity Limits',
          description: 'Max allowed items and auto-hold inventory thresholds',
          type: 'JSON',
          variations: [
            { id: 'var_standard', name: 'Standard Users', value: { maxItems: 50, autoHoldMins: 15 } },
            { id: 'var_vip', name: 'VIP Subscribers', value: { maxItems: 250, autoHoldMins: 60 } }
          ]
        },
        {
          key: 'enable-bnpl-klarna',
          name: 'Enable Klarna & Affirm BNPL',
          description: 'Show Buy-Now-Pay-Later payment options during checkout',
          type: 'BOOLEAN'
        },
        {
          key: 'realtime-shipping-rate-calculator',
          name: 'Realtime Shipping Rate Calculator',
          description: 'Fetch live FedEx, UPS, and DHL rates via API',
          type: 'BOOLEAN'
        }
      ]
    },
    {
      name: 'AI Search & Recommendation Engine',
      description: 'Vector embeddings, Gemini-powered semantic reranking, and search algorithm flags',
      environments: ['Development', 'Staging', 'Production'],
      flags: [
        {
          key: 'enable-gemini-reranking',
          name: 'Enable Gemini Semantic Reranking',
          description: 'Use Gemini LLM to semantically rerank search results',
          type: 'BOOLEAN'
        },
        {
          key: 'llm-query-expander',
          name: 'LLM Query Expander',
          description: 'Expand user search query with LLM-suggested domain synonyms',
          type: 'BOOLEAN',
          parentKey: 'enable-gemini-reranking'
        },
        {
          key: 'vector-search-model-version',
          name: 'Vector Search Model Version',
          description: 'Active embedding model for nearest-neighbor search',
          type: 'MULTIVARIATE',
          variations: [
            { id: 'var_bm25', name: 'BM25 Keyword Engine', value: 'v1-bm25-classic' },
            { id: 'var_bge', name: 'BGE Large Embeddings', value: 'v2-bge-large-en' },
            { id: 'var_gemini', name: 'Gemini Multimodal Embedding', value: 'v3-gemini-embed-002' }
          ]
        },
        {
          key: 'vector-nearest-neighbors-k',
          name: 'Vector Nearest Neighbors (Top-K)',
          description: 'JSON configuration for vector search top-K retrieval and similarity score cutoffs',
          type: 'JSON',
          variations: [
            { id: 'var_fast', name: 'Fast Search (K=10)', value: { topK: 10, minScore: 0.75 } },
            { id: 'var_deep', name: 'Deep Search (K=50)', value: { topK: 50, minScore: 0.60 } }
          ]
        }
      ]
    },
    {
      name: 'Mobile Banking & Digital Wallet',
      description: 'iOS & Android V2 redesign, biometric auth, Zelle integration, and transfer limits',
      environments: ['Alpha', 'Beta', 'Production'],
      flags: [
        {
          key: 'passkey-biometric-auth',
          name: 'Passkey & Biometric Auth',
          description: 'Enable WebAuthn, TouchID, and FaceID passwordless login',
          type: 'BOOLEAN'
        },
        {
          key: 'oled-dark-mode-v2',
          name: 'OLED Dark Mode V2',
          description: 'Redesigned true-black OLED dark theme',
          type: 'BOOLEAN'
        },
        {
          key: 'instant-p2p-zelle-transfers',
          name: 'Instant P2P Zelle Transfers',
          description: 'Realtime peer-to-peer Zelle money transfer service',
          type: 'BOOLEAN'
        },
        {
          key: 'daily-wire-transfer-config',
          name: 'Daily Wire Transfer Limits',
          description: 'Limits and secondary approval requirements per KYC tier',
          type: 'JSON',
          variations: [
            { id: 'tier_tier1', name: 'Tier 1 Verified', value: { dailyLimit: 10000, requireOTP: true } },
            { id: 'tier_tier2', name: 'Tier 2 Institutional', value: { dailyLimit: 250000, requireOTP: true } }
          ]
        }
      ]
    },
    {
      name: 'SaaS Analytics & Reporting Portal',
      description: 'Live event streaming, PDF report generation, anomaly detection, and retention rules',
      environments: ['Development', 'Production'],
      flags: [
        {
          key: 'realtime-websocket-feed',
          name: 'Realtime WebSocket Feed',
          description: 'Live event streaming in analytics dashboard',
          type: 'BOOLEAN'
        },
        {
          key: 'ai-anomaly-detection-alerts',
          name: 'AI Anomaly Detection Alerts',
          description: 'Automated alert trigger when traffic anomalies occur',
          type: 'BOOLEAN'
        },
        {
          key: 'audit-log-retention-config',
          name: 'Audit Log Retention Config',
          description: 'Retention periods and automated archiver rules',
          type: 'JSON',
          variations: [
            { id: 'ret_free', name: 'Free Tier (30 Days)', value: { retentionDays: 30, exportAllowed: false } },
            { id: 'ret_enterprise', name: 'Enterprise Tier (365 Days)', value: { retentionDays: 365, exportAllowed: true } }
          ]
        }
      ]
    }
  ];

  for (const projSpec of projectsData) {
    console.log(`🚀 Creating Project: "${projSpec.name}"...`);
    const projRes = await request('/projects', 'POST', {
      name: projSpec.name,
      description: projSpec.description
    }, token);

    if (projRes.status !== 201) {
      console.warn(`   ⚠️ Project "${projSpec.name}" returned ${projRes.status}`);
      continue;
    }

    const projectId = projRes.data.id;
    console.log(`   ✅ Project created. ID: ${projectId}`);

    // Create Environments
    for (const envName of projSpec.environments) {
      const envRes = await request(`/projects/${projectId}/environments`, 'POST', {
        name: envName,
        isProtected: envName === 'Production'
      }, token);

      if (envRes.status === 201) {
        console.log(`      🔹 Environment: ${envName}`);
      }
    }

    // Map created flags to handle dependencies
    const createdFlagsMap = {};

    // Create Flags
    for (const flagSpec of projSpec.flags) {
      const parentId = flagSpec.parentKey ? createdFlagsMap[flagSpec.parentKey] : undefined;

      const flagRes = await request(`/projects/${projectId}/flags`, 'POST', {
        projectId,
        key: flagSpec.key,
        name: flagSpec.name,
        description: flagSpec.description,
        type: flagSpec.type,
        variations: flagSpec.variations,
        parentFlagId: parentId
      }, token);

      if (flagRes.status === 201) {
        createdFlagsMap[flagSpec.key] = flagRes.data.id;
        const depInfo = parentId ? ` (Depends on parent flag)` : '';
        console.log(`      🚩 Flag created: ${flagSpec.name} [${flagSpec.type}]${depInfo}`);
      } else {
        console.warn(`      ⚠️ Flag "${flagSpec.key}" returned ${flagRes.status}:`, flagRes.raw);
      }
    }
    console.log('');
  }

  console.log('🎉 Seed Completed Successfully! Open http://localhost:3000 to view your new projects.');
}

async function main() {
  const loginRes = await request('/auth/login', 'POST', {
    email: 'admin@example.com',
    password: 'admin123'
  });

  if (loginRes.status !== 200 || !loginRes.data.token) {
    console.error('❌ Authentication failed. Is FlagManagment running?', loginRes.status, loginRes.raw);
    process.exit(1);
  }

  await seedData(loginRes.data.token);
}

if (require.main === module) {
  main();
}

module.exports = { request, seedData };
