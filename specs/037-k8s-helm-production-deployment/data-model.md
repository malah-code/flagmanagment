# Data Model: Helm Values Hierarchy & Resource Mapping

**Branch**: `037-k8s-helm-production-deployment` | **Date**: 2026-10-31

This document defines the schema of `values.yaml` and maps configuration keys to Kubernetes resources.

---

## 1. Values Schema Hierarchy

| Section | Key | Type | Default | Description |
|---|---|---|---|---|
| **global** | `imageRegistry` | string | `""` | Global image registry override |
| **global** | `imagePullSecrets`| list | `[]` | Credentials for private registries |
| **backend** | `replicaCount` | int | `2` | Number of backend pod replicas |
| **backend** | `image.repository` | string | `flagmanagment/backend` | Container image repository |
| **backend** | `image.tag` | string | `latest` | Container image tag |
| **backend** | `resources` | object | `{requests: {cpu: 100m, memory: 128Mi}}` | CPU and Memory bounds |
| **backend** | `autoscaling.enabled` | bool | `true` | Enable HorizontalPodAutoscaler |
| **backend** | `autoscaling.maxReplicas` | int | `10` | Maximum auto-scaled replicas |
| **frontend** | `replicaCount` | int | `2` | Number of frontend pod replicas |
| **frontend** | `image.repository` | string | `flagmanagment/frontend` | Container image repository |
| **ingress** | `enabled` | bool | `true` | Expose services via Kubernetes Ingress |
| **ingress** | `className` | string | `"nginx"` | Ingress class |
| **ingress** | `hosts` | list | `[host: flagmanagment.local]` | Domain hosts |
| **ingress** | `tls` | list | `[]` | TLS certificate secret bindings |
| **postgresql** | `enabled` | bool | `true` | Deploy embedded PostgreSQL StatefulSet |
| **externalDatabase** | `host` | string | `""` | External PostgreSQL hostname |
| **redis** | `enabled` | bool | `true` | Deploy embedded Redis Deployment |
| **externalRedis** | `host` | string | `""` | External Redis hostname |

