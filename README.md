# 🚩 FlagManagment

[![CI Status](https://img.shields.io/badge/build-passing-brightgreen.svg)]()
[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![React Version](https://img.shields.io/badge/React-18-61DAFB?style=flat&logo=react)](https://reactjs.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16+-4169E1?style=flat&logo=postgresql)](https://www.postgresql.org/)
[![Redis](https://img.shields.io/badge/Redis-7+-DC382D?style=flat&logo=redis)](https://redis.io/)
[![License](https://img.shields.io/badge/License-BSL_1.1_(Fair--Source)-blue.svg)](LICENSE)

**FlagManagment** is an enterprise-grade, cloud-native feature flag and remote configuration platform. It provides deterministic evaluations, multi-environment isolation, contextual targeting, emergency kill switches, four-eyes governance change requests, and tamper-evident audit logging with zero vendor lock-in.

---

## ✨ Key Features

- 🚩 **Multi-Type Feature Flags**: First-class support for `BOOLEAN`, `MULTIVARIATE`, and dynamic `JSON` remote configuration flags with tags and lifecycle states.
- 🎯 **Contextual Targeting & Rollouts**: Granular attribute-based rules (`EQUALS`, `CONTAINS`, `REGEX`, etc.) and deterministic percentage rollouts powered by MurmurHash3.
- 🛡️ **Environment Isolation & Governance**:
  - Unlimited isolated environments (`Development`, `Staging`, `Production`) with unique cryptographic SDK server keys.
  - **Protected Environments**: Mandatory **Change Requests** with side-by-side visual diffs and four-eyes review controls prior to production promotion.
- ⚡ **Real-Time Synchronization**: Ultra-fast rule broadcast to streaming server-side and client-side SDKs over Redis Pub/Sub.
- 🚨 **Emergency Controls**: One-click kill switches and automated scheduled state transitions (`ENABLE`/`DISABLE` at future timestamps).
- 👥 **Role-Based Access Control (RBAC)**: Fine-grained access control (`ADMIN`, `RELEASE_MANAGER`, `EDITOR`, `VIEWER`), team member invitation tokens, and user access management.
- 📜 **Tamper-Evident Audit Logging**: Chronological, queryable audit trails capturing actors, IPs, target resources, and JSON state diffs for compliance.
- ⚙️ **System Configuration**: Self-hosted SMTP mail server configuration with encrypted credentials and live connectivity testing.

---

## 🏛️ Architecture

```
                      +-----------------------------+
                      |      React Frontend         |
                      |  (Vite + TS + TailwindCSS)  |
                      +--------------+--------------+
                                     |
                             REST / JSON API
                                     |
                                     v
                      +-----------------------------+
                      |         Go Backend          |
                      |    (Chi Router + RBAC)      |
                      +-------+--------------+------+
                              |              |
                      SQL / Migrations  Pub/Sub Sync
                              |              |
                              v              v
                      +---------------+ +----------+
                      | PostgreSQL 16 | | Redis 7+ |
                      +---------------+ +----------+
```

---

## 🚀 Quick Start (Docker Compose)

### 1. Clone and Start Stack
```bash
git clone https://github.com/malah-code/flagmanagment.git
cd flagmanagment
docker compose up -d --build
```

### 2. Access the Application
- **Frontend Dashboard**: [http://localhost:3000](http://localhost:3000)
- **Backend API**: [http://localhost:8080](http://localhost:8080)
- **Health Endpoint**: `curl http://localhost:8080/healthz`

### 3. Default Credentials
| Account | Email | Password | Default Role |
| :--- | :--- | :--- | :--- |
| **System Administrator** | `admin@example.com` | `admin123` | `ADMIN` (Global) |
| **Developer / Editor** | `dev@example.com` | `admin123` | `EDITOR` |

---

## 📡 API Overview

| Method | Route | Description | Auth Required |
| :--- | :--- | :--- | :--- |
| `POST` | `/api/v1/auth/login` | Authenticate user and obtain JWT token | No |
| `GET` | `/api/v1/projects` | List all projects | Yes (`VIEWER`) |
| `POST` | `/api/v1/projects` | Create a new workspace project | Yes (`ADMIN`) |
| `GET` | `/api/v1/projects/:id/flags` | List feature flag definitions | Yes (`VIEWER`) |
| `POST` | `/api/v1/projects/:id/flags` | Create a feature flag | Yes (`EDITOR`) |
| `PUT` | `/api/v1/projects/:id/environments/:envId/flags/:flagId/state` | Update flag state / targeting rules | Yes (`EDITOR`) |
| `POST` | `/api/v1/projects/:id/flags/:flagId/promote` | Promote flag across environments | Yes (`EDITOR`) |
| `GET` | `/api/v1/environments/:envId/change-requests` | List environment change requests | Yes (`VIEWER`) |
| `POST` | `/api/v1/change-requests/:id/approve` | Approve & apply change request | Yes (`RELEASE_MANAGER`) |
| `GET` | `/api/v1/projects/:id/audit-logs` | Retrieve project audit trail | Yes (`VIEWER`) |
| `GET/PUT`| `/api/v1/config/smtp` | Manage system email server configuration | Yes (`ADMIN`) |
| `POST` | `/api/v1/users/invite` | Issue invitation link with secure token | Yes (`ADMIN`) |

---

## 🛠️ Local Development

### Backend (Go)
```bash
cd backend
go run cmd/server/main.go
```
*Note: Hot reloading is configured via `air` in the Docker setup.*

### Frontend (React + Vite)
```bash
cd frontend
npm install
npm run dev
```

### Running Tests
```bash
# Backend unit tests
cd backend && go test ./...

# Database migrations check
docker compose logs -f postgres
```

---

## 📄 License
Distributed under the **Business Source License 1.1 (BSL 1.1)** / Fair-Source model. Free for internal self-hosted production and development use with zero artificial feature caps, while preventing third parties from reselling the codebase as a competing public SaaS. See [`LICENSE`](LICENSE) for complete terms.
