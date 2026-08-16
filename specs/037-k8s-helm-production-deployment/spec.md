# Feature Specification: Kubernetes Helm Chart & Production Deployment Manifests

**Feature Branch**: `037-k8s-helm-production-deployment`

**Created**: 2026-08-16

**Status**: Ready for Planning

**Input**: User description: "if any missing or need enhance them create spec for it - Kubernetes Helm charts and production self-hosted deployment packages (PRD Section 5.2, 11.2)"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - One-Command Production Helm Deployment (Priority: P1)

As a DevOps Engineer or Platform Lead, I want to deploy FlagManagment to a production Kubernetes cluster (or local minikube/k3s) using a standard Helm chart with customizable `values.yaml`, so that our organization can self-host the platform with high availability, persistent volume claims, automated database migrations, and Redis clustering.

**Why this priority**: Directly fulfills the Section 5.2 requirement in `docs/p-requirements.md` ("Kubernetes manifests or Helm charts for reference deployments") enabling enterprise self-hosting at scale.

**Independent Test**: Can be validated by executing `helm lint` and `helm template`, deploying to a test cluster/k3s, and verifying that all services (backend, frontend, postgres, redis) initialize, run migrations, and pass liveness/readiness probes.

**Acceptance Scenarios**:

1. **Given** a Kubernetes cluster, **When** running `helm install flagmanagment ./deploy/helm/flagmanagment`, **Then** the backend, frontend, PostgreSQL, and Redis pods are created with configured replica counts and reach `Running` and `Ready` states.
2. **Given** an existing PostgreSQL database, **When** configuring external database credentials in `values.yaml`, **Then** the Helm chart runs the automated migration job and connects the backend without launching an embedded database pod.

---

### User Story 2 - Ingress & TLS Termination Configuration (Priority: P2)

As a System Administrator, I want to configure standard Kubernetes Ingress rules with TLS termination and path-based routing (routing `/api/*` and `/grpc.*` to the backend and root traffic to the frontend dashboard), so that external users and SDKs can securely connect over HTTPS.

**Why this priority**: Essential for secure enterprise deployments without manual proxy setup.

**Independent Test**: Can be tested by rendering the Ingress manifest with TLS annotations and verifying host routing and certificate secrets.

**Acceptance Scenarios**:

1. **Given** `ingress.enabled: true` in `values.yaml`, **When** the chart is deployed, **Then** an Ingress resource with specified hostnames, annotations (e.g. `cert-manager.io/cluster-issuer`), and TLS secret configurations is created.

---

### User Story 3 - Health Probes, Resource Limits, and Autoscaling (Priority: P3)

As an SRE, I want configurable CPU/Memory resource requests/limits, HorizontalPodAutoscalers (HPA), and liveness/readiness probes targeting `/healthz` and `/readyz`, so that the platform scales reliably under peak SDK traffic and self-heals during pod restarts.

**Why this priority**: Guarantees zero-downtime rolling upgrades and resilience for millions of flag evaluations.

**Independent Test**: Can be tested by triggering simulated CPU load or crashing a backend pod to verify Kubernetes restarts the pod and passes readiness checks.

**Acceptance Scenarios**:

1. **Given** `autoscaling.enabled: true` with target CPU utilization at 70%, **When** traffic increases, **Then** the HPA controller scales backend pods up to the defined max replicas.

---

### Edge Cases

- **Database Connection Failure at Startup**: The backend pod must wait for the database migration job to complete or retry with exponential backoff rather than crash-looping indefinitely.
- **Air-Gapped / Private VPC Deployments**: The chart must support pulling images from private registries via `imagePullSecrets`.
- **ConfigMap & Secret Updates**: Changes in configuration values must trigger rolling pod updates without downtime.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST provide a structured Helm v3 chart under `deploy/helm/flagmanagment` containing templates for Backend, Frontend, Ingress, ConfigMaps, Secrets, ServiceAccounts, and HPA.
- **FR-002**: The Helm chart MUST provide a comprehensive `values.yaml` with clear defaults for standalone development (embedded Postgres & Redis) and enterprise production (external managed DB & Redis).
- **FR-003**: The chart MUST include a Kubernetes Pre-Install/Pre-Upgrade Hook Job that runs database migrations before new backend pods are routed traffic.
- **FR-004**: System MUST define configurable Liveness (`/healthz`) and Readiness (`/readyz`) probes on all service deployments.
- **FR-005**: System MUST support Ingress definitions with TLS secret bindings, gRPC streaming annotations (e.g. `nginx.ingress.kubernetes.io/backend-protocol: "GRPC"`), and customizable hostnames.
- **FR-006**: System MUST support Horizontal Pod Autoscalers (HPA) targeting CPU and memory utilization thresholds.

### Key Entities

- **Helm Release**: Represents the installed instance of the FlagManagment platform on a Kubernetes cluster.
- **Values Configuration (`values.yaml`)**: The declarative configuration defining replica counts, container image tags, environment secrets, and storage classes.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can deploy a complete, functional self-hosted FlagManagment cluster with a single `helm install` command in under 3 minutes.
- **SC-002**: `helm lint deploy/helm/flagmanagment` passes with 0 errors and 0 warnings.
- **SC-003**: Rolling upgrades complete with zero dropped SDK connections or user requests.
- **SC-004**: Deployment supports both x86_64 and ARM64 container architectures.

## Assumptions

- Users have Helm v3.8+ and access to a Kubernetes cluster (v1.24+).
- Ingress controllers (such as ingress-nginx or cloud load balancers) and CSI storage providers are available in target production clusters.
