# Research: Kubernetes Helm Architecture & Enterprise Deployment Decisions

**Branch**: `037-k8s-helm-production-deployment` | **Date**: 2026-10-31

This research outlines architectural choices for the Helm v3 chart at `deploy/helm/flagmanagment`.

---

## Decision 1: Database Migration Execution Pattern

- **Decision**: Use a Kubernetes `Job` configured with Helm hooks:
  ```yaml
  annotations:
    "helm.sh/hook": pre-install,pre-upgrade
    "helm.sh/hook-weight": "1"
    "helm.sh/hook-delete-policy": before-hook-creation
  ```
- **Rationale**: Running migrations directly inside application containers at startup causes race conditions during multi-replica rolling deployments. A pre-install/pre-upgrade Hook Job ensures that the database schema is up-to-date before any new backend pods receive incoming HTTP/gRPC traffic.
- **Alternatives Considered**:
  - *InitContainers on the backend pod*: Multiple replicas initializing concurrently compete for migration locks and can cause deadlocks.

---

## Decision 2: Datastore Flexibility (Embedded vs External)

- **Decision**: Support two modes in `values.yaml`:
  1. `embedded`: Local PostgreSQL StatefulSet + Redis Deployment with PVC for quick POCs / Minikube / k3s.
  2. `external`: Pointers to managed Cloud SQL / RDS PostgreSQL and Redis cluster with user-supplied Secret references.
- **Rationale**: Enterprise customers never deploy production stateful databases on bare pods without dedicated managed services; developers want a single command `helm install` without provisioning external cloud accounts.

---

## Decision 3: Ingress & gRPC Protocol Handling

- **Decision**: Use standard Kubernetes `networking.k8s.io/v1` `Ingress` with customizable annotations:
  - Supports `ingress-nginx`, AWS ALB Ingress Controller, and Traefik.
  - Automatically configures path routing:
    - `/api` -> `backend:8080`
    - `/grpc` -> `backend:9090`
    - `/` -> `frontend:3000`
- **Rationale**: Path-based routing prevents CORS issues between dashboard and API while allowing single-domain hosting (`flags.myorg.internal`).

