# FlagManagment — Node.js & Express Sample

This is a complete, runnable Express.js application demonstrating how to integrate the **FlagManagment Node SDK** (`@flagmanagment/node-sdk`) using the **CNCF OpenFeature** standard.

---

## Features Demonstrated

1. **CNCF OpenFeature Provider**: Initializing `FlagManagmentProvider` with local evaluation caching and SSE delta streams.
2. **Context-Aware Targeting**: Passing `EvaluationContext` (`targetingKey`, user tier, location) to dynamically branch routes.
3. **Multivariate & String Variants**: Evaluating non-boolean variants (e.g., A/B/C banner experiments).
4. **Resilient Fallbacks**: Graceful default value return when network or backend connections are unreachable.

---

## Quickstart

### 1. Install Dependencies
```bash
npm install
```

### 2. Configure Environment (Optional)
By default, the sample connects to `http://localhost:8080`. You can configure your own server key and API host:
```bash
export FLAGMANAGMENT_SERVER_KEY="flg_srv_your_server_key"
export FLAGMANAGMENT_API_URL="http://localhost:8080"
```

### 3. Run the Server
```bash
npm start
```
The server will start at `http://localhost:4000`.

---

## Test Endpoints

- **Boolean Flag Evaluation with User Targeting**:
  ```bash
  curl "http://localhost:4000/api/checkout?userId=vip-42&tier=premium"
  ```
  *Response*:
  ```json
  {
    "flow": "v2-one-click-checkout",
    "features": ["apple-pay", "google-pay", "instant-address-autofill"],
    "context": { "targetingKey": "vip-42", "tier": "premium", "country": "US" }
  }
  ```

- **Multivariate A/B/C Banner**:
  ```bash
  curl "http://localhost:4000/api/banner?userId=anon-99"
  ```

- **Detailed Resolution Playground**:
  ```bash
  curl -X POST http://localhost:4000/api/evaluate \
    -H "Content-Type: application/json" \
    -d '{"flagKey": "new-checkout-v2", "defaultValue": false, "context": {"targetingKey": "user-1"}}'
  ```

