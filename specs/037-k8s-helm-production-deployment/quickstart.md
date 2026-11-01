# Quickstart: Deploying FlagManagment via Helm

**Branch**: `037-k8s-helm-production-deployment` | **Date**: 2026-10-31

This guide demonstrates how to validate, install, and upgrade the FlagManagment Helm chart.

---

## 1. Lint and Template Dry-Run

```bash
# 1. Lint chart syntax
helm lint deploy/helm/flagmanagment

# 2. Render templates locally
helm template my-release deploy/helm/flagmanagment --debug
```

---

## 2. Deploy to a Cluster (Local k3s / Minikube or Production)

### Option A: Standalone Deployment (Embedded Postgres & Redis)
```bash
helm install flagmanagment ./deploy/helm/flagmanagment \
  --create-namespace \
  --namespace flagmanagment
```

### Option B: Production Deployment (External Managed Databases)
Create custom `prod-values.yaml`:
```yaml
postgresql:
  enabled: false

externalDatabase:
  host: "postgres.mycompany.internal"
  port: 5432
  user: "flag_user"
  database: "flagmanagment"
  passwordSecret: "managed-db-secret"
  passwordKey: "password"

redis:
  enabled: false

externalRedis:
  host: "redis.mycompany.internal"
  port: 6379

ingress:
  enabled: true
  className: "nginx"
  hosts:
    - host: flags.mycompany.com
      paths:
        - path: /api
          pathType: Prefix
          service: backend
        - path: /
          pathType: Prefix
          service: frontend
  tls:
    - secretName: flags-tls-cert
      hosts:
        - flags.mycompany.com
```

Deploy with custom values:
```bash
helm upgrade --install flagmanagment ./deploy/helm/flagmanagment \
  -f prod-values.yaml \
  --namespace flagmanagment
```

---

## 3. Verify Deployment

```bash
kubectl get pods,svc,ingress,hpa -n flagmanagment
```

