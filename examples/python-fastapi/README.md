# FlagManagment — Python & FastAPI Sample

This is a complete, runnable FastAPI web service demonstrating how to integrate the **FlagManagment Python SDK** using the **CNCF OpenFeature** standard.

---

## Features Demonstrated

1. **CNCF OpenFeature Python Provider**: Configuring `FlagManagmentProvider` on the global `openfeature.api`.
2. **Contextual User Evaluation**: Passing `EvaluationContext` with user keys, subscription tiers, and geolocation metadata.
3. **Multivariate AI Model Selection**: Evaluating dynamic model names and variants without service redeployment.
4. **Resilient Local Evaluation**: Default fallbacks ensured under any network outage or latency spike.

---

## Quickstart

### 1. Install Dependencies
```bash
pip install -r requirements.txt
```

### 2. Configure Environment (Optional)
```bash
export FLAGMANAGMENT_SERVER_KEY="flg_srv_your_python_key"
export FLAGMANAGMENT_API_URL="http://localhost:8080"
```

### 3. Run the FastAPI Application
```bash
uvicorn main:app --reload --port 8000
```
FastAPI interactive Swagger UI will be available at: `http://localhost:8000/docs`

---

## Test Endpoints

- **Features Matrix with Targeting**:
  ```bash
  curl "http://localhost:8000/api/features?user_id=usr-101&tier=enterprise"
  ```
- **Dynamic Recommendation Algorithm Gate**:
  ```bash
  curl "http://localhost:8000/api/recommendations?user_id=beta-tester-1"
  ```

---

## Running Tests
```bash
pytest test_main.py
```

