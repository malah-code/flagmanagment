# FlagManagment — Examples & Quickstart Blueprints

This directory contains production-ready blueprints and runnable application samples demonstrating how to integrate FlagManagment across modern tech stacks and infrastructure.

---

## 📱 Application SDK Samples

End-to-end runnable applications demonstrating how to evaluate feature flags with user targeting contexts, OpenFeature standards, and real-time streaming updates:

| Sample Application | Stack / SDK | Standards | Description |
|---|---|---|---|
| [`examples/node-express/`](node-express/) | Node.js, Express | CNCF OpenFeature (`@flagmanagment/node-sdk`) | Express API with contextual user evaluations, A/B banner variants, and fallback resilience. |
| [`examples/python-fastapi/`](python-fastapi/) | Python 3.10+, FastAPI | CNCF OpenFeature (`flagmanagment-sdk`) | High-performance FastAPI service with evaluation context injection and dynamic algorithm gates. |
| [`examples/go-http/`](go-http/) | Go 1.23+, `net/http` | CNCF OpenFeature (`sdk/go`) | Low-latency HTTP service with local in-memory evaluation (<5 µs) and background SSE sync. |
| [`examples/react-vite/`](react-vite/) | React 19, TypeScript, Vite | React Context & Hooks (`@flagmanagment/react-sdk`) | Interactive SPA with `<FlagProvider>`, `useFlag()` hook, dynamic UI toggles, and live identity switcher. |

---

## 🏗️ Infrastructure & Operations Blueprints

Declarative infrastructure-as-code and cloud-native deployment patterns:

| Blueprint | Category | Description |
|---|---|---|
| [`examples/terraform/`](terraform/) | Infrastructure as Code | Complete HashiCorp Terraform module provisioning projects, environments, flags, states, and service accounts. |
| [`examples/kubernetes/`](kubernetes/) | Kubernetes Sidecar | Multi-container pod manifest deploying Edge Proxy as an in-pod zero-hop loopback sidecar (<10 µs latency). |

---

## Running All Samples

Each sample contains a dedicated `README.md` with step-by-step instructions. For automated verification:

```bash
# Node Express
cd examples/node-express && npm test

# Python FastAPI
cd examples/python-fastapi && pytest test_main.py

# Go HTTP
cd examples/go-http && go test ./...

# React Vite
cd examples/react-vite && npm run build
```

