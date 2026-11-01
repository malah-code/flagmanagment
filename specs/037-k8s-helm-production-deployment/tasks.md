# Tasks: Kubernetes Helm Chart & Production Deployment Manifests

**Input**: Design documents from `/specs/037-k8s-helm-production-deployment/`
**Prerequisites**: [plan.md](plan.md), [spec.md](spec.md), [research.md](research.md), [data-model.md](data-model.md), [quickstart.md](quickstart.md)

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Chart scaffolding and Helm v3 metadata

- [x] T001 Initialize Helm chart directory structure in deploy/helm/flagmanagment/
- [x] T002 Create Helm Chart metadata in deploy/helm/flagmanagment/Chart.yaml
- [x] T003 [P] Create Helm file exclusions in deploy/helm/flagmanagment/.helmignore
- [x] T004 [P] Create standard template naming helpers in deploy/helm/flagmanagment/templates/_helpers.tpl

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Shared ConfigMaps, Secrets, RBAC, and Pre-migration Hook

- [x] T005 Create shared ConfigMap template in deploy/helm/flagmanagment/templates/configmap.yaml
- [x] T006 [P] Create Secret template for credentials in deploy/helm/flagmanagment/templates/secrets.yaml
- [x] T007 [P] Create ServiceAccount template in deploy/helm/flagmanagment/templates/serviceaccount.yaml
- [x] T008 Create pre-install/pre-upgrade database migration Job in deploy/helm/flagmanagment/templates/migration-job.yaml

---

## Phase 3: User Story 1 - One-Command Production Helm Deployment (Priority: P1) 🎯 MVP

**Goal**: Enable deployment of full FlagManagment stack with backend, frontend, postgres, and redis via `helm install`.

**Independent Test**: Running `helm lint deploy/helm/flagmanagment` and `helm template test deploy/helm/flagmanagment` renders all workloads with zero errors.

- [x] T009 [US1] Create backend Deployment template in deploy/helm/flagmanagment/templates/backend-deployment.yaml
- [x] T010 [P] [US1] Create backend Service template in deploy/helm/flagmanagment/templates/backend-service.yaml
- [x] T011 [P] [US1] Create frontend Deployment template in deploy/helm/flagmanagment/templates/frontend-deployment.yaml
- [x] T012 [P] [US1] Create frontend Service template in deploy/helm/flagmanagment/templates/frontend-service.yaml
- [x] T013 [P] [US1] Create embedded PostgreSQL StatefulSet template in deploy/helm/flagmanagment/templates/postgres-statefulset.yaml
- [x] T014 [P] [US1] Create embedded Redis Deployment template in deploy/helm/flagmanagment/templates/redis-deployment.yaml
- [x] T015 [US1] Create complete values configuration in deploy/helm/flagmanagment/values.yaml

---

## Phase 4: User Story 2 - Ingress & TLS Termination Configuration (Priority: P2)

**Goal**: Expose API, gRPC, and Dashboard through configurable Kubernetes Ingress with TLS.

**Independent Test**: Rendering with `ingress.enabled=true` produces valid `networking.k8s.io/v1` Ingress resource.

- [x] T016 [US2] Create Ingress template with TLS annotations in deploy/helm/flagmanagment/templates/ingress.yaml

---

## Phase 5: User Story 3 - Health Probes, Resource Limits, and Autoscaling (Priority: P3)

**Goal**: Configure resource requests/limits, probes, and HorizontalPodAutoscalers for high-load scaling.

**Independent Test**: Rendering with `backend.autoscaling.enabled=true` produces `autoscaling/v2` HorizontalPodAutoscaler.

- [x] T017 [US3] Create HorizontalPodAutoscaler template in deploy/helm/flagmanagment/templates/backend-hpa.yaml

---

## Phase 6: Polish & Verification

**Purpose**: Chart validation, dry runs, and documentation

- [x] T018 Validate chart syntax using helm lint deploy/helm/flagmanagment
- [x] T019 Validate template rendering using helm template test deploy/helm/flagmanagment
- [x] T020 Add Helm deployment instructions to docs/GETTING_STARTED.md
