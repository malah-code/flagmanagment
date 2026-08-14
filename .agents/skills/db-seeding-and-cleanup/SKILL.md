---
name: db-seeding-and-cleanup
description: Guidance on seeding realistic sample projects/flags and resetting database state in FlagManagment. Use whenever asked to seed sample data, reset projects, or clean up test data.
---

# Database Seeding and Cleanup Skill

This skill documents how to seed realistic sample data into FlagManagment, reset the project state, and manage database cleanup procedures.

---

## ⚡ Quick Reference Commands

| Command | Purpose | Description |
| :--- | :--- | :--- |
| `make seed` | Append Sample Data | Runs `node scripts/seed.js` to create realistic enterprise projects, environments, and flags via API. |
| `make seed-reset` or `make reseed` | Clean & Seed Fresh | Runs `node scripts/reset-seed.js` to delete all existing projects and seed fresh data. |
| `make clean` | Hard Storage Reset | Runs `docker compose down -v` to destroy PostgreSQL & Redis volumes. |
| `make up` | Start Stack & Migrate | Boots up Docker containers, applies database migrations, and seeds the initial admin user. |

---

## 🔑 Default Admin Credentials

All API-based seeding scripts and manual testing use the seed admin user:
- **Email**: `admin@example.com`
- **Password**: `admin123`

---

## 🏗️ Seeding Architecture & Workflows

### 1. Incremental Seeding (`scripts/seed.js`)
Creates 4 realistic enterprise projects with varied flag types (`BOOLEAN`, `MULTIVARIATE`, `JSON`) and parent-child dependencies:
- **E-Commerce & Checkout Platform**: Checkout V2 funnels, Stripe Elements V3, BNPL, cart limits (JSON).
- **AI Search & Recommendation Engine**: Gemini semantic reranking, vector search model selection, top-K retrieval parameters.
- **Mobile Banking & Digital Wallet**: Passkeys, OLED dark mode, Zelle P2P, daily wire transfer limits.
- **SaaS Analytics & Reporting Portal**: Realtime WebSocket streaming, AI anomaly alerts, audit retention rules.

### 2. API-Level Reset & Reseed (`scripts/reset-seed.js`)
1. Authenticates against `/api/v1/auth/login`.
2. Queries `/api/v1/projects` to list all existing user/test projects.
3. Issues `DELETE /api/v1/projects/:id` for each project (cascades to delete environments, flags, states, rules, and audit entries).
4. Executes `seedData()` to re-populate the standard sample suite.

### 3. Full Infrastructure Hard Reset
When schema migrations or DB constraints change:
```bash
make clean && make up
```
Followed by:
```bash
make seed
```
