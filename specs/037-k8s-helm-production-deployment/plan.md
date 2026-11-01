# Implementation Plan: Kubernetes Helm Chart & Production Deployment Manifests

**Branch**: `037-k8s-helm-production-deployment` | **Date**: 2026-10-31 | **Spec**: [specs/037-k8s-helm-production-deployment/spec.md](spec.md)

**Input**: Feature specification from `/specs/037-k8s-helm-production-deployment/spec.md` (FR-001 through FR-006, SC-001 through SC-004)

---

## Summary

This plan designs the production-grade Helm v3 chart package located at `deploy/helm/flagmanagment`. The package provides complete Kubernetes manifests to self-host FlagManagment with zero manual infrastructure assembly:
- **Backend Deployment**: High-availability Go backend pods with horizontal pod autoscaling (HPA), liveness (`/healthz`) and readiness (`/readyz`) probes, and resource bounds.
- **Frontend Deployment**: React dashboard served via optimized Nginx container with reverse proxy routing to the backend.
- **Pre-Install/Pre-Upgrade Migration Hook**: Kubernetes Job running `golang-migrate` schema migrations against PostgreSQL before rolling deployment traffic.
- **Ingress Controller**: Path-based routing for `/api`, `/grpc`, and root UI with TLS secret termination and cert-manager integration.
- **Datastore Options**: Standalone embedded PostgreSQL 16 & Redis 7 for quick testing or seamless integration with external managed datastores (AWS RDS, GCP Cloud SQL, Redis Enterprise).

---

## Technical Context

**Language/Version**: Helm v3 (YAML Go templating, Kubernetes API v1, apps/v1, networking.k8s.io/v1, autoscaling/v2)  
**Primary Dependencies**: Kubernetes 1.24+, Helm 3.8+, Bitnami PostgreSQL sub-chart (optional dependency), Bitnami Redis sub-chart (optional dependency)  
**Storage**: PersistentVolumeClaim (PVC) for embedded PostgreSQL, ConfigMap & Secret for runtime environment variables  
**Testing**: `helm lint deploy/helm/flagmanagment`, `helm template flagmanagment ./deploy/helm/flagmanagment`, `kubeval`/`conftest`  
**Target Platform**: Any standard Kubernetes cluster (EKS, GKE, AKS, k3s, Minikube, OpenShift)  
**Project Type**: Infrastructure as Code / Helm Chart Packaging  
**Performance Goals**: Sub-3 minute complete deployment; zero-downtime rolling upgrades  
**Constraints**: Support both x86_64 and ARM64 container architectures; support air-gapped private registries via `imagePullSecrets`  

---

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Requirement | Plan Alignment | Status |
|---|---|---|---|
| **VIII. Cloud-Native Portability** | All services containerized with reference Kubernetes manifests or Helm charts | Creates official `deploy/helm/flagmanagment` chart fulfilling Principle VIII | ✅ PASS |
| **II. Environment Isolation** | Secure tokens, no cross-environment leakage | Secrets and ConfigMaps strictly isolated with Kubernetes RBAC and Secret masking | ✅ PASS |
| **III. Governance by Default** | Unlimited projects/environments, no feature caps in self-hosted | Full platform components deployed out-of-the-box with no enterprise locks | ✅ PASS |
| **V. Test-First Quality Gates** | Validated via static analysis and automated checks | Validated by `helm lint` and `helm template` dry runs | ✅ PASS |

---

## Project Structure

### Documentation (this feature)

```text
specs/037-k8s-helm-production-deployment/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output (Values hierarchy & template mapping)
├── quickstart.md        # Phase 1 output (Installation & upgrade guide)
└── tasks.md             # Phase 2 output
```

### Source Code (repository root)

```text
deploy/
└── helm/
    └── flagmanagment/
        ├── Chart.yaml                  # Chart metadata and dependencies
        ├── values.yaml                 # Default configuration values
        ├── .helmignore                 # Packaged file exclusions
        └── templates/
            ├── _helpers.tpl            # Standard template naming helpers
            ├── configmap.yaml          # Shared backend/frontend ConfigMaps
            ├── secrets.yaml            # JWT secrets, DB passwords, salt keys
            ├── serviceaccount.yaml     # Kubernetes ServiceAccount definition
            ├── migration-job.yaml      # Pre-upgrade hook job for SQL migrations
            ├── backend-deployment.yaml # Go API & gRPC backend deployment
            ├── backend-service.yaml    # Service exposing HTTP (8080) and gRPC (9090)
            ├── backend-hpa.yaml        # HorizontalPodAutoscaler for backend
            ├── frontend-deployment.yaml# React UI deployment
            ├── frontend-service.yaml   # Service exposing UI (3000)
            ├── ingress.yaml            # Ingress rules with TLS & path routing
            ├── postgres-statefulset.yaml # Optional embedded PostgreSQL 16
            └── redis-deployment.yaml   # Optional embedded Redis 7
```

**Structure Decision**: The Helm chart lives directly in `deploy/helm/flagmanagment`, following the standard enterprise repository convention.

---

## Complexity Tracking

*No violations. Follows standard CNCF Helm 3 Best Practices.*
