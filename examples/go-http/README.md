# FlagManagment — Go HTTP Microservice Sample

This is a complete, runnable Go HTTP microservice demonstrating how to integrate the **FlagManagment Go SDK** (`sdk/go`) using the **CNCF OpenFeature Go SDK**.

---

## Features Demonstrated

1. **Sub-Microsecond In-Memory Evaluation**: Zero network hops per evaluation using local cached rule evaluation.
2. **CNCF OpenFeature Provider**: Global provider registration with OpenFeature client APIs.
3. **Background SSE Streaming**: Automatic live flag updates from FlagManagment backend or Edge Proxy.
4. **Resilient Local Fallback**: Fast default evaluation when upstream connections are severed.

---

## Quickstart

### 1. Download Dependencies
```bash
go mod tidy
```

### 2. Configure Environment (Optional)
```bash
export FLAGMANAGMENT_SERVER_KEY="flg_srv_your_go_key"
export FLAGMANAGMENT_API_URL="http://localhost:8080"
```

### 3. Run the Microservice
```bash
go run main.go
```
The server will start at `http://localhost:5000`.

---

## Test Endpoints

- **Health Check**:
  ```bash
  curl http://localhost:5000/health
  ```

- **User Context Evaluation**:
  ```bash
  curl "http://localhost:5000/api/checkout?userId=alice&tier=premium"
  ```

- **Detailed Evaluation Resolution**:
  ```bash
  curl "http://localhost:5000/api/evaluate?flag=new-checkout-v2&userId=alice"
  ```

---

## Run Unit Tests
```bash
go test -v ./...
```

