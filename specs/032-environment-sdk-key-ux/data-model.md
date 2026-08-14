# Data Model: Environment SDK Key UX & Integration Guide

No database schema changes required. The `Environment` model already contains `client_key` (or API key token) associated with the environment entity.

- **Environment**:
  - `id`: UUID
  - `name`: String
  - `client_key` / `apiKey`: String (Exposed to authenticated dashboard users for frontend SDK setup)
